package router

import (
	"go-movies-crud/middleware"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Router() *mux.Router {
	r := mux.NewRouter()

	// API endpoints
	r.HandleFunc("/movies", middleware.GetMovies).Methods("GET", "OPTIONS")
	r.HandleFunc("/movies/{id}", middleware.GetMovie).Methods("GET", "OPTIONS")
	r.HandleFunc("/movies", middleware.CreateMovies).Methods("POST", "OPTIONS")
	r.HandleFunc("/movies/{id}", middleware.UpdateMovie).Methods("PUT", "OPTIONS")
	r.HandleFunc("/movies/{id}", middleware.DeleteMovie).Methods("DELETE", "OPTIONS")

	// Serve Swagger UI and Swagger JSON
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // Adjust this URL as per your setup
	))

	return r
}
