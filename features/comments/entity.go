package comments

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Comment struct {
	ID        uint
	Komentar  string
	UserID    uint
	PostID    uint
	Username  string
	Image     string
	CreatedAt time.Time
}

type Handler interface {
	Add() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
}

type Service interface {
	CreateComment(token *jwt.Token, newComment Comment) (Comment, error)
	PutComment(token *jwt.Token, updatedComment Comment) (Comment, error)
	DeleteComment(token *jwt.Token, CommentID uint) error
}

type Repository interface {
	InsertComment(UserID uint, newComment Comment) (Comment, error)
	UpdateComment(UserID uint, updatedComment Comment) (Comment, error)
	DeleteComment(UserID uint, CommentID uint) error
}
