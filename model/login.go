package model

//Login represents the data structure for a login
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//LoginResponse represents the data returned when login is successful
type LoginResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}
