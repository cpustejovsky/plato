package router

import (
	"fmt"
	"net/http"

	"github.com/cpustejovsky/plato/middleware"
)

// GET takes a url and prepends "GET" to it to specify the allow method
func GET(url string) string {
	return "GET " + url
}

// POST takes a url and prepends "POST" to it to specify the allow method
func POST(url string) string {
	return "POST " + url
}

// PUT takes a url and prepends "PUT" to it to specify the allow method
func PUT(url string) string {
	return "PUT " + url
}

// PATCH takes a url and prepends "PATCH" to it to specify the allow method
func PATCH(url string) string {
	return "PATCH " + url
}

// DELETE takes a url and prepends "DELETE" to it to specify the allow method
func DELETE(url string) string {
	return "DELETE " + url
}

type Router struct {
	*http.ServeMux
	middlewares []middleware.MiddlewareFunc
}

func New() (*Router, error) {
	mux := http.NewServeMux()
	mux.Handle(POST("/foo"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "created foo") }))
	mux.Handle(GET("/foo"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "got foos") }))
	mux.Handle(GET("/foo/{id}"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "got foo with id: %v", r.PathValue("id"))
	}))
	mux.Handle(PUT("/foo/{id}"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "updated foo with id: %v", r.PathValue("id"))
	}))

	mux.Handle(PATCH("/foo/{id}"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "edited foo with id: %v", r.PathValue("id"))
	}))
	mux.Handle(DELETE("/foo/{id}"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "deleted foo with id: %v", r.PathValue("id"))
	}))

	var mws []middleware.MiddlewareFunc

	mws = append(mws, middleware.LogRequest)
	mws = append(mws, middleware.DrainAndClose)

	r := Router{
		ServeMux:    mux,
		middlewares: mws,
	}
	return &r, nil
}

func (r *Router) ServeHTTP(w http.ResponseWriter, rq *http.Request) {
	var h http.Handler = r.ServeMux

	for _, mw := range r.middlewares {
		h = mw(h)
	}
	h.ServeHTTP(w, rq)
}
