package auth

// data used for sign in
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
