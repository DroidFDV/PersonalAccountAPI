package models

import "time"

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

// NOTE: куда-то надо перенести
type S3Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
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
