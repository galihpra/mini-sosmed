package handler

type CommentRequest struct {
	Komentar string `json:"komentar"`
	PostID   uint   `json:"post_id"`
}

type CommentResponse struct {
	Komentar string `json:"komentar"`
}
