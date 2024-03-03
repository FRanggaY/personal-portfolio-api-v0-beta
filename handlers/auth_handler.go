package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/FRanggaY/personal-portfolio-api/config"
	"github.com/FRanggaY/personal-portfolio-api/helper"
	"github.com/FRanggaY/personal-portfolio-api/models"
	"github.com/FRanggaY/personal-portfolio-api/repositories"
	"github.com/golang-jwt/jwt/v5"
)

// login godoc
// @Summary login user
// @Description login user with the provided details
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.UserLoginForm true "User input"
// @Success 200 {object} map[string]string "Success"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	// define input from json
	var userInput models.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	// validation
	var user models.User
	userRepo := repositories.NewUserRepository()
	exist_user, err := userRepo.ReadByUsername(userInput.Username)
	if err != nil {
		// Handle error
		response := map[string]string{"message": "Username or password invalid"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// check password valid or not
	if err := userRepo.CompareUserPassword(exist_user.Password, userInput.Password); err != nil {
		response := map[string]string{"message": "Username or password invalid"}
		helper.ResponseJSON(w, http.StatusUnauthorized, response)
		return
	}

	// create token jwt (expired 24 hour)
	expTime := time.Now().Add(24 * time.Hour)
	claims := &config.JWTClaim{
		Id:       exist_user.Id,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-jwt-mux",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	// declare algorithm for sign
	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// signed token
	token, err := tokenAlgo.SignedString(config.JWT_KEY)
	if err != nil {
		response := map[string]string{"message": err.Error()}
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	// set token yang ke cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    token,
		HttpOnly: true,
	})

	response := map[string]string{
		"message":      "success",
		"access_token": token,
		"expired":      expTime.String(),
	}
	helper.ResponseJSON(w, http.StatusOK, response)
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with the provided details
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.UserCreateForm true "User input"
// @Success 201 {object} map[string]string "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Router /register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	// define input from json
	var userInput models.UserCreateForm
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		log.Fatal("Error decoding new user: ")
	}
	defer r.Body.Close()

	// validate username unique
	userRepo := repositories.NewUserRepository()
	exist_user, _ := userRepo.ReadByUsername(userInput.Username)
	if exist_user != nil {
		// Handle error
		response := map[string]string{"message": "Username already used"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	// insert to database
	if newUser, err := userRepo.Create(&userInput); err != nil {
		// Handle error
		response := map[string]string{"message": "Error creating new user"}
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	} else {
		response := map[string]interface{}{
			"message": "success",
			"data": map[string]interface{}{
				"id": newUser.Id,
			},
		}
		helper.ResponseJSON(w, http.StatusOK, response)
		return
	}
}

// Logout godoc
// @Summary Logout user
// @Description Logout user
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Success"
// @Router /logout [get]
func Logout(w http.ResponseWriter, r *http.Request) {
	// remove token from cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})

	response := map[string]string{"message": "logout success"}
	helper.ResponseJSON(w, http.StatusOK, response)
}

// Profile godoc
// @Summary Profile user
// @Description Profile user
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Success"
// @Router /profile [get]
func Profile(w http.ResponseWriter, r *http.Request) {
	jwtClaim, _ := helper.GetJWTClaim(r)

	userRepo := repositories.NewUserRepository()
	user, err := userRepo.Read(jwtClaim.Id)
	if err != nil {
		response := map[string]string{"message": "User ID not found"}
		helper.ResponseJSON(w, http.StatusNotFound, response)
		return
	}

	response := map[string]interface{}{
		"message": "success",
		"data": map[string]interface{}{
			"id":       user.Id,
			"username": user.Username,
			"name":     user.Name,
		},
	}

	helper.ResponseJSON(w, http.StatusOK, response)
}
