package service

import (
	"errors"
	"testing"

	"BE-Sosmed/features/comments"
	"BE-Sosmed/features/postings"
	"BE-Sosmed/features/postings/mocks"
	"BE-Sosmed/features/users"
	um "BE-Sosmed/features/users/mocks"
	"BE-Sosmed/helper/jwt"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

var userID = uint(1)
var str, _ = jwt.GenerateJWT(userID)
var token, _ = gojwt.Parse(str, func(t *gojwt.Token) (interface{}, error) {
	return []byte("$!1gnK3yyy!!!"), nil
})

func TestTambahPosting(t *testing.T) {
	mockRepo := mocks.NewRepository(t)
	mockUser := um.NewService(t)

	postingService := New(mockRepo, mockUser)

	newPosting := postings.Posting{
		Artikel: "Test Artikel",
		Gambar:  "Test Gambar",
	}
	repoData := postings.Posting{
		Artikel: "Test Artikel",
		Gambar:  "Test Gambar",
	}

	t.Run("Success Case", func(t *testing.T) {
		mockRepo.On("InsertPosting", userID, newPosting).Return(repoData, nil).Once()
		result, err := postingService.TambahPosting(token, newPosting)

		mockRepo.AssertExpectations(t)

		assert.Nil(t, err)
		assert.Equal(t, newPosting.Artikel, result.Artikel)
		assert.Equal(t, newPosting.Gambar, result.Gambar)
	})

	t.Run("Failed case", func(t *testing.T) {
		mockRepo.On("InsertPosting", userID, newPosting).Return(postings.Posting{}, errors.New("internal server error")).Once()
		result, err := postingService.TambahPosting(token, newPosting)

		mockRepo.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, postings.Posting{}, result)
		assert.Equal(t, "terjadi kesalahan server", err.Error())

	})

	t.Run("Duplicate case", func(t *testing.T) {
		mockRepo.On("InsertPosting", userID, newPosting).Return(postings.Posting{}, errors.New("duplicate")).Once()
		result, err := postingService.TambahPosting(token, newPosting)

		mockRepo.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, postings.Posting{}, result)
		assert.Equal(t, "posting sudah ada pada sistem", err.Error())

	})
}
func TestGetComment(t *testing.T) {
	mockRepo := mocks.NewRepository(t)
	mockUser := um.NewService(t)

	postingService := New(mockRepo, mockUser)

	dataComment := []comments.Comment{
		{ID: 1, PostID: 1, Komentar: "Amazing", UserID: 1},
		{ID: 2, PostID: 1, Komentar: "Amazing", UserID: 2},
	}

	t.Run("Success Case", func(t *testing.T) {
		mockRepo.On("GetComment", uint(1)).Return(dataComment, nil).Once()
		mockUser.On("GetUserById", uint(1)).Return(users.User{ID: 1, Username: "User1", Image: "Image1"}, nil).Once()
		mockUser.On("GetUserById", uint(2)).Return(users.User{ID: 2, Username: "User2", Image: "Image2"}, nil).Once()
		result, err := postingService.AmbilComment(uint(1))

		mockRepo.AssertExpectations(t)
		mockUser.AssertExpectations(t)

		assert.Nil(t, err)
		assert.Equal(t, "User1", result[0].Username)
		assert.Equal(t, "Image1", result[0].Image)
		assert.Equal(t, "User2", result[1].Username)
		assert.Equal(t, "Image2", result[1].Image)
	})

	t.Run("Failed Case", func(t *testing.T) {
		mockRepo.On("GetComment", uint(1)).Return(nil, errors.New("failed")).Once()

		result, err := postingService.AmbilComment(uint(1))

		mockRepo.AssertExpectations(t)
		mockUser.AssertExpectations(t)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "terjadi kesalahan server", err.Error())

	})

}

func TestGetAllPosting(t *testing.T) {
	mockRepo := mocks.NewRepository(t)
	mockUser := um.NewService(t)

	postingService := New(mockRepo, mockUser)

	dataPosting := []postings.Posting{
		{ID: 1, Artikel: "Artikel 1", Gambar: "Gambar1", UserID: 1},
		{ID: 2, Artikel: "Artikel 2", Gambar: "Gambar2", UserID: 2},
	}

	t.Run("Success Case", func(t *testing.T) {
		mockRepo.On("GetAllPost").Return(dataPosting, nil).Once()
		mockUser.On("GetUserById", uint(1)).Return(users.User{ID: 1, Username: "User1", Image: "Image1"}, nil).Once()
		mockUser.On("GetUserById", uint(2)).Return(users.User{ID: 2, Username: "User2", Image: "Image2"}, nil).Once()
		result, err := postingService.SemuaPosting()

		mockRepo.AssertExpectations(t)
		mockUser.AssertExpectations(t)

		assert.Nil(t, err)
		assert.Equal(t, "User1", result[0].Username)
		assert.Equal(t, "Image1", result[0].Image)
		assert.Equal(t, "User2", result[1].Username)
		assert.Equal(t, "Image2", result[1].Image)
	})

	t.Run("Failed Case", func(t *testing.T) {
		mockRepo.On("GetAllPost").Return(nil, errors.New("terjadi kesalahan server")).Once()
		result, err := postingService.SemuaPosting()

		mockRepo.AssertExpectations(t)
		mockUser.AssertExpectations(t)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "terjadi kesalahan server", err.Error())
	})

}

