package models

var UploadsDir string

type UserRequest struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
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
