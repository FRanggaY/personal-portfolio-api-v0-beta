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

	apiProtect.HandleFunc("/user-skill", handlers.CreateUserSkill).Methods("POST")
	apiProtect.HandleFunc("/user-skill/{id}", handlers.DeleteUserSkill).Methods("DELETE")

	apiProtect.HandleFunc("/user-experience", handlers.CreateUserExperience).Methods("POST")
	apiProtect.HandleFunc("/user-experience/{id}", handlers.DeleteUserExperience).Methods("DELETE")

	apiProtect.HandleFunc("/user-position", handlers.CreateUserPosition).Methods("POST")
	apiProtect.HandleFunc("/user-position/{id}", handlers.DeleteUserPosition).Methods("DELETE")

	apiProtect.HandleFunc("/user-education", handlers.CreateUserEducation).Methods("POST")
	apiProtect.HandleFunc("/user-education/{id}", handlers.DeleteUserEducation).Methods("DELETE")

	apiProtect.HandleFunc("/user-experience-translation", handlers.CreateUserExperienceTranslation).Methods("POST")
	apiProtect.HandleFunc("/user-experience-translation/{id}", handlers.DeleteUserExperienceTranslation).Methods("DELETE")

	apiProtect.HandleFunc("/user-education-translation", handlers.CreateUserEducationTranslation).Methods("POST")
	apiProtect.HandleFunc("/user-education-translation/{id}", handlers.DeleteUserEducationTranslation).Methods("DELETE")

	apiProtect.HandleFunc("/user-attachment", handlers.CreateUserAttachment).Methods("POST")

	apiProtect.HandleFunc("/school", handlers.GetFilteredPaginatedSchools).Methods("GET")
	apiProtect.HandleFunc("/school", handlers.CreateSchool).Methods("POST")
	apiProtect.HandleFunc("/school/{id}", handlers.ReadSchool).Methods("GET")

	apiProtect.HandleFunc("/company", handlers.GetFilteredPaginatedCompanies).Methods("GET")
	apiProtect.HandleFunc("/company", handlers.CreateCompany).Methods("POST")
	apiProtect.HandleFunc("/company/{id}", handlers.ReadCompany).Methods("GET")

	apiProtect.HandleFunc("/language", handlers.GetFilteredPaginatedLanguages).Methods("GET")
	apiProtect.HandleFunc("/language", handlers.CreateLanguage).Methods("POST")
	apiProtect.HandleFunc("/language/{id}", handlers.ReadLanguage).Methods("GET")

	apiProtect.HandleFunc("/skill", handlers.GetFilteredPaginatedSkills).Methods("GET")
	apiProtect.HandleFunc("/skill", handlers.CreateSkill).Methods("POST")
	apiProtect.HandleFunc("/skill/{id}", handlers.ReadSkill).Methods("GET")

	apiProtect.HandleFunc("/skill-translation", handlers.CreateSkillTranslation).Methods("POST")
	apiProtect.HandleFunc("/skill-translation/{id}", handlers.DeleteSkillTranslation).Methods("DELETE")
	apiProtect.Use(middlewares.JWTMiddleware)

	log.Fatal(http.ListenAndServe(":8080", r))
}
