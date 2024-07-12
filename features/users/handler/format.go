package handler

import "time"

type RegisterRequest struct {
	FirstName string `json:"nama_depan" form:"nama_depan" validate:"required,alpha"`
	LastName  string `json:"nama_belakang" form:"nama_belakang" validate:"required,alpha"`
	Gender    string `json:"gender" form:"gender" validate:"required,alpha"`
	Hp        string `json:"hp" form:"hp" validate:"required,numeric"`
	Email     string `json:"email" form:"email" validate:"required,email"`
	Password  string `json:"password" form:"password" validate:"required"`
	Username  string `json:"username" form:"username"  validate:"required"`
	Image     string `json:"foto_profil" form:"foto_profil"`
}

type RegisterResponse struct {
	Username  string `json:"username"`
	FirstName string `json:"nama_depan"`
}

type LoginRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

type LoginResponse struct {
	FirstName string `json:"nama_depan"`
	Username  string `json:"username"`
	Token     string `json:"token"`
}

type GetResponse struct {
	Username  string    `json:"username"`
	FirstName string    `json:"nama_depan"`
	LastName  string    `json:"nama_belakang"`
	Gender    string    `json:"gender"`
	Hp        string    `json:"hp"`
	Email     string    `json:"email"`
	Image     string    `json:"foto_profil" form:"foto_profil"`
	CreatedAt time.Time `json:"created_at"`
}
