package main

import (
	"log"
	"net/http"

	"github.com/FRanggaY/personal-portfolio-api/docs"
	"github.com/FRanggaY/personal-portfolio-api/handlers"
	"github.com/FRanggaY/personal-portfolio-api/middlewares"
	"github.com/FRanggaY/personal-portfolio-api/models"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	r := mux.NewRouter()

	models.ConnectDatabase()

	var basePathRoute = "/api/v1"

	docs.SwaggerInfo.Title = "Swagger Personal Portfolio API"
	docs.SwaggerInfo.Description = "This is a sample personal portfolio server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = basePathRoute
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	r.PathPrefix("/documentation/").Handler(httpSwagger.WrapHandler)

	// static
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	api := r.PathPrefix(basePathRoute).Subrouter()
	api.HandleFunc("/login", handlers.Login).Methods("POST")
	api.HandleFunc("/register", handlers.Register).Methods("POST")
	api.HandleFunc("/logout", handlers.Logout).Methods("GET")

	apiProtect := r.PathPrefix(basePathRoute).Subrouter()
	apiProtect.HandleFunc("/profile", handlers.Profile).Methods("GET")
	apiProtect.HandleFunc("/user", handlers.GetFilteredPaginatedUsers).Methods("GET")

	var userDetailRoute = "/user/{id}"
	apiProtect.HandleFunc(userDetailRoute, handlers.GetUser).Methods("GET")
	apiProtect.HandleFunc(userDetailRoute, handlers.UpdateUser).Methods("PUT")
	apiProtect.HandleFunc(userDetailRoute, handlers.DeleteUser).Methods("DELETE")

	apiProtect.HandleFunc("/school", handlers.GetFilteredPaginatedSchools).Methods("GET")
	apiProtect.HandleFunc("/school", handlers.CreateSchool).Methods("POST")
	apiProtect.HandleFunc("/school/{id}", handlers.ReadSchool).Methods("GET")
	apiProtect.Use(middlewares.JWTMiddleware)

	log.Fatal(http.ListenAndServe(":8080", r))
}
