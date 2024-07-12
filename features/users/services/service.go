package service

import (
	"BE-Sosmed/features/users"
	"BE-Sosmed/helper/enkrip"
	"BE-Sosmed/helper/jwt"
	"errors"
	"strings"

	golangjwt "github.com/golang-jwt/jwt/v5"
)

type userService struct {
	repo users.Repository
	h    enkrip.HashInterface
}

func New(r users.Repository, h enkrip.HashInterface) users.Service {
	return &userService{
		repo: r,
		h:    h,
	}
}

func (us *userService) Register(newUser users.User) (users.User, error) {
	// validasi

	// enkripsi password
	enkPassword, err := us.h.HashPassword(newUser.Password)

	if err != nil {
		return users.User{}, errors.New("terjadi masalah saat memproses data")
	}

	newUser.Password = enkPassword

	result, err := us.repo.InsertUser(newUser)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return users.User{}, errors.New("data yang dimasukkan sudah terdaftar")
		}
		return users.User{}, errors.New("terjadi kesalahan pada sistem")
	}

	return result, nil
}
func (us *userService) Login(email string, password string) (users.User, error) {
	result, err := us.repo.Login(email)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return users.User{}, errors.New("data tidak ditemukan")
		}
		return users.User{}, errors.New("terjadi kesalahan pada sistem")
	}

	err = us.h.Compare(result.Password, password)
	if err != nil {
		return users.User{}, errors.New("password yang diinputkan salah")
	}

	return result, nil
}

func (us *userService) GetUserById(UserID uint) (users.User, error) {
	result, err := us.repo.ReadUserById(UserID)

	if err != nil {
		return users.User{}, errors.New("terjadi kesalahan server")
	}

	return result, nil
}

func (us *userService) PutUser(token *golangjwt.Token, updatedUser users.User) (users.User, error) {
	UserID, err := jwt.ExtractToken(token)
	if err != nil {
		return users.User{}, err
	}

	enkPassword, err := us.h.HashPassword(updatedUser.Password)

	if err != nil {
		return users.User{}, errors.New("terjadi masalah saat memproses data")
	}

	updatedUser.Password = enkPassword
	result, err := us.repo.UpdateUser(UserID, updatedUser)

	if err != nil {
		return users.User{}, errors.New("terjadi kesalahan pada sistem")
	}

	return result, nil
}

func (us *userService) DeleteUser(token *golangjwt.Token) error {
	UserID, err := jwt.ExtractToken(token)
	if err != nil {
		return err
	}
	err = us.repo.DeleteUser(UserID)

	if err != nil {
		return errors.New("terjadi kesalahan pada sistem")
	}

	return nil
}

func (us *userService) GetUserByUsername(username string) (users.User, error) {
	return us.repo.ReadUserByUsername(username)
}