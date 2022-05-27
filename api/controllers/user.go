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

	"github.com/gorilla/context"
	"gorm.io/gorm"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	// get user id from context
	userID := context.Get(r, "id").(uint32)

	// connect db
	db, err := database.ConnectDB()
	if err != nil || db == nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse(err.Error())), nil)
		return
	}

	// check user in db by id
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

	// creating response base on entity user response
	var userRes models.UserResponse
	if err := userRes.UserFromModel(dbUser); err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse(err.Error())), nil)
		return
	}

	utils.JSONResponseWriter(&w, http.StatusOK, userRes, nil)
	return
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := context.Get(r, "id").(uint32)

	// get json object from request body
	var userUpdateReq models.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&userUpdateReq); err != nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest, *(models.NewErrorResponse("invalid body format")), nil)
		return
	}

	// create entity for update
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

	// check authority id (optional)
	if userID != dbUser.ID {
		utils.JSONResponseWriter(&w, http.StatusForbidden, *(models.NewErrorResponse("can't do action at this user")), nil)
		return
	}

	// check when send with same username and email
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

	// update base on entity
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

	// find user by id
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

	// deleting raw of user
	if err := db.Delete(&dbUser, userID).Error; err != nil {
		utils.JSONResponseWriter(&w, http.StatusForbidden, *(models.NewErrorResponse(err.Error())), nil)
		return
	}

	// make response from entity to response
	var userRes models.UserResponse
	userRes.UserFromModel(dbUser)

	utils.JSONResponseWriter(&w, http.StatusOK, userRes, nil)
	return
}

func Logout(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("access_token")
	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusBadRequest, *(models.NewErrorResponse("Already logged out")), nil)
		return
	}

	userId := context.Get(r, "id").(uint32)
	email := context.Get(r, "email").(string)

	// generate jwt
	expiration := time.Now().Add(-(2 * time.Hour))
	accessToken, err := auth.CreateJWTToken(userId, email, expiration)
	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse(err.Error())), nil)
		return
	}
	utils.AddCookie(w, "access_token", accessToken, expiration)

	refreshToken, err := auth.CreateRefreshToken(accessToken, expiration)
	if err != nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError, *(models.NewErrorResponse(err.Error())), nil)
		return
	}
	utils.AddCookie(w, "refresh_token", refreshToken, expiration)

	utils.JSONResponseWriter(&w, http.StatusOK, map[string]interface{}{"message": "successfully logged out"}, nil)
	return
}
