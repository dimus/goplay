package contextimpl

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTodoBackground(t *testing.T) {
	todo := fmt.Sprint(TODO())
	bg := fmt.Sprint(Background())
	assert.NotEqualf(t, todo, bg, "%q and %q shoud not be equal", todo, bg)
}

func TestWithCancel(t *testing.T) {
	ctx, cancel := WithCancel(Background())
	assert.Nilf(t, ctx.Err(), "error should be nil at first %v", ctx.Err())
	cancel()
	<-ctx.Done()
	assert.Equal(t, ctx.Err(), ErrCanceled, "error should not be nil after cancel %v", ctx.Err())
}

func TestWithCancelConcurrent(t *testing.T) {
	ctx, cancel := WithCancel(Background())

	time.AfterFunc(1*time.Second, cancel)

	assert.Nilf(t, ctx.Err(), "error should be nil first, got %v", ctx.Err())
	<-ctx.Done()
	assert.Equalf(t, ctx.Err(), ErrCanceled, "error should be Canceled, got %v", ctx.Err())
}

func TestWithCancelPropagation(t *testing.T) {
	ctxA, cancelA := WithCancel(Background())
	ctxB, cancelB := WithCancel(ctxA)
	defer cancelB()

	cancelA()

	select {
	case <-ctxB.Done():
	case <-time.After(1 * time.Second):
		t.Errorf("time out")
	}

	assert.Equalf(t, ctxB.Err(), ErrCanceled, "error should be canceled now, got %v", ctxB.Err())
}

func TestWithDeadline(t *testing.T) {
	deadline := time.Now().Add(2 * time.Second)
	ctx, cancel := WithDeadline(Background(), deadline)

	d, ok := ctx.Deadline()
	msg := fmt.Sprintf("expected deadline %v; got %v", deadline, d)
	assert.True(t, ok, msg)
	assert.Equal(t, d, deadline, msg)

	then := time.Now()
	<-ctx.Done()
	if d := time.Since(then); math.Abs(d.Seconds()-2.0) > 0.1 {
		t.Errorf("should have been done after 2.0 seconds, took %v", d)
	}
	assert.Equalf(t, ctx.Err(), ErrDeadlineExceeded,
		"error should be DeadlineExceeded, got %v", ctx.Err())

	cancel()
	assert.Equalf(t, ctx.Err(), ErrDeadlineExceeded,
		"error should still be DeadlineExceeded, got %v", ctx.Err())
}

func TestWithTimeout(t *testing.T) {
	timeout := 2 * time.Second
	deadline := time.Now().Add(timeout)
	ctx, cancel := WithTimeout(Background(), timeout)

	if d, ok := ctx.Deadline(); !ok || d.Sub(deadline) > time.Millisecond {
		t.Errorf("expected deadline %v; got %v", deadline, d)
	}

	then := time.Now()
	<-ctx.Done()
	if d := time.Since(then); math.Abs(d.Seconds()-2.0) > 0.1 {
		t.Errorf("should have been done after 2.0 seconds, took %v", d)
	}
	if err := ctx.Err(); err != ErrDeadlineExceeded {
		t.Errorf("error should be DeadlineExceeded, got %v", err)
	}

	cancel()
	if err := ctx.Err(); err != ErrDeadlineExceeded {
		t.Errorf("error should still be DeadlineExceeded, got %v", err)
	}
}

func TestWithValue(t *testing.T) {
	tc := []struct {
		key, val, keyRet, valRet interface{}
		shouldPanic              bool
	}{
		{"a", "b", "a", "b", false},
		{"a", "b", "c", nil, false},
		{42, true, 42, true, false},
		{42, true, int64(42), nil, false},
		{nil, true, nil, nil, true},
		{[]int{1, 2, 3}, true, []int{1, 2, 3}, nil, true},
	}

	for _, tt := range tc {
		var panicked interface{}
		func() {
			defer func() { panicked = recover() }()

			ctx := WithValue(Background(), tt.key, tt.val)
			if val := ctx.Value(tt.keyRet); val != tt.valRet {
				t.Errorf("expected value %v, got %v", tt.valRet, val)
			}
		}()

		if panicked != nil && !tt.shouldPanic {
			t.Errorf("unexpected panic: %v", panicked)
		}
		if panicked == nil && tt.shouldPanic {
			t.Errorf("expected panic, but didn't get it")
		}
	}
}
