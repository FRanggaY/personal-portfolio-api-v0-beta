package main

import (
	"log"
	"net/http"

	"github.com/FRanggaY/personal-portfolio-api/docs"
	"github.com/FRanggaY/personal-portfolio-api/handlers/authHandler"
	"github.com/FRanggaY/personal-portfolio-api/models"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	r := mux.NewRouter()

	models.ConnectDatabase()

	docs.SwaggerInfo.Title = "Swagger Personal Portfolio API"
	docs.SwaggerInfo.Description = "This is a sample personal portfolio server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	r.PathPrefix("/documentation/").Handler(httpSwagger.WrapHandler)

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/v1/login", authHandler.Login).Methods("POST")
	api.HandleFunc("/v1/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/v1/logout", authHandler.Logout).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
