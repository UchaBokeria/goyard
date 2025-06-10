package upload

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// Response represents the result of an upload attempt.
type Response struct {
	ID      int
	Message string
	Success bool
}

// Dir is the directory where uploaded files are stored. You may overwrite this
// variable at runtime before calling Save to change the destination.
var Dir = "./public/uploads/"

// Save stores the uploaded file on disk using a SHA-256 hash as the filename
// (plus its original extension). It returns a Response describing the result.
func Save(file *multipart.FileHeader) *Response {
	src, err := file.Open()
	if err != nil {
		return &Response{ID: -1, Message: "Error opening received file", Success: false}
	}
	defer src.Close()

	// Hash the file contents.
	h := sha256.New()
	if _, err := io.Copy(h, src); err != nil {
		return &Response{ID: -1, Message: "Error calculating hash", Success: false}
	}

	if _, err := src.Seek(0, 0); err != nil {
		return &Response{ID: -1, Message: "Error resetting file reader", Success: false}
	}

	ext := getFileExtension(file)
	hashName := hex.EncodeToString(h.Sum(nil))

	if err := os.MkdirAll(Dir, 0o755); err != nil {
		return &Response{ID: -1, Message: "Unable to create uploads directory", Success: false}
	}

	dstPath := filepath.Join(Dir, hashName+ext)
	dst, err := os.Create(dstPath)
	if err != nil {
		return &Response{ID: -1, Message: "Error creating file on server", Success: false}
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return &Response{ID: -1, Message: "Error copying file to destination", Success: false}
	}

	return &Response{ID: 0, Message: "Successfully uploaded", Success: true}
}

func getFileExtension(file *multipart.FileHeader) string {
	for i := len(file.Filename) - 1; i >= 0; i-- {
		if file.Filename[i] == '.' {
			return strings.ToLower(file.Filename[i:])
		}
	}
	return ""
}
