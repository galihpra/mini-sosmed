package service_test

import (
	"BE-Sosmed/features/users"
	"BE-Sosmed/features/users/mocks"
	service "BE-Sosmed/features/users/services"
	enkMock "BE-Sosmed/helper/enkrip/mocks"
	"BE-Sosmed/helper/jwt"
	"errors"
	"testing"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	repo := mocks.NewRepository(t)
	enkrip := enkMock.NewHashInterface(t)
	s := service.New(repo, enkrip)

	var inputData = users.User{FirstName: "Galih", LastName: "Prayoga", Username: "galih123", Email: "galih@gmail.com", Gender: "male", Password: "admin", Hp: "081229081229"}
	var repoData = users.User{FirstName: "Galih", LastName: "Prayoga", Username: "galih123", Email: "galih@gmail.com", Gender: "male", Password: "some string", Hp: "081229081229"}
	var successReturnData = users.User{ID: uint(1), FirstName: "Galih", Username: "galih123"}
	var errorReturnData = users.User{}
	t.Run("Success Case", func(t *testing.T) {
		enkrip.On("HashPassword", inputData.Password).Return("some string", nil).Once()
		repo.On("InsertUser", repoData).Return(successReturnData, nil).Once()
		result, err := s.Register(inputData)

		enkrip.AssertExpectations(t)
		repo.AssertExpectations(t)

		assert.Nil(t, err)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "", result.Password)
	})

	t.Run("Hashing Error Case", func(t *testing.T) {
		enkrip.On("HashPassword", inputData.Password).Return("", errors.New("hashing error")).Once()

		res, err := s.Register(inputData)

		enkrip.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, uint(0), res.ID)
		assert.Equal(t, "", res.FirstName)
	})

	t.Run("Duplicate Case", func(t *testing.T) {
		enkrip.On("HashPassword", inputData.Password).Return("some string", nil).Once()
		repo.On("InsertUser", repoData).Return(errorReturnData, errors.New("duplicate entry")).Once()

		res, err := s.Register(inputData)

		enkrip.AssertExpectations(t)
		repo.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, uint(0), res.ID)
		assert.Equal(t, "", res.FirstName)
	})

}

func TestLogin(t *testing.T) {
	repo := mocks.NewRepository(t)
	enkrip := enkMock.NewHashInterface(t)
	s := service.New(repo, enkrip)

	var inputData = users.User{Email: "galih@gmail.com", Password: "admin"}
	var successReturnData = users.User{FirstName: "Galih", Username: "galih123", Password: "hashed"}
	var errReturnData = users.User{}

	t.Run("Success Case", func(t *testing.T) {
		repo.On("Login", inputData.Email).Return(successReturnData, nil).Once()
		enkrip.On("Compare", successReturnData.Password, inputData.Password).Return(nil).Once()
		res, err := s.Login(inputData.Email, inputData.Password)

		enkrip.AssertExpectations(t)
		repo.AssertExpectations(t)

		assert.Nil(t, err)
		assert.Equal(t, "Galih", res.FirstName)
		assert.Equal(t, "galih123", res.Username)
	})

	t.Run("Failed Not Found Case", func(t *testing.T) {
		repo.On("Login", inputData.Email).Return(errReturnData, errors.New("data tidak ditemukan")).Once()
		res, err := s.Login(inputData.Email, inputData.Password)

		enkrip.AssertExpectations(t)
		repo.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, errReturnData, res)
		assert.Equal(t, "terjadi kesalahan pada sistem", err.Error())
	})

	t.Run("Failed Incorrect Password Case", func(t *testing.T) {
		repo.On("Login", inputData.Email).Return(successReturnData, nil).Once()
		enkrip.On("Compare", successReturnData.Password, inputData.Password).Return(errors.New("terjadi kesalahan pada sistem")).Once()
		res, err := s.Login(inputData.Email, inputData.Password)

		enkrip.AssertExpectations(t)
		repo.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, errReturnData, res)
		assert.Equal(t, "password yang diinputkan salah", err.Error())
	})
}

var userID = uint(1)
var str, _ = jwt.GenerateJWT(userID)
var token, _ = gojwt.Parse(str, func(t *gojwt.Token) (interface{}, error) {
	return []byte("$!1gnK3yyy!!!"), nil
})

