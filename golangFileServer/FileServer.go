package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST requests are allowed", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Printf("Uploaded file: %+v\n", header.Filename)

	// Save the file
	dst, err := os.Create("./uploads/" + header.Filename)
	if err != nil {
		http.Error(w, "Error saving the file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error saving the file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File successfully uploaded"))
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the file name from the request, for example, from URL parameters
	fileName := r.URL.Query().Get("filename")

	// Open the file on the server
	filePath := "./uploads/" + fileName
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Obtain file information
	fileStat, err := file.Stat()
	if err != nil {
		http.Error(w, "Could not get file information", http.StatusInternalServerError)
		return
	}


	// Set the appropriate headers for the response
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileStat.Size()))

	// Copy the file content to the response writer
	io.Copy(w, file)
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download", downloadHandler)

	// Specify the paths to the TLS certificate and key files
	// currently using test files and browsers etc will give warnings that it is not secure
	// can use certificate from somewhere like Let's Encrypt
	certFile := "./TestPublicPrivateKey/certificate.crt"
	keyFile := "./TestPublicPrivateKey/private.key"

	// create server
	s := &http.Server{
		Addr:           ":8443",
		Handler:        nil, // opportunity here to use a non default multiplexer for routing to handle functions
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	
	log.Println("Starting the server on :8443 (HTTPS)")
	
	// Handle graceful shutdown
	
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.ListenAndServeTLS(certFile, keyFile); err != nil {
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
