package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CMullaney01/FileServer/handlers"
)

func main() {
	http.Handle("/upload", handlers.CORSHandler(http.HandlerFunc(handlers.UploadHandler)))
	http.Handle("/download", handlers.CORSHandler(http.HandlerFunc(handlers.DownloadHandler)))
	http.Handle("/list", handlers.CORSHandler(http.HandlerFunc(handlers.ListFilesHandler)))

	s := &http.Server{
		Addr:           ":8080",
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("Starting the server on :8080 (HTTP)")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-sigCh
	log.Println("Shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Error during server shutdown: %v", err)
	}
}
