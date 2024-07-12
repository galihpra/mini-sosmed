package users

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type User struct {
	ID        uint
	FirstName string
	LastName  string
	Gender    string
	Hp        string
	Email     string
	Password  string
	Image     string
	Username  string
	CreatedAt time.Time
}

type Handler interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc
	ReadById() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
	ReadByUsername() echo.HandlerFunc
}

type Service interface {
	Register(newUser User) (User, error)
	Login(email string, password string) (User, error)
	GetUserById(UserID uint) (User, error)
	PutUser(token *jwt.Token, updatedUser User) (User, error)
	DeleteUser(token *jwt.Token) error
	GetUserByUsername(username string) (User, error)
}

type Repository interface {
	InsertUser(newUser User) (User, error)
	Login(email string) (User, error)
	ReadUserById(UserID uint) (User, error)
	UpdateUser(UserID uint, updatedUser User) (User, error)
	DeleteUser(UserID uint) error
	ReadUserByUsername(username string) (User, error)
}
