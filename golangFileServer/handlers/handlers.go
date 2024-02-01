package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func AuthStatus(w http.ResponseWriter, r *http.Request) {
	// Retrieve session token from the request
    c, err := r.Cookie("session_token")
    if err != nil {
        // No session token found, user is not authenticated
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

	sessionToken := c.Value
	// Check if the session token is valid
	isAuthenticated := true

	userSession, exists := Sessions[sessionToken]  // Update this line to use the correct variable name
	if !exists || userSession.IsExpired() {
		isAuthenticated = false
	}

    if isAuthenticated {
        // User is authenticated
        w.WriteHeader(http.StatusOK)
    } else {
        // User is not authenticated
        w.WriteHeader(http.StatusUnauthorized)
    }
}

// ListFilesHandler handles the listing of files.
func ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("List Files Called")
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fileNames)
}

// UploadHandler handles file uploads.
func UploadHandler(w http.ResponseWriter, r *http.Request) {
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

// DownloadHandler handles file downloads.
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("filename")

	filePath := "./upload/" + fileName
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		http.Error(w, "Could not get file information", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileStat.Size()))

	io.Copy(w, file)
}

// Sign in and Welcome Handlers