func TestGetUserByID(t *testing.T) {
	repo := mocks.NewRepository(t)
	enkrip := enkMock.NewHashInterface(t)
	s := service.New(repo, enkrip)

	dataResult := users.User{
		ID:        1,
		Username:  "galih123",
		FirstName: "Galih",
		LastName:  "Prayoga",
		Email:     "galih@gmail.com",
		Hp:        "081229081229",
		Image:     "profile.png",
		Gender:    "male",
	}

	t.Run("Success Case", func(t *testing.T) {
		repo.On("ReadUserById", dataResult.ID).Return(dataResult, nil).Once()
		result, err := s.GetUserById(dataResult.ID)

		repo.AssertExpectations(t)

		assert.Nil(t, err)
		assert.Equal(t, "Galih", result.FirstName)
		assert.Equal(t, "galih123", result.Username)
	})

	t.Run("Failed Case", func(t *testing.T) {
		repo.On("ReadUserById", userID).Return(users.User{}, errors.New("failed")).Once()
		result, err := s.GetUserById(userID)

		repo.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, "terjadi kesalahan server", err.Error())
		assert.Equal(t, users.User{}, result)
	})

}

func TestUpdateUser(t *testing.T) {
	repo := mocks.NewRepository(t)
	enkrip := enkMock.NewHashInterface(t)
	s := service.New(repo, enkrip)

	existingUser := users.User{
		ID:        1,
		Username:  "galih123",
		FirstName: "Galih",
		LastName:  "Prayoga",
		Email:     "galih@gmail.com",
		Hp:        "081229081229",
		Image:     "profile.png",
		Gender:    "male",
		Password:  "some string",
	}

	updatedUserData := users.User{
		ID:        1,
		Username:  "galih_update",
		FirstName: "UpdatedFirstName",
		LastName:  "UpdatedLastName",
		Email:     "updated_email@gmail.com",
		Hp:        "081234567890",
		Image:     "updated_profile.png",
		Gender:    "female",
		Password:  "some string",
	}

	t.Run("Success Case", func(t *testing.T) {
		enkrip.On("HashPassword", existingUser.Password).Return("some string", nil).Once()
		repo.On("UpdateUser", userID, updatedUserData).Return(updatedUserData, nil).Once()

		result, err := s.PutUser(token, updatedUserData)

		repo.AssertExpectations(t)
		enkrip.AssertExpectations(t)

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "UpdatedFirstName", result.FirstName)
		assert.Equal(t, "galih_update", result.Username)
		assert.Equal(t, "updated_email@gmail.com", result.Email)
		assert.Equal(t, "some string", result.Password)
	})

	t.Run("Failed Case - User Not Found", func(t *testing.T) {
		enkrip.On("HashPassword", existingUser.Password).Return("some string", nil).Once()
		repo.On("UpdateUser", userID, updatedUserData).Return(users.User{}, errors.New("user not found")).Once()

		result, err := s.PutUser(token, updatedUserData)

		repo.AssertExpectations(t)
		enkrip.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, "terjadi kesalahan pada sistem", err.Error())
		assert.Equal(t, users.User{}, result)
	})

	t.Run("Failed Case - Repository Error", func(t *testing.T) {
		enkrip.On("HashPassword", existingUser.Password).Return("some string", nil).Once()
		repo.On("UpdateUser", userID, updatedUserData).Return(users.User{}, errors.New("repository error")).Once()

		result, err := s.PutUser(token, updatedUserData)

		repo.AssertExpectations(t)
		enkrip.AssertExpectations(t)

		assert.Error(t, err)
		assert.Equal(t, "terjadi kesalahan pada sistem", err.Error())
		assert.Equal(t, users.User{}, result)
	})
}

func TestDeleteUser(t *testing.T) {
	repo := mocks.NewRepository(t)
	enkrip := enkMock.NewHashInterface(t)
	s := service.New(repo, enkrip)

	t.Run("Success Case", func(t *testing.T) {
		repo.On("DeleteUser", userID).Return(nil).Once()

		err := s.DeleteUser(token)

		repo.AssertExpectations(t)

		assert.Nil(t, err)
	})

	t.Run("Error Case", func(t *testing.T) {
		repo.On("DeleteUser", userID).Return(errors.New("test error")).Once()

		err := s.DeleteUser(token)

		repo.AssertExpectations(t)

		assert.Error(t, err)
		assert.EqualError(t, err, "terjadi kesalahan pada sistem")
	})
}
