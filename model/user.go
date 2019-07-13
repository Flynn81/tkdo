package model

//User represents the data structure for a user
type User struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Hash         []byte
	Status       string
	ClientID     string
	ClientSecret string
}

//UserAccess functions to work with users
type UserAccess interface {
	Create(u *User) *User
	Get(email string) (*User, error)
}
