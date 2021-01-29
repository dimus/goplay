// Package contextimpl is for understanding how stdlib context package works.
package contextimpl

import (
	"errors"
	"reflect"
	"sync"
	"time"
)

// This is the standard interface of Context
type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}

// We use int as type so we can TODO and Background contexts are not identical.
type emptyCtx int

func (e emptyCtx) Deadline() (deadline time.Time, ok bool) {
	return
}
func (e emptyCtx) Done() <-chan struct{} {
	return nil
}
func (e emptyCtx) Err() error {
	return nil
}
func (e emptyCtx) Value(key interface{}) interface{} {
	return nil
}

var (
	// having background and todo as vars allows to allocate mem for them only
	// one time.
	background = new(emptyCtx)
	todo       = new(emptyCtx)
	// ErrCanceled is error for WithCancel
	ErrCanceled = errors.New("context canceled")
	// ErrDeadlineExceeded is error for WithDeadline and WithTimeout
	ErrDeadlineExceeded = errors.New("deadline exceeded")
)

// TODO context indicates that it is a placeholder and context is not yet
// used at all.
func TODO() Context {
	return todo
}

// Background is an empty context that cancels nothing.
func Background() Context {
	return background
}

// CacnelFunc is a signature of a cancel function contexts return. Different
// contexts have a bit different implementations of it.
type CancelFunc func()

// cancelCtx is the basis of other cancelling contexts.
type cancelCtx struct {
	Context
	done chan struct{}
	err  error
	mu   sync.Mutex
}

func (ctx *cancelCtx) Done() <-chan struct{} {
	// context aware methods will wait for Done channel to be closed.
	// Nothing ever sent to the channel, it only reacts on closing.
	return ctx.done
}

// Err implementation that is reused by all canceling contexts.
func (ctx *cancelCtx) Err() error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	return ctx.err
}

// WithCancel uses cancelCtx to generate canceling context and cancel
// function.
func WithCancel(parent Context) (Context, CancelFunc) {
	ctx := &cancelCtx{
		Context: parent,
		done:    make(chan struct{}),
	}

	cancel := func() {
		// ctx.cancel ingores canceling if there is already an error, if not
		// it creates a new error and closes the channel for Done()
		// it uses mutex to be thread-safe.
		ctx.cancel(ErrCanceled)
	}

	// this func ignores ctx.Done() as it should be picked by the "main" cancel
	// function, but if parent cancels itself, it generates closes Done() chan
	// and creates error from the parent. This allows to propagate context
	// again and again throuth the chain of methods/functions.
	go func() {
		select {
		case <-parent.Done():
			ctx.cancel(parent.Err())
		case <-ctx.Done():
		}
	}()
	return ctx, cancel
}

func (ctx *cancelCtx) cancel(err error) {
	// mutex makes it thread safe. I wonder what happens when ctx is sent to
	// a different computer. My understanding that because context can never
	// change parent, it starts a "new life" on a remote machine, and so it is
	// not a problem.
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	if ctx.err != nil {
		return
	}
	ctx.err = err
	close(ctx.done)
}

// deadlineCtx implements deadline and timeout
type deadlineCtx struct {
	*cancelCtx
	deadline time.Time
}

// Deadline ok means that deadline is set. All non-deadline context like
// Background, or simple cancel return false for ok.
func (dctx *deadlineCtx) Deadline() (deadline time.Time, ok bool) {
	return dctx.deadline, true
}

func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc) {
	// first we create with cancel context, because it would be indempotant for
	// in case parent is also cancel.
	cctx, cancel := WithCancel(parent)
	// then we wrap cancel context with deadline
	dctx := &deadlineCtx{
		cancelCtx: cctx.(*cancelCtx),
		// now deadline can be returned by Deadline method of deadlineCtx
		deadline: deadline,
	}

	// AfterFunc waits until it can run cancel funtion, and also returns a timer.
	t := time.AfterFunc(time.Until(deadline), func() {
		dctx.cancel(ErrDeadlineExceeded)
	})

	// we need to stop the timer to cleanup long-going methods with long
	// deadlines, if parent is canceled.
	stop := func() {
		t.Stop()
		cancel()
	}
	// now we wrap one cancel func into another cancel func (called stop here)
	return dctx, stop
}

// Timeout is a convenience function that reuses deadlineCtx
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	return WithDeadline(parent, time.Now().Add(timeout))
}

type valueCtx struct {
	Context
	key, value interface{}
}

func (vctx *valueCtx) Value(key interface{}) interface{} {
	// we try to find a key, if it is found, we return the value stored
	// in the context.
	if vctx.key == key {
		return vctx.value
	}
	return vctx.Context.Value(key)
}

// Here we set key and value. They should be for information only, and
// directly connected to request functions.
func WithValue(parent Context, key, value interface{}) Context {
	// key cannot be nil
	if key == nil {
		panic("key is nil")
	}
	// key must recognize ==, if not (like slices for example), abort.
	if !reflect.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}
	// Now context is wrapped with values.
	return &valueCtx{
		Context: parent,
		key:     key,
		value:   value,
	}
}
