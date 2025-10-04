package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"blog-api/internal/database"
	"blog-api/internal/middleware"
	"blog-api/internal/models"

	"github.com/gorilla/mux"
)

type CommentRequest struct {
	Content string `json:"content"`
}

// CreateComment - Buat comment baru pada post (dengan transaksi)
func CreateComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	postID, err := strconv.ParseUint(vars["post_id"], 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	var req CommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validasi input
	valid, errMsg := ValidateRequired(map[string]string{
		"content": req.Content,
	})
	if !valid {
		HandleValidationError(w, errMsg)
		return
	}

	if !ValidateStringLength(req.Content, 1, 1000) {
		HandleValidationError(w, "Content must be between 1 and 1000 characters")
		return
	}

	// Mulai transaksi
	tx := database.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Cek apakah post exists
	var post models.Post
	if err := tx.First(&post, postID).Error; err != nil {
		tx.Rollback()
		respondError(w, http.StatusNotFound, "Post not found")
		return
	}

	comment := models.Comment{
		Content: req.Content,
		UserID:  userID,
		PostID:  uint(postID),
	}

	if err := tx.Create(&comment).Error; err != nil {
		tx.Rollback()
		respondError(w, http.StatusInternalServerError, "Failed to create comment")
		return
	}

	// Preload user data
	if err := tx.Preload("User").First(&comment, comment.ID).Error; err != nil {
		tx.Rollback()
		respondError(w, http.StatusInternalServerError, "Failed to load comment data")
		return
	}

	tx.Commit()
	respondJSON(w, http.StatusCreated, comment)
}

// GetComments - Ambil semua comments untuk post tertentu
func GetComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.ParseUint(vars["post_id"], 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	// Cek apakah post exists
	var post models.Post
	if err := database.GetDB().First(&post, postID).Error; err != nil {
		respondError(w, http.StatusNotFound, "Post not found")
		return
	}

	var comments []models.Comment
	if err := database.GetDB().Preload("User").Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch comments")
		return
	}

	respondJSON(w, http.StatusOK, comments)
}

// DeleteComment - Hapus comment (dengan transaksi, soft delete)
func DeleteComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	postID, err := strconv.ParseUint(vars["post_id"], 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	commentID, err := strconv.ParseUint(vars["comment_id"], 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	// Mulai transaksi
	tx := database.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var comment models.Comment
	if err := tx.Where("id = ? AND post_id = ?", commentID, postID).First(&comment).Error; err != nil {
		tx.Rollback()
		respondError(w, http.StatusNotFound, "Comment not found")
		return
	}

	// Cek ownership
	if comment.UserID != userID {
		tx.Rollback()
		respondError(w, http.StatusForbidden, "You can only delete your own comments")
		return
	}

	// Soft delete comment
	if err := tx.Delete(&comment).Error; err != nil {
		tx.Rollback()
		respondError(w, http.StatusInternalServerError, "Failed to delete comment")
		return
	}

	tx.Commit()
	respondJSON(w, http.StatusOK, map[string]string{"message": "Comment deleted successfully"})
}
