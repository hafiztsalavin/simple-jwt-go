# SIMPLE JWT GO

This is just a simple API in Go, with basic authentication using JWT tokens, user management and sets it to cookie for further Access. All written in GO with a PostgreSQL databases.

## Goals of this project:

Learn about JWT Authentication in Golang REST APIs and securing it with Authentication Middleware

## To-do:

- [x] Creates the migrations/seeds for the database.
- [x] Request for login, returning JWT token and set as a cookies.
- [x] CRUD For the users (Delete, Patch and Read for user).
- [x] Returning access token and refresh token to refresh the access tokens, if current refresh token is valid.
- [ ] Unit Test everything.
- [ ] Documentation Spec API.
- [ ] Docker later.

## User Registration, User Login and Authorization process.
This is diagram to show how to User Registration, Login and Authorization process.
![jwt-authentication-flow](token_generate.png)


And this is for Refresh Token:
![refresh-token-jwt-flow](refresh_token.png)
