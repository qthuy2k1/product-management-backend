package rest

import (
	"github.com/qthuy2k1/product-management/internal/controllers"
)

type Handler struct {
	Controller controllers.IController
}

// NewHandler returns the Handler struct that contains the Controller
func NewHandler(controller controllers.IController) *Handler {
	return &Handler{Controller: controller}
}
