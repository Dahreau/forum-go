package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const uploadDir = "assets/img/uploads/"

func UploadImageHandler(w http.ResponseWriter, r *http.Request) (string, error) {
	// Parse multipart form data
	err := r.ParseMultipartForm(20 * 1024 * 1024) // 20MB limit
	if err != nil {
		http.Error(w, "File size exceeds limit", http.StatusBadRequest)
		return "", err
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Invalid file upload", http.StatusBadRequest)
		return "", err
	}
	defer file.Close()

	// Check file extension
	ext := filepath.Ext(handler.Filename)
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif"}
	if !contains(allowedExtensions, ext) {
		http.Error(w, "Unsupported file type", http.StatusUnsupportedMediaType)
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}

	// Create destination file
	newFileName := fmt.Sprintf("%d%s", time.Now().Unix(), ext)
	destPath := filepath.Join(uploadDir, newFileName)

	out, err := os.Create(destPath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return "", err
	}
	defer out.Close()

	// Write to the destination
	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return "", err
	}

	return destPath, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
