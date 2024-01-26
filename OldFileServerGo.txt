package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// CORSHandler adds CORS headers to the provided http.Handler.
func CORSHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w, r)
	})
}


func listFilesHandler(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir("./upload")
	if err != nil {
		http.Error(w, "Error reading directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	log.Printf("Files: %+v\n", fileNames)

	// Convert file names to JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fileNames)
}


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
	dst, err := os.Create("./upload/" + header.Filename)
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
	filePath := "./upload/" + fileName
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
	http.Handle("/upload", CORSHandler(http.HandlerFunc(uploadHandler)))
	http.Handle("/download", CORSHandler(http.HandlerFunc(downloadHandler)))
	http.Handle("/list", CORSHandler(http.HandlerFunc(listFilesHandler)))
	// Specify the paths to the TLS certificate and key files
	// For development, use a self-signed certificate. For production, replace with a valid certificate.
	// certFile := "./TestPublicPrivateKey/certificate.crt"
	// keyFile := "./TestPublicPrivateKey/private.key"

	// For development, you can use HTTP
	// For production, uncomment the TLSConfig section and use HTTPS
	s := &http.Server{
		Addr:           ":8080", // Use a different port for HTTP
		Handler:        nil,      // opportunity here to use a non-default multiplexer for routing to handle functions
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		// TLSConfig: &tls.Config{
		// 	InsecureSkipVerify: true, // Skip TLS verification -- only for development
		// },
	}

	log.Println("Starting the server on :8080 (HTTP)")

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// if err := s.ListenAndServeTLS(certFile, keyFile); err != nil { ## http option
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
