package models

var UploadsDir string = "./uploads"

type UserRequest struct {
	Id       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

// type UserDTO struct {}
// type IDResponse struct {}
// type UserResponse struct {}
