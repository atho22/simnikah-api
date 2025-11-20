package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

// ImgBBResponse represents the response from ImgBB API
type ImgBBResponse struct {
	Success bool `json:"success"`
	Data    struct {
		ID         string `json:"id"`
		Title      string `json:"title"`
		URL        string `json:"url"`
		DisplayURL string `json:"display_url"`
		Size       int    `json:"size"`
		DeleteURL  string `json:"delete_url"`
	} `json:"data"`
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

// UploadToImgBB uploads an image file to ImgBB (free cloud storage)
// Returns the image URL if successful
func UploadToImgBB(fileData []byte, filename string) (string, error) {
	// ImgBB API endpoint
	apiURL := "https://api.imgbb.com/1/upload"
	
	// Get API key from environment (optional, but recommended)
	// If not set, ImgBB will use anonymous upload (limited)
	apiKey := os.Getenv("IMGBB_API_KEY")
	if apiKey == "" {
		// Use anonymous upload (free but with limitations)
		// For production, get free API key from https://api.imgbb.com/
		apiKey = "anonymous"
	}

	// Create multipart form
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add API key
	if err := writer.WriteField("key", apiKey); err != nil {
		return "", fmt.Errorf("failed to write API key: %v", err)
	}

	// Add image file
	part, err := writer.CreateFormFile("image", filename)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %v", err)
	}
	if _, err := part.Write(fileData); err != nil {
		return "", fmt.Errorf("failed to write file data: %v", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", apiURL, &requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send request with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	// Parse JSON response
	var imgbbResp ImgBBResponse
	if err := json.Unmarshal(body, &imgbbResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	// Check if upload was successful
	if !imgbbResp.Success {
		return "", fmt.Errorf("ImgBB upload failed: %s", imgbbResp.Error)
	}

	// Return the image URL
	return imgbbResp.Data.URL, nil
}

// UploadFileFromMultipart uploads a multipart file to ImgBB
func UploadFileFromMultipart(fileData io.Reader, filename string) (string, error) {
	// Read file data
	data, err := io.ReadAll(fileData)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// Upload to ImgBB
	return UploadToImgBB(data, filename)
}

