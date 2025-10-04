package models

import "time"

var UploadsDir string

type UserRequest struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type WrapUser struct {
	User UserRequest
	TTL  time.Time
}

type UserDTO struct {
	ID       int
	Login    string
	Password string
}

func (ur *UserRequest) ToDTO() UserDTO {
	return UserDTO{
		ID:       ur.ID,
		Login:    ur.Login,
		Password: ur.Password,
	}
}

// type IDResponse struct {}
// type UserResponse struct {}
