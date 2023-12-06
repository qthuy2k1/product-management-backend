package handlers

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/qthuy2k1/product-management/internal/handlers/rest"
)

// MethodNotAllowedHandler is a middleware that writes the 405 status code to header and returns the error: "method not allowed"
func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(405)
	render.Render(w, r, rest.ErrMethodNotAllowed)
}

// NotFoundHandler is a middleware that writes the 400 status code to header and returns the error: "resource not found"
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(400)
	render.Render(w, r, rest.ErrNotFound)
}
