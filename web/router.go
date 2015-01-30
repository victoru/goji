package web

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type netHTTPWrap struct {
	http.Handler
}

func (h netHTTPWrap) ServeHTTPC(c C, w http.ResponseWriter, r *http.Request) {
	h.Handler.ServeHTTP(w, r)
}

func New() *Router {
	return &Router{mux.NewRouter()}
}

type Router struct {
	*mux.Router
}

func (r *Router) NotFound(h interface{}) {
	r.Router.NotFoundHandler = ParseHandler(h)

}

// Wrap wraps a handler with multiple middlewares
func Wrap(hf interface{}, mf ...interface{}) Handler {
	return wrap(hf, mf...)
}

func wrap(ih interface{}, middlewares ...interface{}) Handler {
	var h = ParseHandler(ih)
	return HandlerFunc(func(c C, w http.ResponseWriter, r *http.Request) {
		// new handler instance to prevent duplicate wraps
		var h Handler = h
		for i := len(middlewares); i > 0; i-- {
			mw := ParseMiddleware(middlewares[i-1])
			h = mw(c, h)
		}

		h.ServeHTTPC(c, w, r)
	})
}

func (r *Router) Handle(path string, h interface{}, m ...interface{}) *mux.Route {
	return r.Router.Handle(path, wrap(h, m...))
}

func (r *Router) HandleFunc(
	path string,
	f func(C, http.ResponseWriter, *http.Request),
	middlewares ...interface{},
) *mux.Route {
	return r.Router.Handle(path, wrap(f, middlewares...))
}

func ParseHandler(h interface{}) Handler {
	switch f := h.(type) {
	case Handler:
		return f
	case http.Handler:
		return netHTTPWrap{f}
	case func(c C, w http.ResponseWriter, r *http.Request):
		return HandlerFunc(f)
	case func(w http.ResponseWriter, r *http.Request):
		return netHTTPWrap{http.HandlerFunc(f)}
	default:
		log.Fatalf("Unknown handler type %v. Expected a web.Handler, "+
			"a http.Handler, or a function with signature func(C, "+
			"http.ResponseWriter, *http.Request) or "+
			"func(http.ResponseWriter, *http.Request)", h)
	}
	panic("log.Fatalf does not return")
}

func ParseMiddlewares(m ...interface{}) []func(C, Handler) Handler {
	var middlewares []func(C, Handler) Handler
	for _, m := range m {
		middlewares = append(middlewares, ParseMiddleware(m))
	}
	return middlewares
}

func ParseMiddleware(m interface{}) func(C, Handler) Handler {
	switch f := m.(type) {
	case func(Handler) Handler:
		return func(c C, h Handler) Handler {
			return f(h)
		}
	case func(C, Handler) Handler:
		return f
	default:
		log.Fatalf(`Unknown middleware type %#v. Expected a function `+
			`with signature "func(web.Handler) web.Handler" or `+
			`"func(*web.C, web.Handler) web.Handler".`, m)
	}
	panic("log.Fatalf does not return")
}

// helper methods
func (r *Router) Get(path string, h interface{}) *mux.Route {
	return r.Handle(path, h).Methods("GET")
}
func (r *Router) Post(path string, h interface{}) *mux.Route {
	return r.Handle(path, h).Methods("POST")
}
func (r *Router) Patch(path string, h interface{}) *mux.Route {
	return r.Handle(path, h).Methods("PATCH")
}
func (r *Router) Put(path string, h interface{}) *mux.Route {
	return r.Handle(path, h).Methods("PUT")
}
func (r *Router) Delete(path string, h interface{}) *mux.Route {
	return r.Handle(path, h).Methods("DELETE")
}
func (r *Router) Trace(path string, h interface{}) *mux.Route {
	return r.Handle(path, h).Methods("TRACE")
}
func (r *Router) Connect(path string, h interface{}) *mux.Route {
	return r.Handle(path, h).Methods("CONNECT")
}
func (r *Router) Options(path string, h interface{}) *mux.Route {
	return r.Handle(path, h).Methods("Options")
}
func (r *Router) Head(path string, h interface{}) *mux.Route {
	return r.Handle(path, h).Methods("HEAD")
}
