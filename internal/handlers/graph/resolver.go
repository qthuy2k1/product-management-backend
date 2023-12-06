package graph

import "github.com/qthuy2k1/product-management/internal/controllers"

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Controller controllers.IController
}
