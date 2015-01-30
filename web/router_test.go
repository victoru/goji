package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// These tests can probably be DRY'd up a bunch

func chHandler(ch chan string, s string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ch <- s
	})
}

var methods = []string{"CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "PATCH",
	"POST", "PUT", "TRACE", "OTHER"}

func TestMethods(t *testing.T) {
	t.Parallel()
	m := New()
	ch := make(chan string, 1)

	m.Connect("/", chHandler(ch, "CONNECT"))
	m.Delete("/", chHandler(ch, "DELETE"))
	m.Head("/", chHandler(ch, "HEAD"))
	m.Get("/", chHandler(ch, "GET"))
	m.Options("/", chHandler(ch, "OPTIONS"))
	m.Patch("/", chHandler(ch, "PATCH"))
	m.Post("/", chHandler(ch, "POST"))
	m.Put("/", chHandler(ch, "PUT"))
	m.Trace("/", chHandler(ch, "TRACE"))
	m.Handle("/", chHandler(ch, "OTHER"))

	for _, method := range methods {
		r, _ := http.NewRequest(method, "/", nil)
		w := httptest.NewRecorder()
		m.ServeHTTP(w, r)
		select {
		case val := <-ch:
			if val != method {
				t.Errorf("Got %q, expected %q", val, method)
			}
		case <-time.After(5 * time.Millisecond):
			t.Errorf("Timeout waiting for method %q", method)
		}
	}
}

type testPattern struct{}

func (t testPattern) Prefix() string {
	return ""
}

func (t testPattern) Match(r *http.Request, c *C) bool {
	return true
}
func (t testPattern) Run(r *http.Request, c *C) {
}

type testHandler chan string

func (t testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t <- "http"
}
func (t testHandler) ServeHTTPC(c C, w http.ResponseWriter, r *http.Request) {
	t <- "httpc"
}

var testHandlerTable = map[string]string{
	"/a": "http fn",
	"/b": "http handler",
	"/c": "web fn",
	"/d": "web handler",
	"/e": "httpc",
}

func TestHandlerTypes(t *testing.T) {
	t.Parallel()
	m := New()
	ch := make(chan string, 1)

	m.Get("/a", func(w http.ResponseWriter, r *http.Request) {
		ch <- "http fn"
	})
	m.Get("/b", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ch <- "http handler"
	}))
	m.Get("/c", func(c C, w http.ResponseWriter, r *http.Request) {
		ch <- "web fn"
	})
	m.Get("/d", HandlerFunc(func(c C, w http.ResponseWriter, r *http.Request) {
		ch <- "web handler"
	}))
	m.Get("/e", testHandler(ch))

	for route, response := range testHandlerTable {
		r, _ := http.NewRequest("GET", route, nil)
		w := httptest.NewRecorder()
		m.ServeHTTP(w, r)
		select {
		case resp := <-ch:
			if resp != response {
				t.Errorf("Got %q, expected %q", resp, response)
			}
		case <-time.After(5 * time.Millisecond):
			t.Errorf("Timeout waiting for path %q", route)
		}

	}
}

func TestNotFound(t *testing.T) {
	t.Parallel()
	m := New()

	r, _ := http.NewRequest("post", "/", nil)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, r)
	if w.Code != 404 {
		t.Errorf("Expected 404, got %d", w.Code)
	}

	m.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "I'm a teapot!", http.StatusTeapot)
	})

	r, _ = http.NewRequest("POST", "/", nil)
	w = httptest.NewRecorder()
	m.ServeHTTP(w, r)
	if w.Code != http.StatusTeapot {
		t.Errorf("Expected a teapot, got %d", w.Code)
	}
}

func TestPrefix(t *testing.T) {
	t.Parallel()
	m := New()
	ch := make(chan string, 1)

	m.Handle("/hello/world", func(w http.ResponseWriter, r *http.Request) {
		ch <- r.URL.Path
	})

	r, _ := http.NewRequest("GET", "/hello/world", nil)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, r)
	select {
	case val := <-ch:
		if val != "/hello/world" {
			t.Errorf("Got %q, expected /hello/world", val)
		}
	case <-time.After(5 * time.Millisecond):
		t.Errorf("Timeout waiting for hello")
	}
}

//var validMethodsTable = map[string][]string{
//"/hello/carl":       {"DELETE", "GET", "HEAD", "PATCH", "POST", "PUT"},
//"/hello/bob":        {"DELETE", "GET", "HEAD", "PATCH", "PUT"},
//"/hola/carl":        {"DELETE", "GET", "HEAD", "PUT"},
//"/hola/bob":         {"DELETE"},
//"/does/not/compute": {},
//}

//func TestValidMethods(t *testing.T) {
//t.Parallel()
//m := New()
//ch := make(chan []string, 1)

//m.NotFound(func(c C, w http.ResponseWriter, r *http.Request) {
//if c.Env == nil {
//ch <- []string{}
//return
//}
//methods, ok := c.Env[ValidMethodsKey]
//if !ok {
//ch <- []string{}
//return
//}
//ch <- methods.([]string)
//})

//m.Get("/hello/carl", http.NotFound)
//m.Post("/hello/carl", http.NotFound)
//m.Head("/hello/bob", http.NotFound)
//m.Get("/hello/:name", http.NotFound)
//m.Put("/hello/:name", http.NotFound)
//m.Patch("/hello/:name", http.NotFound)
//m.Get("/:greet/carl", http.NotFound)
//m.Put("/:greet/carl", http.NotFound)
//m.Delete("/:greet/:anyone", http.NotFound)

//for path, eMethods := range validMethodsTable {
//r, _ := http.NewRequest("BOGUS", path, nil)
//m.ServeHTTP(httptest.NewRecorder(), r)
//aMethods := <-ch
//if !reflect.DeepEqual(eMethods, aMethods) {
//t.Errorf("For %q, expected %v, got %v", path, eMethods,
//aMethods)
//}
//}

//// This should also work when c.Env has already been initalized
//m.Use(func(c *C, h http.Handler) http.Handler {
//return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//c.Env = make(map[string]interface{})
//h.ServeHTTP(w, r)
//})
//})
//for path, eMethods := range validMethodsTable {
//r, _ := http.NewRequest("BOGUS", path, nil)
//m.ServeHTTP(httptest.NewRecorder(), r)
//aMethods := <-ch
//if !reflect.DeepEqual(eMethods, aMethods) {
//t.Errorf("For %q, expected %v, got %v", path, eMethods,
//aMethods)
//}
//}
//}
