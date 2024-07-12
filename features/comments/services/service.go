package services

import (
	"BE-Sosmed/features/comments"
	"BE-Sosmed/helper/jwt"
	"errors"
	"strings"

	golangjwt "github.com/golang-jwt/jwt/v5"
)

type CommentService struct {
	m comments.Repository
}

func New(model comments.Repository) comments.Service {
	return &CommentService{
		m: model,
	}
}

func (cs *CommentService) CreateComment(token *golangjwt.Token, newComment comments.Comment) (comments.Comment, error) {
	userID, err := jwt.ExtractToken(token)
	if err != nil {
		return comments.Comment{}, err
	}

	result, err := cs.m.InsertComment(userID, newComment)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return comments.Comment{}, errors.New("komentar sudah ada di dalam sistem")
		}
		return comments.Comment{}, errors.New("terjadi kesalahan server")
	}

	return result, nil
}

func (cs *CommentService) DeleteComment(token *golangjwt.Token, commentID uint) error {
	userID, err := jwt.ExtractToken(token)
	if err != nil {
		return err
	}

	err = cs.m.DeleteComment(userID, commentID)
	if err != nil {
		return err
	}

	return nil
}

func (cs *CommentService) PutComment(token *golangjwt.Token, updateComment comments.Comment) (comments.Comment, error) {
	userID, err := jwt.ExtractToken(token)
	if err != nil {
		return comments.Comment{}, err
	}

	result, err := cs.m.UpdateComment(userID, updateComment)
	if err != nil {
		return comments.Comment{}, err
	}

	return result, nil
}
