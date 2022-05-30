package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"simple-jwt-go/api/auth"
	"simple-jwt-go/api/database"
	"simple-jwt-go/api/models"
	"simple-jwt-go/api/utils"
	"time"

	"gorm.io/gorm"
)

func SignUp(w http.ResponseWriter, r *http.Request) {

	// get json object from request body
	var registrationReq models.RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&registrationReq); err != nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest, *(models.NewErrorResponse("invalid body format")), nil)
		return
	}

	// create db entity level application
	var user, dbUser models.User
	if err := registrationReq.CreateUserModel(&user); err != nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest, *(models.NewErrorResponse(err.Error())), nil)
		return
	}

	// hashing password
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
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse("username or email already used, pls use another email or username")), nil)
		return
	}

	// insert user into db
	if err := db.Select("username", "email", "password").Create(&user).Error; err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse(err.Error())), nil)
		return
	}

	// respond with created user
	utils.JSONResponseWriter(&w, http.StatusCreated, map[string]interface{}{"message": "registration success"}, nil)
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	authorizationCookie, err := r.Cookie("access_token")
	if authorizationCookie != nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest, *(models.NewErrorResponse("Already logged in")), nil)
		return
	}

	var creds auth.Credentials
	var user models.User

	// get json object from request body
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

	// check username in db by email
	if err := db.Where("email = ?", creds.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JSONResponseWriter(&w, http.StatusUnauthorized, *(models.NewErrorResponse("wrong password or username")), nil)
			return
		}

		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse(err.Error())), nil)
		return
	}

	// check password and comparing
	passTrue := utils.CheckPassword(creds.Password, user.Password)
	if !passTrue {
		utils.JSONResponseWriter(&w, http.StatusUnauthorized, *(models.NewErrorResponse("wrong password or username")), nil)
		return
	}

	// generate access token and refresh token
	expirationAccessToken := time.Now().Add(time.Minute * 1)
	accessToken, err := auth.CreateJWTToken(user.ID, user.Email, expirationAccessToken)
	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse(err.Error())), nil)
		return
	}
	utils.AddCookie(w, "access_token", accessToken, expirationAccessToken)

	expirationRefreshToken := time.Now().Add(time.Minute * 5)
	refreshToken, err := auth.CreateRefreshToken(accessToken, expirationRefreshToken)
	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse(err.Error())), nil)
		return
	}
	utils.AddCookie(w, "refresh_token", refreshToken, expirationRefreshToken)

	utils.JSONResponseWriter(&w, http.StatusOK, map[string]interface{}{"message": "login success"}, nil)
	return
}
