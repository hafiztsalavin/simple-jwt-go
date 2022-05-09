package auth

// Credentials is data used for sign in
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
