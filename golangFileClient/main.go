package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

func uploadFile(filePath string, targetURL string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", targetURL, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Upload failed with status: %v", resp.Status)
	}

	return nil
}

func downloadFile(fileName string, targetURL string) error {
	tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
	client := &http.Client{Transport: tr}
	
	resp, err := client.Get(targetURL + "?filename=" + fileName)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("Download failed with status: %v", resp.Status)
    }

    content, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    err = ioutil.WriteFile(fileName, content, 0644)
    if err != nil {
        return err
    }

    return nil
}

func listFiles(targetURL string) error {
	resp, err := http.Get(targetURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("List files failed with status: %v", resp.Status)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("List of files:")
	fmt.Println(string(content))
	return nil
}


func main() {
	uploadURL := "http://localhost:8080/upload"   // Use HTTP for development
	downloadURL := "http://localhost:8080/download" // Existing download URL
	listURL := "http://localhost:8080/list"         // New list URL

	// Upload file
	err := uploadFile("./test2.txt", uploadURL)
	if err != nil {
		fmt.Println("Upload failed:", err)
		return
	}
	fmt.Println("File uploaded successfully")

	// Download file
	err = downloadFile("./test.txt", downloadURL)
	if err != nil {
		fmt.Println("Download failed:", err)
		return
	}
	fmt.Println("File downloaded successfully")

	// List files
	err = listFiles(listURL)
	if err != nil {
		fmt.Println("List files failed:", err)
		return
	}
}
