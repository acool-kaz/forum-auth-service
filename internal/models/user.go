package models

type User struct {
	Id        uint   `json:"id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
}

type userCtxSelect string

var (
	Email    userCtxSelect = "email"
	Username userCtxSelect = "username"
)
