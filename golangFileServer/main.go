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
	"github.com/CMullaney01/FileServer/middleware"
)

func main() {
	// Existing routes for file handling
	http.Handle("/upload", middleware.AuthCORSHandler(http.HandlerFunc(handlers.UploadHandler)))
	http.Handle("/download", middleware.AuthCORSHandler(http.HandlerFunc(handlers.DownloadHandler)))
	http.Handle("/list", middleware.AuthCORSHandler(http.HandlerFunc(handlers.ListFilesHandler)))

	// New routes for authentication -- redirects should be handled on the client side of things
	http.Handle("/signin", middleware.CORSHandler(http.HandlerFunc(handlers.Signin)))
	http.Handle("/welcome",  middleware.CORSHandler(http.HandlerFunc(handlers.Welcome)))
	http.Handle("/refresh",  middleware.CORSHandler(http.HandlerFunc(handlers.Refresh)))
	http.Handle("/logout",  middleware.CORSHandler(http.HandlerFunc(handlers.Logout)))


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
