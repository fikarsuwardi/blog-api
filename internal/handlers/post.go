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

type PostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// CreatePost - Buat post baru (dengan transaksi)
func CreatePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req PostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validasi input
	valid, errMsg := ValidateRequired(map[string]string{
		"title":   req.Title,
		"content": req.Content,
	})
	if !valid {
		HandleValidationError(w, errMsg)
		return
	}

	if !ValidateStringLength(req.Title, 3, 200) {
		HandleValidationError(w, "Title must be between 3 and 200 characters")
		return
	}

	if !ValidateStringLength(req.Content, 10, 10000) {
		HandleValidationError(w, "Content must be between 10 and 10000 characters")
		return
	}

	// Mulai transaksi
	tx := database.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}

	if err := tx.Create(&post).Error; err != nil {
		tx.Rollback()
		respondError(w, http.StatusInternalServerError, "Failed to create post")
		return
	}

	// Preload user data
	if err := tx.Preload("User").First(&post, post.ID).Error; err != nil {
		tx.Rollback()
		respondError(w, http.StatusInternalServerError, "Failed to load post data")
		return
	}

	tx.Commit()
	respondJSON(w, http.StatusCreated, post)
}

// GetPosts - Ambil semua posts
func GetPosts(w http.ResponseWriter, r *http.Request) {
	var posts []models.Post

	if err := database.GetDB().Preload("User").Preload("Comments").Find(&posts).Error; err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch posts")
		return
	}

	respondJSON(w, http.StatusOK, posts)
}

// GetPost - Ambil single post by ID
func GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	var post models.Post
	if err := database.GetDB().Preload("User").Preload("Comments.User").First(&post, postID).Error; err != nil {
		respondError(w, http.StatusNotFound, "Post not found")
		return
	}

	respondJSON(w, http.StatusOK, post)
}

// UpdatePost - Update post (dengan transaksi)
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	postID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	var req PostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validasi input
	valid, errMsg := ValidateRequired(map[string]string{
		"title":   req.Title,
		"content": req.Content,
	})
	if !valid {
		HandleValidationError(w, errMsg)
		return
	}

	if !ValidateStringLength(req.Title, 3, 200) {
		HandleValidationError(w, "Title must be between 3 and 200 characters")
		return
	}

	if !ValidateStringLength(req.Content, 10, 10000) {
		HandleValidationError(w, "Content must be between 10 and 10000 characters")
		return
	}

	// Mulai transaksi
	tx := database.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var post models.Post
	if err := tx.First(&post, postID).Error; err != nil {
		tx.Rollback()
		respondError(w, http.StatusNotFound, "Post not found")
		return
	}

	// Cek ownership
	if post.UserID != userID {
		tx.Rollback()
		respondError(w, http.StatusForbidden, "You can only update your own posts")
		return
	}

	// Update post
	post.Title = req.Title
	post.Content = req.Content

	if err := tx.Save(&post).Error; err != nil {
		tx.Rollback()
		respondError(w, http.StatusInternalServerError, "Failed to update post")
		return
	}

	// Preload user data
	if err := tx.Preload("User").First(&post, post.ID).Error; err != nil {
		tx.Rollback()
		respondError(w, http.StatusInternalServerError, "Failed to load post data")
		return
	}

	tx.Commit()
	respondJSON(w, http.StatusOK, post)
}

// DeletePost - Hapus post (dengan transaksi, soft delete)
func DeletePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	postID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	// Mulai transaksi
	tx := database.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var post models.Post
	if err := tx.First(&post, postID).Error; err != nil {
		tx.Rollback()
		respondError(w, http.StatusNotFound, "Post not found")
		return
	}

	// Cek ownership
	if post.UserID != userID {
		tx.Rollback()
		respondError(w, http.StatusForbidden, "You can only delete your own posts")
		return
	}

	// Soft delete post (dan comments akan ikut ter-cascade karena foreign key)
	if err := tx.Delete(&post).Error; err != nil {
		tx.Rollback()
		respondError(w, http.StatusInternalServerError, "Failed to delete post")
		return
	}

	tx.Commit()
	respondJSON(w, http.StatusOK, map[string]string{"message": "Post deleted successfully"})
}
