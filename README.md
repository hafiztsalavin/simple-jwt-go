# SIMPLE JWT GO

This is just a simple API in Go, with basic authentication using JWT tokens, user management and sets it to cookie for further Access. All written in GO with a PostgreSQL databases.

## Goals of this project:

Learn about JWT Authentication in Golang REST APIs and securing it with Authentication Middleware.

## To-do:

- [x] Creates the migrations/seeds for the database.
- [x] Request for login, returning JWT token and set as a cookies.
- [x] CRUD For the users (Delete, Patch and Read for user).
- [x] Returning access token and refresh token to refresh the access tokens, if current refresh token is valid.
- [x] Documentation Spec API.
- [ ] Docker later.

## User Registration, User Login and Authorization process.
I set access token will be expired in 1 minute and for Refresh token will be expired in 30 minute (just for testing).

This is diagram to show how to User Registration, Login and Authorization process.
![jwt-authentication-flow](https://github.com/hafiztsalavin/simple-jwt-go/blob/main/docs/documentation/token_generate.png)

And this is for Refresh Token:
![refresh-token-jwt-flow](https://github.com/hafiztsalavin/simple-jwt-go/blob/main/docs//documentation/refresh_token.png)


## Folder Description

```
├── api : root folder where the api is running
|   ├── auth : contains auth file that related to auth system
│   ├── controllers : contains controller file that manage requestmodel and  view
│   ├── database : contains file to integrate database with database
│   ├── logger : basically to write out activity the api into log
│   ├── middlewares : contains functions in intermediate that passed by request and response
│   ├── models : represent the app data model
│   ├── routers : define route endpoint api path
|   ├── utils : contains file that help other function
│   └── server.go : contains function that become root startup the api
├── doc : contains documentation of the project included api documentation
├── migrations : contains file for migrate model database to the database
├── seeder : contains file to seed data to database
├── main.go : gateway of the project app
```

### Prerequisite
1. Go >=1.17.
2. Postgresql >= 11.

### Installation

1. Add `.env` file and copy `env.sample` to `.env`. and adjust the environment variable to your own environment.
```
cp env.sample .env
```
2. Install Go library that was listed in go.mod. 
```
go mod tidy
```
3. Make sure postgres is already based on your `.env` file
4. Create empty database to migrate database if you haven't. 
```
go run main.go -- migrate
```
5. Make a seed data. 
```
go run main.go -- seed
```
7. Run api.
```
go run main.go
```