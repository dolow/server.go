package server

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Handler struct {
	DocumentRoot string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	status := 200

	if path == "" || path == "/" {
		path = "/index.html"
	}

	fullPath := fmt.Sprintf("%s%s", h.DocumentRoot, filepath.FromSlash(path))

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		status = 404
		path = "/404.html"
	}

	ext := strings.ToLower(filepath.Ext(path))
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", contentType)

	w.WriteHeader(status)
	w.Write(data)
}
