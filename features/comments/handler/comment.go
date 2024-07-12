package handler

import (
	"BE-Sosmed/features/comments"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	gojwt "github.com/golang-jwt/jwt/v5"
)

type commentHandler struct {
	s comments.Service
}

func New(s comments.Service) comments.Handler {
	return &commentHandler{
		s: s,
	}
}

func (cc *commentHandler) Add() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(CommentRequest)
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

		var inputProcess = new(comments.Comment)
		inputProcess.Komentar = input.Komentar
		inputProcess.PostID = input.PostID

		result, err := cc.s.CreateComment(c.Get("user").(*gojwt.Token), *inputProcess)

		if err != nil {
			c.Logger().Error("ERROR Register, explain:", err.Error())
			if strings.Contains(err.Error(), "duplicate") {
				return c.JSON(http.StatusBadRequest, map[string]any{
					"message": "data yang diinputkan sudah terdaftar ada sistem",
				})
			}
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"message": "terjadi permasalahan ketika memproses data",
			})
		}

		var response = new(CommentResponse)
		response.Komentar = result.Komentar

		return c.JSON(http.StatusCreated, map[string]any{
			"message": "success create comment",
			"data":    response,
		})
	}
}

func (cc *commentHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		commentID, err := strconv.ParseUint(c.Param("commentId"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "invalid comment ID",
			})
		}

		err = cc.s.DeleteComment(c.Get("user").(*gojwt.Token), uint(commentID))
		if err != nil {
			c.Logger().Error("Error deleting comment:", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "error deleting comment",
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success delete comment",
		})
	}
}

func (cc *commentHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		commentID, err := strconv.ParseUint(c.Param("commentId"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "invalid comment ID",
			})
		}

		var updateRequest = new(CommentRequest)
		if err := c.Bind(updateRequest); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "input yang diberikan tidak sesuai",
			})
		}

		updatedComment, err := cc.s.PutComment(c.Get("user").(*gojwt.Token), comments.Comment{
			ID:       uint(commentID),
			Komentar: updateRequest.Komentar,
		})

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "terjadi permasalahan ketika memproses data",
			})
		}

		var response = new(CommentResponse)
		response.Komentar = updatedComment.Komentar

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success update comment",
			"data":    response,
		})
	}

}
