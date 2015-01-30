package goji

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/victoru/goji/web"
)

// The default web.Mux.
var DefaultMux *mux.Router

func init() {
	DefaultMux = mux.NewRouter()
}

// Wrap wraps a given handler with a list of given middlewares
func Wrap(hf interface{}, mfns ...interface{}) {
	var h web.Handler
	h = web.ParseHandler(hf)
	for _, mfn := range mfns {
		h = web.HandlerFunc(func(c web.C, w http.ResponseWriter, r *http.Request) {
			web.ParseMiddleware(mfn)(c, h)
		})
	}
}

// Handle adds a route to the default Mux. See the documentation for web.Mux for
// more information about what types this function accepts.
func Handle(path string, handler interface{}) *mux.Route {
	return DefaultMux.Handle(path, web.ParseHandler(handler))
}

// Connect adds a CONNECT route to the default Mux. See the documentation for
// web.Mux for more information about what types this function accepts.
func Connect(path string, handler interface{}) *mux.Route {
	return DefaultMux.Handle(path, web.ParseHandler(handler)).Methods("CONNECT")
}

// Delete adds a DELETE route to the default Mux. See the documentation for
// web.Mux for more information about what types this function accepts.
func Delete(path string, handler interface{}) *mux.Route {
	return DefaultMux.Handle(path, web.ParseHandler(handler)).Methods("DELETE")
}

// Get adds a GET route to the default Mux. See the documentation for web.Mux for
// more information about what types this function accepts.
func Get(path string, handler interface{}) *mux.Route {
	return DefaultMux.Handle(path, web.ParseHandler(handler)).Methods("GET")
}

// Head adds a HEAD route to the default Mux. See the documentation for web.Mux
// for more information about what types this function accepts.
func Head(path string, handler interface{}) *mux.Route {
	return DefaultMux.Handle(path, web.ParseHandler(handler)).Methods("HEAD")
}

// Options adds a OPTIONS route to the default Mux. See the documentation for
// web.Mux for more information about what types this function accepts.
func Options(path string, handler interface{}) *mux.Route {
	return DefaultMux.Handle(path, web.ParseHandler(handler)).Methods("HEAD")
}

// Patch adds a PATCH route to the default Mux. See the documentation for web.Mux
// for more information about what types this function accepts.
func Patch(path string, handler interface{}) *mux.Route {
	return DefaultMux.Handle(path, web.ParseHandler(handler)).Methods("HEAD")
}

// Post adds a POST route to the default Mux. See the documentation for web.Mux
// for more information about what types this function accepts.
func Post(path string, handler interface{}) *mux.Route {
	return DefaultMux.Handle(path, web.ParseHandler(handler)).Methods("HEAD")
}

// Put adds a PUT route to the default Mux. See the documentation for web.Mux for
// more information about what types this function accepts.
func Put(path string, handler interface{}) *mux.Route {
	return DefaultMux.Handle(path, web.ParseHandler(handler)).Methods("HEAD")
}

// Trace adds a TRACE route to the default Mux. See the documentation for
// web.Mux for more information about what types this function accepts.
func Trace(path string, handler interface{}) *mux.Route {
	return DefaultMux.Handle(path, web.ParseHandler(handler)).Methods("HEAD")
}

// NotFound sets the NotFound handler for the default Mux. See the documentation
// for web.Mux.NotFound for more information.
func NotFound(handler interface{}) {
	DefaultMux.NotFoundHandler = web.ParseHandler(handler)
}
