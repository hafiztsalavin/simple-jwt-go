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
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"gorm.io/gorm"
)

func SignUp(w http.ResponseWriter, r *http.Request) {

	// parse body paylod json request
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

	utils.JSONResponseWriter(&w, http.StatusCreated, map[string]interface{}{"message": "registration success"}, nil)
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	var creds auth.Credentials
	var user models.User

	// parse body paylod json request
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
	if err := db.Where("email = ?", creds.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JSONResponseWriter(&w, http.StatusUnauthorized, *(models.NewErrorResponse("wrong password or username")), nil)
			return
		}

		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse(err.Error())), nil)
		return
	}

	// check password
	passTrue := utils.CheckPassword(creds.Password, user.Password)

	if !passTrue {
		utils.JSONResponseWriter(&w, http.StatusUnauthorized, *(models.NewErrorResponse("wrong password or username")), nil)
		return
	}

	expirationTime := time.Now().Add(time.Minute * 5)

	// build
	claims := &auth.Claims{
		ID:       user.ID,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
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

func GetUser(w http.ResponseWriter, r *http.Request) {
	// get user id from context
	userID := context.Get(r, "id").(uint32)

	// connect db
	db, err := database.ConnectDB()
	if err != nil || db == nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse(err.Error())), nil)
		return
	}

	var dbUser models.User
	if err := db.Where("id = ?", userID).First(&dbUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JSONResponseWriter(&w, http.StatusNotFound,
				*(models.NewErrorResponse("can't find this user")), nil)
			return
		}

		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	var userRes models.UserResponse
	if err := userRes.InsertFromModel(dbUser); err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse(err.Error())), nil)
		return
	}

	utils.JSONResponseWriter(&w, http.StatusOK, userRes, nil)
	return
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := context.Get(r, "id").(uint32)

	// parse body paylod json request
	var userUpdateReq models.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&userUpdateReq); err != nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest, *(models.NewErrorResponse("invalid body format")), nil)
		return
	}

	var user, dbUser models.User
	user.ID = userID
	if err := userUpdateReq.UpdateUserModel(&user); err != nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest, *(models.NewErrorResponse(err.Error())), nil)
		return
	}

	// connect db
	db, err := database.ConnectDB()
	if err != nil || db == nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse(err.Error())), nil)
		return
	}

	// check id in db
	if err := db.Where("id = ?", user.ID).First(&dbUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JSONResponseWriter(&w, http.StatusNotFound, *(models.NewErrorResponse("can't find this user")), nil)
			return
		}
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse(err.Error())), nil)
		return
	}

	if userID != dbUser.ID {
		utils.JSONResponseWriter(&w, http.StatusForbidden, *(models.NewErrorResponse("can't do action at this user")), nil)
		return
	}

	dbUser = models.User{}
	err = db.Where("username = ? AND email = ?", user.Username, user.Email).First(&dbUser).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse(err.Error())), nil)
		return
	} else if err == nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest, *(models.NewErrorResponse("existing email or username")), nil)
		return
	}

	if user.Password != "" {
		if user.Password, err = utils.HashPassword(user.Password); err != nil {
			utils.JSONResponseWriter(&w, http.StatusForbidden, *(models.NewErrorResponse(err.Error())), nil)
			return
		}
	}

	if err := db.Model(&user).Updates(user).Error; err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse("username or email already exists")), nil)
		return
	}

	utils.JSONResponseWriter(&w, http.StatusNoContent, nil, nil)
	return
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := context.Get(r, "id").(uint32)

	// connect db
	db, err := database.ConnectDB()
	if err != nil || db == nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse(err.Error())), nil)
		return
	}

	var dbUser models.User
	if err := db.Where("id = ?", userID).First(&dbUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JSONResponseWriter(&w, http.StatusNotFound,
				*(models.NewErrorResponse("can't find this user")), nil)
			return
		}

		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	if err := db.Delete(&dbUser, userID).Error; err != nil {
		utils.JSONResponseWriter(&w, http.StatusForbidden, *(models.NewErrorResponse(err.Error())), nil)
		return
	}

	var userRes models.UserResponse
	userRes.InsertFromModel(dbUser)

	utils.JSONResponseWriter(&w, http.StatusOK, userRes, nil)
	return
}
