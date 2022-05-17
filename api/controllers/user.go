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
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
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
	// reference of
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

	// expirationTime := time.Now().Add(time.Minute * 5)
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

func GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	userID := context.Get(r, "id").(uint32)

	if idStr == "" || !utils.IsInteger(idStr) {
		utils.JSONResponseWriter(&w, http.StatusBadRequest, *(models.NewErrorResponse("invalid id format")), nil)
		return
	}

	idStr64, _ := strconv.ParseUint(idStr, 10, 64)
	id := uint32(idStr64)

	// connect db
	db, err := database.ConnectDB()
	if err != nil || db == nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse(err.Error())), nil)
		return
	}

	var dbUser models.User
	if err := db.Where("id = ?", id).First(&dbUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JSONResponseWriter(&w, http.StatusNotFound,
				*(models.NewErrorResponse("can't find this user")), nil)
			return
		}

		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	if userID != id {
		utils.JSONResponseWriter(&w, http.StatusForbidden,
			*(models.NewErrorResponse("can't do any action at this user")), nil)
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
	// ambil id user dari context
	userID := context.Get(r, "id").(uint32)

	// parse body paylod json request
	var userUpdateReq models.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&userUpdateReq); err != nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest, *(models.NewErrorResponse("invalid body format")), nil)
		return
	}

	var user, dbUser models.User
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
			utils.JSONResponseWriter(&w, http.StatusNotFound, *(models.NewErrorResponse("can't find specified user")), nil)
			return
		}
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse(err.Error())), nil)
		return
	}

	if userID != dbUser.ID {
		utils.JSONResponseWriter(&w, http.StatusForbidden, *(models.NewErrorResponse("can't do specified action as this user")), nil)
		return
	}

	dbUser = models.User{}
	err = db.Where("id <> ? AND email = ?", user.ID, user.Email).First(&dbUser).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse(err.Error())), nil)
		return
	} else if err == nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest, *(models.NewErrorResponse("existing email")), nil)
		return
	}

	if user.Password != "" {
		if user.Password, err = utils.HashPassword(user.Password); err != nil {
			utils.JSONResponseWriter(&w, http.StatusForbidden, *(models.NewErrorResponse(err.Error())), nil)
			return
		}
	}

	if err := db.Model(&user).Updates(user).Error; err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse(err.Error())), nil)
		return
	}

	utils.JSONResponseWriter(&w, http.StatusNoContent, nil, nil)
	return
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	userID := context.Get(r, "id").(uint32)

	if idStr == "" || !utils.IsInteger(idStr) {
		utils.JSONResponseWriter(&w, http.StatusBadRequest, *(models.NewErrorResponse("invalid id format")), nil)
		return
	}

	idStr64, _ := strconv.ParseUint(idStr, 10, 64)
	id := uint32(idStr64)

	// connect db
	db, err := database.ConnectDB()
	if err != nil || db == nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse(err.Error())), nil)
		return
	}

	var dbUser models.User
	if err := db.Where("id = ?", id).First(&dbUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.JSONResponseWriter(&w, http.StatusNotFound,
				*(models.NewErrorResponse("can't find this user")), nil)
			return
		}

		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			*(models.NewErrorResponse(err.Error())), nil)
		return
	}

	if userID != id {
		utils.JSONResponseWriter(&w, http.StatusForbidden, *(models.NewErrorResponse("can't do any action at this user")), nil)
		return
	}

	if err := db.Delete(&dbUser, id).Error; err != nil {
		utils.JSONResponseWriter(&w, http.StatusForbidden, *(models.NewErrorResponse(err.Error())), nil)
		return
	}

	var userRes models.UserResponse
	userRes.InsertFromModel(dbUser)

	utils.JSONResponseWriter(&w, http.StatusOK, userRes, nil)
	return
}
