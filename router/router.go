package router

import (
	"go-movies-crud/controller"
	"go-movies-crud/middleware"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Router() *mux.Router {
	r := mux.NewRouter()

	// Apply the TimeRequest middleware to all routes
	r.Use(middleware.TimeRequest)

	// API endpoints
	r.HandleFunc("/movies", controller.GetMovies).Methods("GET", "OPTIONS")
	r.HandleFunc("/movies/{id}", controller.GetMovie).Methods("GET", "OPTIONS")
	r.HandleFunc("/movies", controller.CreateMovies).Methods("POST", "OPTIONS")
	r.HandleFunc("/movies/{id}", controller.UpdateMovie).Methods("PUT", "OPTIONS")
	r.HandleFunc("/movies/{id}", controller.DeleteMovie).Methods("DELETE", "OPTIONS")

	// Serve Swagger UI and Swagger JSON
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // Adjust this URL as per your setup
	))

	return r
}
