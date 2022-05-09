package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"simple-jwt-go/api/auth"
	"simple-jwt-go/api/database"
	"simple-jwt-go/api/models"
	"simple-jwt-go/api/utils"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

// SignUp is function for create New User
// @Title Sign Up
// @Description Create a new user from JSON-formatted request body
// @Param post body auth.Credentials true "auth.Credential"
// @Success 201 object "Created - No Body"
// @Success 400 object models.ErrorResponse "ErrorResponse JSON"
// @Success 500 object models.ErrorResponse "Created - No Body"
// @Route /signup [post]

func SignUp(w http.ResponseWriter, r *http.Request) {

	// check body paylod json request
	var registrationReq models.RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&registrationReq); err != nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest, *(models.NewErrorResponse("invalid body format")), nil)
		return
	}

	// create user in db level application
	var user, dbUser models.User
	if err := registrationReq.CreateUser(&user); err != nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest, *(models.NewErrorResponse(err.Error())), nil)
		return
	}

	// hash password
	var err error
	user.Password, err = utils.HashPassword(user.Password)
	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse(err.Error())), nil)
		return
	}

	// connect db
	db, err := database.ConnectDB()
	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse(err.Error())), nil)
	}

	dbUser = models.User{}
	err = db.Where("username = ? OR email = ?", user.Username, user.Email).First(&dbUser).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse(err.Error())), nil)
		return
	} else if err == nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse(err.Error())), nil)
		return
	}

	// create user in db
	if err := db.Select("username", "email", "password").Create(&user).Error; err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse(err.Error())), nil)
		return
	}

	utils.JSONResponseWriter(&w, http.StatusCreated, nil, nil)
}

// SignIn is function for sign in also to get token JWT for creds/auth
// @Title Sign In
// @Description Sign in with JSON-formatted request body
// @Param body models.RegistrationRequest
// @Success 200 object auth.Claims "auth.Claims"
// @Success 400 object models.ErrorResponse "ErrorResponse JSON"
// @Success 500 object models.ErrorResponse "Created - No Body"
// @Route /signup

func SignIn(w http.ResponseWriter, r *http.Request) {
	var creds auth.Credentials
	var user models.User

	// check post body paylod json request
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest, nil, nil)
		return
	}

	// connect db
	db, err := database.ConnectDB()
	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse(err.Error())), nil)
		return
	}

	// check username in db
	if err := db.Where("username = ?", creds.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JSONResponseWriter(&w, http.StatusUnauthorized, *(models.NewErrorResponse("wrong password or username")), nil)
			return
		}

		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse(err.Error())), nil)
		return
	}

	passTrue := utils.CheckPassword(creds.Password, user.Password)

	if !passTrue {
		utils.JSONResponseWriter(&w, http.StatusUnauthorized, *(models.NewErrorResponse("wrong password or username")), nil)
		return
	}

	// build
	claims := &auth.Claims{
		ID:             user.ID,
		Username:       creds.Username,
		StandardClaims: jwt.StandardClaims{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse(err.Error())), nil)
		return
	}

	utils.JSONResponseWriter(&w, http.StatusOK, map[string]interface{}{"token": tokenString}, nil)
	return
}
