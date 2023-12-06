package main

import (
	"database/sql"

	"net/http"

	"github.com/qthuy2k1/product-management/internal/handlers"
	"github.com/redis/go-redis/v9"

	"github.com/qthuy2k1/product-management/internal/handlers/graph"
	"github.com/qthuy2k1/product-management/internal/handlers/rest"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/qthuy2k1/product-management/internal/controllers"
	"github.com/qthuy2k1/product-management/internal/repositories"

	graphHandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

// InitRoutes initializes a new router and sets up routes
func InitRoutes(db *sql.DB, redis *redis.Client) http.Handler {
	// create new router
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.MethodNotAllowed(handlers.MethodNotAllowedHandler)
	r.NotFound(handlers.NotFoundHandler)

	// create new handler
	repository := repositories.NewRepository(db, redis)
	controller := controllers.NewController(repository)

	initRest(r, controller)
	initGraph(r, controller)

	return r
}

// initRest initializes the REST API for the application
func initRest(r *chi.Mux, controller controllers.IController) {
	restHandler := rest.NewHandler(controller)

	//* user router
	r.Route("/users", func(r chi.Router) {
		r.Post("/", restHandler.CreateUser)
		r.Get("/{userID}", restHandler.GetUser)
	})

	//* product category router
	r.Route("/product-categories", func(r chi.Router) {
		r.Post("/", restHandler.CreateProductCategory)
	})

	//* product router
	r.Route("/products", func(r chi.Router) {
		r.Post("/", restHandler.CreateProduct)
		r.Get("/", restHandler.GetProducts)
		r.Route("/{productID}", func(r chi.Router) {
			r.Put("/", restHandler.UpdateProduct)
			r.Delete("/", restHandler.DeleteProduct)
		})
		r.Post("/import-csv", restHandler.ImportProductsFromCSV)
		r.Get("/export-csv", restHandler.ExportProductsToCSV)
	})
}

// initGraph initializes the GraphQL API for the application
func initGraph(r *chi.Mux, controller controllers.IController) {
	srv := graphHandler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Controller: controller}}))

	r.Handle("/graphql-playground", playground.Handler("GraphQL playground", "/graphql"))
	r.Handle("/graphql", srv)
}
