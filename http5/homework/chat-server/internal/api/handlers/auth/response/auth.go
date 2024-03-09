package response

type RegistrationResponse struct {
	Username string `json:"username"`
}

type JWTResponse struct {
	Token string `json:"token"`
}
