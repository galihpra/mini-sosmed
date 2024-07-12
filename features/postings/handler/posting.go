package handler

import (
	"BE-Sosmed/features/postings"
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

type PostingHandler struct {
	s postings.Service
}

func New(s postings.Service) postings.Handler {
	return &PostingHandler{
		s: s,
	}
}

func (pc *PostingHandler) Add() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input = new(AddRequest)
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": "input yang diberikan tidak sesuai",
			})
		}

		var urlCloudinary = "cloudinary://533421842888945:Oish5XyXkCiiV6oTW2sEo0lEkGg@dlxvvuhph"

		fileHeader, err := c.FormFile("gambar")

		validate := validator.New(validator.WithRequiredStructEnabled())

		if err := validate.Struct(input); err != nil {
			c.Echo().Logger.Error("Input error :", err.Error())
			return c.JSON(http.StatusBadRequest, map[string]any{
				"message": err.Error(),
				"data":    nil,
			})
		}

		var inputProcess = new(postings.Posting)
		if err != nil {
			inputProcess.Artikel = input.Artikel
		} else {
			log.Println(fileHeader.Filename)

			file, _ := fileHeader.Open()

			var ctx = context.Background()

			cldService, _ := cloudinary.NewFromURL(urlCloudinary)
			resp, _ := cldService.Upload.Upload(ctx, file, uploader.UploadParams{})
			log.Println(resp.SecureURL)

			inputProcess.Artikel = input.Artikel
			inputProcess.Gambar = resp.SecureURL
		}

		result, err := pc.s.TambahPosting(c.Get("user").(*gojwt.Token), *inputProcess)

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

		var response = new(AddResponse)
		response.Artikel = result.Artikel
		response.Gambar = result.Gambar

		return c.JSON(http.StatusCreated, map[string]any{
			"message": "success create data",
			"data":    response,
		})
	}
}

func (pc *PostingHandler) GetAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		posts, err := pc.s.SemuaPosting()

		if err != nil {
			c.Logger().Error("Error getting all posts:", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "Error getting all posts",
			})
		}

		var response = make([]GetResponse, len(posts))

		for i, post := range posts {
			comments, err := pc.s.AmbilComment(post.ID)
			if err != nil {
				c.Logger().Error("Error getting comments for post:", err.Error())
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"message": "Error getting comments for post",
				})
			}

			var commentInfo = make([]CommentInfo, len(comments))
			for j, comment := range comments {
				commentInfo[j] = CommentInfo{
					ID:        comment.ID,
					Komentar:  comment.Komentar,
					PostID:    comment.PostID,
					Username:  comment.Username,
					Image:     comment.Image,
					CreatedAt: comment.CreatedAt,
				}
			}

			response[i] = GetResponse{
				ID:        post.ID,
				Artikel:   post.Artikel,
				Gambar:    post.Gambar,
				Likes:     post.Likes,
				Username:  post.Username,
				Image:     post.Image,
				CreatedAt: post.CreatedAt,
				Comments:  commentInfo,
			}
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success get all posts",
			"data":    response,
		})
	}
}

func (pc *PostingHandler) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "invalid post ID",
			})
		}

		var updateRequest = new(UpdateRequest)
		if err := c.Bind(updateRequest); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "invalid update request",
			})
		}

		updatePosting := postings.Posting{
			ID:      uint(postID),
			Artikel: updateRequest.Artikel,
			Gambar:  updateRequest.Gambar,
		}

		updatedPost, err := pc.s.UpdatePosting(c.Get("user").(*gojwt.Token), updatePosting)
		if err != nil {
			c.Logger().Error("Error updating post:", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "error updating post",
			})
		}

		var response = UpdateResponse{
			Artikel: updatedPost.Artikel,
			Gambar:  updatedPost.Gambar,
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success update post",
			"data":    response,
		})
	}

}

func (pc *PostingHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "invalid post ID",
			})
		}

		err = pc.s.DeletePosting(c.Get("user").(*gojwt.Token), uint(postID))
		if err != nil {
			c.Logger().Error("Error deleting post:", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "error deleting post",
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success deleting post",
		})
	}
}

func (pc *PostingHandler) GetByPostID() echo.HandlerFunc {
	return func(c echo.Context) error {
		postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "invalid post ID",
			})
		}

		post, err := pc.s.AmbilPostingByPostID(uint(postID))
		if err != nil {
			c.Logger().Error("Error getting post:", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "Error getting post",
			})
		}

		comments, err := pc.s.AmbilCommentForDetailPost(post.ID)
		if err != nil {
			c.Logger().Error("Error getting comments for post:", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "Error getting comments for post",
			})
		}

		var commentInfo = make([]CommentInfo, len(comments))
		for j, comment := range comments {
			commentInfo[j] = CommentInfo{
				ID:        comment.ID,
				Komentar:  comment.Komentar,
				PostID:    comment.PostID,
				Username:  comment.Username,
				Image:     comment.Image,
				CreatedAt: comment.CreatedAt,
			}
		}

		response := GetResponse{
			ID:        post.ID,
			Artikel:   post.Artikel,
			Gambar:    post.Gambar,
			Username:  post.Username,
			Image:     post.Image,
			CreatedAt: post.CreatedAt,
			Comments:  commentInfo,
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success get post",
			"data":    response,
		})
	}
}

func (pc *PostingHandler) GetByUsername() echo.HandlerFunc {
	return func(c echo.Context) error {
		posts, err := pc.s.AmbilPostingByUsername(c.Param("username"))

		if err != nil {
			c.Logger().Error("Error getting all posts:", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "Error getting all posts",
			})
		}

		var response = make([]GetResponse, len(posts))

		for i, post := range posts {
			comments, err := pc.s.AmbilComment(post.ID)
			if err != nil {
				c.Logger().Error("Error getting comments for post:", err.Error())
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"message": "Error getting comments for post",
				})
			}

			var commentInfo = make([]CommentInfo, len(comments))
			for j, comment := range comments {
				commentInfo[j] = CommentInfo{
					ID:        comment.ID,
					Komentar:  comment.Komentar,
					PostID:    comment.PostID,
					Username:  comment.Username,
					CreatedAt: comment.CreatedAt,
					Image:     comment.Image,
				}
			}

			response[i] = GetResponse{
				ID:        post.ID,
				Artikel:   post.Artikel,
				Gambar:    post.Gambar,
				Username:  post.Username,
				Image:     post.Image,
				CreatedAt: post.CreatedAt,
				Comments:  commentInfo,
			}
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success get posts",
			"data":    response,
		})
	}
}


func (ph *PostingHandler) LikePost() echo.HandlerFunc {
    return func(c echo.Context) error {
        postIDParam := c.Param("id")
        postIDInt, err := strconv.Atoi(postIDParam)
        if err != nil {
            return c.JSON(http.StatusBadRequest, map[string]interface{}{
                "message": "Invalid postID",
            })
        }
        postID := uint(postIDInt)

        result, err := ph.s.LikePosting(c.Get("user").(*gojwt.Token), postID)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, map[string]interface{}{
                "message": err.Error(),
            })
        }

        return c.JSON(http.StatusOK, map[string]interface{}{
            "message": "Post liked successfully",
            "data":    result,
        })
    }
}