func TestUpdatePosting(t *testing.T) {
	mockRepo := mocks.NewRepository(t)
	mockUser := um.NewService(t)

	postingService := New(mockRepo, mockUser)

	newPosting := postings.Posting{
		ID:      1,
		Artikel: "Updated Artikel",
		Gambar:  "Updated Gambar",
		UserID:  1,
	}
	repoData := postings.Posting{
		Artikel: "Updated Artikel",
		Gambar:  "Updated Gambar",
	}

	t.Run("Success Case", func(t *testing.T) {
		mockRepo.On("UpdatePost", userID, newPosting).Return(repoData, nil).Once()
		result, err := postingService.UpdatePosting(token, newPosting)

		mockRepo.AssertExpectations(t)

		assert.Nil(t, err)
		assert.Equal(t, newPosting.Artikel, result.Artikel)
		assert.Equal(t, newPosting.Gambar, result.Gambar)
	})

	t.Run("Failed case", func(t *testing.T) {
		mockRepo.On("UpdatePost", userID, newPosting).Return(postings.Posting{}, errors.New("terjadi kesalahan server")).Once()
		result, err := postingService.UpdatePosting(token, newPosting)

		mockRepo.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, postings.Posting{}, result)
		assert.Equal(t, "terjadi kesalahan server", err.Error())

	})
}

func TestDeletePosting(t *testing.T) {
	mockRepo := mocks.NewRepository(t)
	mockUser := um.NewService(t)

	postingService := New(mockRepo, mockUser)

	t.Run("Success Case", func(t *testing.T) {
		mockRepo.On("DeletePost", uint(1), uint(123)).Return(nil).Once()

		err := postingService.DeletePosting(token, uint(123))

		mockRepo.AssertExpectations(t)

		assert.Nil(t, err)
	})

	t.Run("Failed Case", func(t *testing.T) {
		mockRepo.On("DeletePost", uint(1), uint(123)).Return(errors.New("terjadi kesalahan server")).Once()

		err := postingService.DeletePosting(token, uint(123))

		mockRepo.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, "terjadi kesalahan server", err.Error())
	})

}

func TestGetPostingByPostID(t *testing.T) {
	mockRepo := mocks.NewRepository(t)
	mockUser := um.NewService(t)

	postingService := New(mockRepo, mockUser)

	dataPosting := postings.Posting{ID: 1, Artikel: "Artikel 1", Gambar: "Gambar1", UserID: 1}

	t.Run("Success Case", func(t *testing.T) {
		mockRepo.On("GetPostByPostID", uint(123)).Return(dataPosting, nil).Once()
		mockUser.On("GetUserById", uint(1)).Return(users.User{ID: 1, Username: "User1", Image: "Image1"}, nil).Once()
		result, err := postingService.AmbilPostingByPostID(uint(123))

		mockRepo.AssertExpectations(t)
		mockUser.AssertExpectations(t)

		assert.Nil(t, err)
		assert.Equal(t, "User1", result.Username)
		assert.Equal(t, "Image1", result.Image)
	})

}

func TestGetPostingByUsername(t *testing.T) {
	mockRepo := mocks.NewRepository(t)
	mockUser := um.NewService(t)

	postingService := New(mockRepo, mockUser)

	dataPosting := []postings.Posting{
		{ID: 1, Artikel: "Artikel 1", Gambar: "Gambar1", UserID: 1},
		{ID: 2, Artikel: "Artikel 2", Gambar: "Gambar2", UserID: 2},
	}

	t.Run("Success Case", func(t *testing.T) {
		mockRepo.On("GetPostByUsername", "user1").Return(dataPosting, nil).Once()
		mockUser.On("GetUserById", uint(1)).Return(users.User{ID: 1, Username: "User1", Image: "Image1"}, nil).Once()
		mockUser.On("GetUserById", uint(2)).Return(users.User{ID: 2, Username: "User2", Image: "Image2"}, nil).Once()
		result, err := postingService.AmbilPostingByUsername("user1")

		mockRepo.AssertExpectations(t)
		mockUser.AssertExpectations(t)

		assert.Nil(t, err)
		assert.Equal(t, "User1", result[0].Username)
		assert.Equal(t, "Image1", result[0].Image)
		assert.Equal(t, "User2", result[1].Username)
		assert.Equal(t, "Image2", result[1].Image)
	})

	t.Run("Failed Case", func(t *testing.T) {
		mockRepo.On("GetPostByUsername", "user1").Return(nil, errors.New("terjadi kesalahan server")).Once()
		result, err := postingService.AmbilPostingByUsername("user1")

		mockRepo.AssertExpectations(t)
		mockUser.AssertExpectations(t)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "terjadi kesalahan server", err.Error())
	})

}
