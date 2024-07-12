package services_test

import (
	"BE-Sosmed/features/comments"
	"BE-Sosmed/features/comments/mocks"
	"BE-Sosmed/features/comments/services"
	"BE-Sosmed/helper/jwt"
	"errors"
	"testing"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

var userID = uint(1)
var str, _ = jwt.GenerateJWT(userID)
var token, _ = gojwt.Parse(str, func(t *gojwt.Token) (interface{}, error) {
	return []byte("$!1gnK3yyy!!!"), nil
})

var invalidToken, _ = gojwt.Parse(str, func(t *gojwt.Token) (interface{}, error) {
	return []byte("$!1gnK3yyx!!!"), nil
})

func TestCreateComment(t *testing.T) {
	repo := mocks.NewRepository(t)
	s := services.New(repo)

	var inputData = comments.Comment{PostID: uint(2), Komentar: "Wow, amazing"}
	var successReturnData = comments.Comment{Komentar: "Wow, amazing"}
	t.Run("Success Case", func(t *testing.T) {
		repo.On("InsertComment", userID, inputData).Return(successReturnData, nil).Once()

		result, err := s.CreateComment(token, inputData)

		repo.AssertExpectations(t)

		assert.Nil(t, err)
		assert.Equal(t, inputData.Komentar, result.Komentar)
	})

	t.Run("Duplicate Case", func(t *testing.T) {
		repo.On("InsertComment", userID, inputData).Return(comments.Comment{}, errors.New("duplicate entry")).Once()

		result, err := s.CreateComment(token, inputData)

		repo.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, comments.Comment{}, result)
		assert.Equal(t, "komentar sudah ada di dalam sistem", err.Error())
	})

}

func TestUpdateComment(t *testing.T) {
	repo := mocks.NewRepository(t)
	s := services.New(repo)

	var inputData = comments.Comment{PostID: uint(2), Komentar: "Wow, amazing"}
	var successReturnData = comments.Comment{Komentar: "Wow, amazing"}

	t.Run("Success Case", func(t *testing.T) {
		repo.On("UpdateComment", userID, inputData).Return(successReturnData, nil).Once()

		result, err := s.PutComment(token, inputData)

		repo.AssertExpectations(t)

		assert.Nil(t, err)
		assert.Equal(t, inputData.Komentar, result.Komentar)
	})

	t.Run("Failed Case", func(t *testing.T) {
		repo.On("UpdateComment", userID, inputData).Return(comments.Comment{}, errors.New("failed")).Once()

		result, err := s.PutComment(token, inputData)

		repo.AssertExpectations(t)

		assert.Error(t, err)
		assert.EqualError(t, err, "failed")
		assert.Equal(t, comments.Comment{}, result)
	})

}

func TestDeleteComment(t *testing.T) {
	repo := mocks.NewRepository(t)
	s := services.New(repo)

	commentID := uint(2)
	t.Run("Success Case", func(t *testing.T) {
		repo.On("DeleteComment", userID, commentID).Return(nil).Once()

		err := s.DeleteComment(token, commentID)

		repo.AssertExpectations(t)

		assert.Nil(t, err)
	})

	t.Run("Error Case", func(t *testing.T) {
		repo.On("DeleteComment", userID, commentID).Return(errors.New("failed")).Once()

		err := s.DeleteComment(token, commentID)

		repo.AssertExpectations(t)

		assert.Error(t, err)
		assert.EqualError(t, err, "failed")
	})
}

func TestCreateComment_InvalidToken(t *testing.T) {
	repo := mocks.NewRepository(t)
	s := services.New(repo)

	var inputData = comments.Comment{PostID: uint(2), Komentar: "Wow, amazing"}

	result, err := s.CreateComment(invalidToken, inputData)

	assert.Error(t, err)
	assert.Equal(t, comments.Comment{}, result)
}

func TestUpdateComment_InvalidToken(t *testing.T) {
	repo := mocks.NewRepository(t)
	s := services.New(repo)

	var inputData = comments.Comment{PostID: uint(2), Komentar: "Wow, amazing"}

	result, err := s.PutComment(invalidToken, inputData)

	assert.Error(t, err)
	assert.Equal(t, comments.Comment{}, result)
}

func TestDeleteComment_InvalidToken(t *testing.T) {
	repo := mocks.NewRepository(t)
	s := services.New(repo)

	err := s.DeleteComment(invalidToken, uint(22))

	assert.Error(t, err)
}
