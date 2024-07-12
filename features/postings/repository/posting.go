package repository

import (
	"BE-Sosmed/features/comments"
	cr "BE-Sosmed/features/comments/repository"
	"BE-Sosmed/features/postings"
	"errors"

	"gorm.io/gorm"
)

type PostingModel struct {
	gorm.Model
	Artikel  string
	Gambar   string
	Likes    int
	UserID   uint
	Comments []cr.CommentModel `gorm:"foreignKey:PostID"`
}

type postingQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) postings.Repository {
	return &postingQuery{
		db: db,
	}
}

func (pq *postingQuery) InsertPosting(userID uint, newPosting postings.Posting) (postings.Posting, error) {
	var inputData = new(PostingModel)
	inputData.UserID = userID
	inputData.Artikel = newPosting.Artikel
	inputData.Gambar = newPosting.Gambar

	if err := pq.db.Create(&inputData).Error; err != nil {
		return postings.Posting{}, err
	}

	newPosting.ID = inputData.ID

	return newPosting, nil
}

func (pq *postingQuery) GetComment(PostID uint) ([]comments.Comment, error) {
	var commentModels []cr.CommentModel

	if err := pq.db.Where("post_id = ?", PostID).Order("created_at desc").Limit(2).Find(&commentModels).Error; err != nil {
		return nil, err
	}

	var result []comments.Comment
	for _, model := range commentModels {
		result = append(result, comments.Comment{
			ID:        model.ID,
			Komentar:  model.Komentar,
			PostID:    model.PostID,
			UserID:    model.UserID,
			CreatedAt: model.CreatedAt,
		})
	}

	return result, nil
}

func (pq *postingQuery) GetCommentForDetailPost(PostID uint) ([]comments.Comment, error) {
	var commentModels []cr.CommentModel

	if err := pq.db.Where("post_id = ?", PostID).Order("created_at desc").Find(&commentModels).Error; err != nil {
		return nil, err
	}

	var result []comments.Comment
	for _, model := range commentModels {
		result = append(result, comments.Comment{
			ID:        model.ID,
			Komentar:  model.Komentar,
			PostID:    model.PostID,
			UserID:    model.UserID,
			CreatedAt: model.CreatedAt,
		})
	}

	return result, nil
}

func (pq *postingQuery) GetAllPost() ([]postings.Posting, error) {
	var posts []PostingModel

	if err := pq.db.Order("created_at desc").Find(&posts).Error; err != nil {
		return nil, err
	}
	var result []postings.Posting
	for _, post := range posts {
		result = append(result, postings.Posting{
			ID:        post.ID,
			Artikel:   post.Artikel,
			Gambar:    post.Gambar,
			Likes:     post.Likes,
			UserID:    post.UserID,
			CreatedAt: post.CreatedAt,
		})
	}

	return result, nil
}

func (pq *postingQuery) UpdatePost(userID uint, updatePosting postings.Posting) (postings.Posting, error) {
	var existingPost PostingModel

	if err := pq.db.First(&existingPost, updatePosting.ID).Error; err != nil {
		return postings.Posting{}, errors.New("posting not found")
	}

	if existingPost.UserID != userID {
		return postings.Posting{}, errors.New("you are not authorized to update this post")
	}

	if err := pq.db.Model(&existingPost).Updates(PostingModel{
		Artikel: updatePosting.Artikel,
		Gambar:  updatePosting.Gambar,
	}).Error; err != nil {
		return postings.Posting{}, err
	}

	updatedPost := postings.Posting{
		ID:      existingPost.ID,
		Artikel: existingPost.Artikel,
		Gambar:  existingPost.Gambar,
		UserID:  existingPost.UserID,
	}

	return updatedPost, nil
}

func (pq *postingQuery) DeletePost(userID uint, postID uint) error {
	var existingPost PostingModel

	if err := pq.db.First(&existingPost, postID).Error; err != nil {
		return errors.New("posting not found")
	}

	if existingPost.UserID != userID {
		return errors.New("you are not authorized to delete this post")
	}

	if err := pq.db.Delete(&existingPost).Error; err != nil {
		return err
	}

	return nil
}

func (pq *postingQuery) GetPostByPostID(PostID uint) (postings.Posting, error) {
	var post PostingModel

	if err := pq.db.Where("id = ?", PostID).Order("created_at desc").First(&post).Error; err != nil {
		return postings.Posting{}, err
	}

	result := postings.Posting{
		ID:        post.ID,
		Artikel:   post.Artikel,
		Gambar:    post.Gambar,
		UserID:    post.UserID,
		CreatedAt: post.CreatedAt,
	}

	return result, nil
}

func (pq *postingQuery) GetPostByUsername(Username string) ([]postings.Posting, error) {
	var posts []PostingModel
	err := pq.db.Model(&PostingModel{}).
		Select("posting_models.*, user_models.username, user_models.image").
		Joins("JOIN user_models on posting_models.user_id = user_models.id").
		Where("username = ?", Username).
		Scan(&posts).
		Error

	if err != nil {
		return nil, err
	}

	var result []postings.Posting
	for _, post := range posts {
		result = append(result, postings.Posting{
			ID:        post.ID,
			Artikel:   post.Artikel,
			Gambar:    post.Gambar,
			UserID:    post.UserID,
			CreatedAt: post.CreatedAt,
		})
	}

	return result, nil
}


func (pq *postingQuery) LikePosts(userID, postID uint, updatePosting postings.Posting) (postings.Posting, error) {
    var existingPost PostingModel

    if err := pq.db.First(&existingPost, postID).Error; err != nil {
        return postings.Posting{}, errors.New("posting not found")
    }

    if existingPost.UserID != userID {
        return postings.Posting{}, errors.New("you are not authorized to like this post")
    }

    if err := pq.db.Model(&existingPost).Updates(PostingModel{
        Likes: existingPost.Likes + 1,
    }).Error; err != nil {
        return postings.Posting{}, err
    }

    updatedPost := postings.Posting{
        ID:      existingPost.ID,
        Artikel: existingPost.Artikel,
        Gambar:  existingPost.Gambar,
        Likes:   existingPost.Likes,
        UserID:  existingPost.UserID,
    }

    return updatedPost, nil
}
