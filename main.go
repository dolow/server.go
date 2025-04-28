package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dolow/server.go/server"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: main <port> <absolute document root>")
		return
	}

	port := os.Args[1]
	docRoot := os.Args[2]

	handler := &server.Handler{
		DocumentRoot: docRoot,
	}
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handler,
	}

	log.Printf("launch server on port %s\n", port)

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGTERM, os.Interrupt)
		<-sigint

		httpServer.Shutdown(context.Background())

		log.Println("server shutdown")
	}()

	httpServer.ListenAndServe()

	log.Println("finishing...")
}
