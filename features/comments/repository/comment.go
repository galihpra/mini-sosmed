package repository

import (
	"BE-Sosmed/features/comments"
	"errors"

	"gorm.io/gorm"
)

type CommentModel struct {
	gorm.Model
	Komentar string
	UserID   uint
	PostID   uint
}

type commentQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) comments.Repository {
	return &commentQuery{
		db: db,
	}
}

func (cq *commentQuery) InsertComment(userID uint, newComment comments.Comment) (comments.Comment, error) {
	var inputData = new(CommentModel)
	inputData.UserID = userID
	inputData.Komentar = newComment.Komentar
	inputData.PostID = newComment.PostID

	if err := cq.db.Create(&inputData).Error; err != nil {
		return comments.Comment{}, err
	}

	newComment.ID = inputData.ID

	return newComment, nil
}

func (cq *commentQuery) DeleteComment(UserID uint, CommentID uint) error {
	var getComment CommentModel

	if err := cq.db.First(&getComment, CommentID).Error; err != nil {
		return errors.New("comment not found")
	}

	if getComment.UserID != UserID {
		return errors.New("you are not authorized to update this comment")
	}

	if err := cq.db.Delete(&CommentModel{}, CommentID).Error; err != nil {
		return err
	}

	return nil
}

func (cq *commentQuery) UpdateComment(UserID uint, updatedComment comments.Comment) (comments.Comment, error) {
	var getComment CommentModel

	if err := cq.db.First(&getComment, updatedComment.ID).Error; err != nil {
		return comments.Comment{}, errors.New("comment not found")
	}

	if getComment.UserID != UserID {
		return comments.Comment{}, errors.New("you are not authorized to update this comment")
	}

	if err := cq.db.Model(&getComment).Updates(CommentModel{
		Komentar: updatedComment.Komentar,
	}).Error; err != nil {
		return comments.Comment{}, err
	}

	result := comments.Comment{
		ID:       getComment.ID,
		Komentar: getComment.Komentar,
		UserID:   getComment.UserID,
	}

	return result, nil
}
