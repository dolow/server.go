package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

type Handler struct {
	DocumentRoot string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	ext := strings.ToLower(filepath.Ext(path))
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	data, err := ioutil.ReadFile(fmt.Sprintf("%s%s", h.DocumentRoot, path))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", contentType)

	w.WriteHeader(200)
	w.Write(data)
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: main <port> <absolute document root>")
		return
	}

	port := os.Args[1]
	docRoot := os.Args[2]

	handler := &Handler{
		DocumentRoot: docRoot,
	}
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handler,
	}

	log.Printf("launch server on port %s\n", port)

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGTERM, os.Interrupt)
		<-sigint

		server.Shutdown(context.Background())

		log.Println("server shutdown")
	}()

	server.ListenAndServe()

	log.Println("finishing...")
}
