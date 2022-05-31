package controllers_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"simple-jwt-go/api/controllers"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/steinfletcher/apitest"
)

type TestStruct struct {
	name     string
	body     map[string]interface{}
	success  bool
	expected int
}

func Test_SignUp(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Cannot load .env", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/register", controllers.SignUp).Methods("POST")

	ts := httptest.NewServer(r)

	cases := []TestStruct{
		// {
		// 	name:     "Register Test",
		// 	body:     map[string]interface{}{"username": "Register Test", "email": "testregister@gmail.com", "password": "123"},
		// 	success:  true,
		// 	expected: http.StatusOK,
		// },
		{
			name:     "Bad Request",
			body:     map[string]interface{}{"email": "testregister", "password": "password123"},
			expected: http.StatusBadRequest,
		},
		{
			name:     "Empty Field JSON Request",
			body:     map[string]interface{}{"email": "testregister@gmail.com", "password": "password123"},
			expected: http.StatusBadRequest,
		},
		{
			name:     "Username or Email already use",
			body:     map[string]interface{}{"username": "Register Test", "email": "testregister@gmail.com", "password": "123"},
			success:  true,
			expected: http.StatusBadRequest,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			body, _ := json.Marshal(c.body)

			test := apitest.New().Debug().Handler(r).Post("/register").JSON(body).Expect(t).Status(c.expected)

			test.End()

		})
	}

	ts.Close()
}

func Test_SignIn(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Cannot load .env", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/login", controllers.SignIn).Methods("POST")

	ts := httptest.NewServer(r)

	cases := []TestStruct{
		{
			name:     "AccountDoesntExist",
			body:     map[string]interface{}{"email": "testLoginUser1", "password": "password123"},
			expected: http.StatusUnauthorized,
		},
		{
			name:     "WrongPass",
			body:     map[string]interface{}{"email": "testregister@gmail.com", "password": "password123"},
			expected: http.StatusUnauthorized,
		},
		{
			name:     "AccountExistsGetToken",
			body:     map[string]interface{}{"email": "testregister@gmail.com", "password": "123"},
			success:  true,
			expected: http.StatusOK,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			body, _ := json.Marshal(c.body)

			test := apitest.New().Debug().Handler(r).Post("/login").JSON(body).Expect(t).Status(c.expected)

			if c.success {
				test.CookiePresent("access_token")
				test.CookiePresent("refresh_token")
			} else {
				test.CookieNotPresent("access_token")
				test.CookieNotPresent("refresh_token")
			}
			test.End()

		})
	}
	ts.Close()
}
