package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"github.com/chickiexd/zenful_shopping/internal/service"

	"github.com/go-chi/chi/v5"
)

type ImageHandler struct {
	service *service.Service
}

func (h *ImageHandler) Get(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}
	image_path := filepath.Join("uploads/recipes", filename)
	if _, err := os.Stat(image_path); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}
	switch {
	case strings.HasSuffix(filename, ".png"):
		w.Header().Set("Content-Type", "image/png")
	case strings.HasSuffix(filename, ".jpg"), strings.HasSuffix(filename, ".jpeg"):
		w.Header().Set("Content-Type", "image/jpeg")
	case strings.HasSuffix(filename, ".webp"):
		w.Header().Set("Content-Type", "image/webp")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	http.ServeFile(w, r, image_path)
}
