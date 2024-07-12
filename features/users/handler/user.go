package handler

import (
	"BE-Sosmed/features/users"
	"BE-Sosmed/helper/jwt"
	"BE-Sosmed/helper/responses"
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/go-playground/validator/v10"
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type userHandler struct {
	s users.Service
}

func New(s users.Service) users.Handler {
	return &userHandler{
		s: s,
	}
}

func (uh *userHandler) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(RegisterRequest)
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input yang diberikan tidak sesuai",
			})
		}

		validate := validator.New(validator.WithRequiredStructEnabled())

		if err := validate.Struct(input); err != nil {
			c.Echo().Logger.Error("Input error :", err.Error())
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": err.Error(),
				"data":    nil,
			})
		}

		var inputProcess = new(users.User)
		inputProcess.FirstName = input.FirstName
		inputProcess.LastName = input.LastName
		inputProcess.Gender = input.Gender
		inputProcess.Hp = input.Hp
		inputProcess.Email = input.Email
		inputProcess.Password = input.Password
		inputProcess.Username = input.Username

		result, err := uh.s.Register(*inputProcess)

		if err != nil {
			c.Logger().Error("ERROR Register, explain:", err.Error())
			if strings.Contains(err.Error(), "Duplicate entry") {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"message": "data yang diinputkan sudah terdaftar pada sistem",
				})
			}
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"message": "terjadi permasalahan ketika memproses data",
			})
		}

		var response = new(RegisterResponse)
		response.Username = result.Username
		response.FirstName = result.FirstName

		return responses.PrintResponse(c, http.StatusCreated, "success create data", response)
		// return c.JSON(http.StatusCreated, map[string]any{
		// 	"message": "success create data",
		// 	"data":    response,
		// })
	}
}

func (uh *userHandler) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(LoginRequest)
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input yang diberikan tidak sesuai",
			})
		}

		validate := validator.New(validator.WithRequiredStructEnabled())

		if err := validate.Struct(input); err != nil {
			c.Echo().Logger.Error("Input error :", err.Error())
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": err.Error(),
				"data":    nil,
			})
		}

		result, err := uh.s.Login(input.Email, input.Password)

		if err != nil {
			c.Logger().Error("ERROR Login, explain:", err.Error())
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusBadRequest, map[string]any{
					"message": "data yang diinputkan tidak ditemukan",
				})
			}
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"message": "terjadi permasalahan ketika memproses data",
			})
		}

		strToken, err := jwt.GenerateJWT(result.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"message": "terjadi permasalahan ketika mengenkripsi data",
			})
		}

		var response = new(LoginResponse)
		response.FirstName = result.FirstName
		response.Username = result.Username
		response.Token = strToken

		return c.JSON(http.StatusOK, map[string]any{
			"message": "login success",
			"data":    response,
		})
	}
}

func (uh *userHandler) ReadById() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Param("id")

		// Konversi string ID ke uint
		id, err := strconv.ParseUint(userID, 10, 32)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "data yang diinputkan tidak sesuai format",
			})
		}

		result, err := uh.s.GetUserById(uint(id))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"message": "terjadi permasalahan ketika memproses data",
			})
		}

		var response = new(RegisterResponse)
		response.Username = result.Username
		response.FirstName = result.FirstName

		return c.JSON(http.StatusOK, map[string]any{
			"message": "read data success",
			"data":    response,
		})
	}
}

func (uh *userHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		var updateRequest = new(RegisterRequest)
		if err := c.Bind(updateRequest); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "input yang diberikan tidak sesuai",
			})
		}

		var urlCloudinary = "cloudinary://351966992153882:J1ZB7xXKOl_27eVbba5HN_zQ7sQ@dhxzinjxp"

		fileHeader, err := c.FormFile("foto_profil")

		var inputProcess = new(users.User)
		if err != nil {
			inputProcess.FirstName = updateRequest.FirstName
			inputProcess.LastName = updateRequest.LastName
			inputProcess.Password = updateRequest.Password
			inputProcess.Username = updateRequest.Username
			inputProcess.Email = updateRequest.Email
			inputProcess.Hp = updateRequest.Hp
			inputProcess.Gender = updateRequest.Gender
		} else {
			log.Println(fileHeader.Filename)

			file, _ := fileHeader.Open()

			var ctx = context.Background()

			cldService, _ := cloudinary.NewFromURL(urlCloudinary)

			resp, _ := cldService.Upload.Upload(ctx, file, uploader.UploadParams{})

			log.Println(resp.SecureURL)

			inputProcess.FirstName = updateRequest.FirstName
			inputProcess.LastName = updateRequest.LastName
			inputProcess.Username = updateRequest.Username
			inputProcess.Email = updateRequest.Email
			inputProcess.Hp = updateRequest.Hp
			inputProcess.Gender = updateRequest.Gender
			inputProcess.Image = resp.SecureURL
			if updateRequest.Password != "" {
				inputProcess.Password = updateRequest.Password
			}
		}

		updatedUser, err := uh.s.PutUser(c.Get("user").(*gojwt.Token), *inputProcess)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "terjadi permasalahan ketika memproses data",
			})
		}

		var response = new(RegisterResponse)
		response.Username = updatedUser.Username
		response.FirstName = updatedUser.FirstName

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "update data success",
			"data":    response,
		})
	}
}

func (uh *userHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		err := uh.s.DeleteUser(c.Get("user").(*gojwt.Token))

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "terjadi permasalahan ketika memproses data",
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "delete data success",
		})
	}
}

func (uh *userHandler) ReadByUsername() echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Param("username")
		user, err := uh.s.GetUserByUsername(username)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Failed to get user",
			})
		}
		var response = new(GetResponse)
		response.Username = user.Username
		response.FirstName = user.FirstName
		response.LastName = user.LastName
		response.Email = user.Email
		response.Gender = user.Gender
		response.Hp = user.Hp
		response.CreatedAt = user.CreatedAt
		response.Image = user.Image

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "read data success",
			"data":    response,
		})
	}
}
