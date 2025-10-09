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

type UserResponse struct {
	ID       int
	Login    string
	Password string
}

func (u *UserRequest) ToDTO() UserDTO {
	return UserDTO{
		ID:       u.ID,
		Login:    u.Login,
		Password: u.Password,
	}
}

func (u *UserDTO) ToResponce() UserResponse {
	return UserResponse{
		ID:       u.ID,
		Login:    u.Login,
		Password: u.Password,
	}
}

func (u *UserRequest) ToResponce() UserResponse {
	return UserResponse{
		ID:       u.ID,
		Login:    u.Login,
		Password: u.Password,
	}
}

// type IDResponse struct {}
// type UserResponse struct {}
