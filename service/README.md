# Ideas from Mat Ryer about organizing web-services

## Code principles:

1. Maintainability

2. Glancability

3. Code has to be boring

4. Consistency. Make code slightly more complex it it brings it to "boring"
   pattern used everywhere else.



## Use of constructors sparingly, direct init is explains intent better.

```go
func newService() *server {
  s := &server{}
  s.routes()

  return s
}
```

## Make a server an http handler:

* Implement `serveHTTP` interface that makes your service a handler.

```go
func (s *server) serveHTTP(w http.ResponseWriter, r *http.Request) {
  s.router.ServeHTTP()
}
```

## Place all routes in one place

```go
func (s *server) routes() {
  s.router.Get("/api", s.handleApi())
  s.router.Get("/", s.handleHome())
  s.router.Get("/doc", s.handleDoc())
}
```

## Handler functions

```go
func(s *server) handleApi() http.HandlerFunc{
  ...
}
```

## Middleware are just go functions

```go
func(s *server) adminOnly(h http.HandlerFunc) http.HandlerFunc{
  return func(w http.ResponseWriter, r http.Request) {
    if !currentUser().isAdmin {
      http.NotFound(w, r)
      return
    }
    h(w, r)
  }
}
```

## To reduce "magic" call middleware explicitly in routes

```go
func (s *server) routes() {
  s.router.Get("/api", s.handleApi())
  s.router.Get("/", s.handleHome())
  s.router.Get("/doc", s.handleDoc())
  s.router.Get("/admin", s.isAdmin(s.handle.AdminHome()))
}
```

## For expensive things use `sync.Once`

```go
// lazy init
func (s *server) handleTemplate(files string...) http.HandlerFunc {
  var init sync.Once
  return func(w http.ResponseWriter, r http.Request) {
    init.Do(expensiveThings)
    ...
  }
}
```

## Server made this way is testable

Create a new server for each test.

Only set dependencies that you need inside the test.

```go
func TestHandleAbout(t *testing.T) {
  is := is.New(t)
  srv := newServer()
  r := httptest.NewRequest("GET", "/about", nil)
  w := httptest.NewRecorder()
  srv.ServeHTTP(w, r)
  is.Equal(w.StatusCode, http.StatusOK)
}

```

You can test the whole stack: ``srv.ServeHTTP(w,r)``.

You can also test the handler only: ``srv.handleAbout()
