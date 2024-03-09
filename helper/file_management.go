package helper

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// upload file saves the uploaded file to the specified directory with the given filename.
func UploadFile(file io.Reader, directory, filename string) (string, error) {
	// Create the directory if it doesn't exist
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err := os.MkdirAll(directory, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create directory: %v", err)
		}
	}

	// Create the file on the server
	filePath := fmt.Sprintf("%s/%s", directory, filename)
	f, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer f.Close()

	// Copy the uploaded file to the server file
	_, err = io.Copy(f, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file: %v", err)
	}

	// Return the URL or path to access the uploaded file
	return filePath, nil
}

func GetFullImageUrl(ImageUrl string, r *http.Request) string {
	var scheme string
	if r.TLS != nil {
		scheme = "https"
	} else {
		scheme = "http"
	}

	baseURL := fmt.Sprintf("%s://%s", scheme, r.Host)
	var fullImageURL string
	if ImageUrl != "" {
		// Construct the full image URL
		fullImageURL = baseURL + ImageUrl
	}
	fullImageURL = strings.Replace(fullImageURL, "./", "/", 1)

	return fullImageURL
}

func RemoveFile(filePath string) error {
	// Attempt to remove the file
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to remove file: %v", err)
	}

	return nil
}
