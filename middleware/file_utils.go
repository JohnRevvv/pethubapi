package middleware

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
)

// ConvertFileToBase64 converts uploaded file to a base64 string
func ConvertFileToBase64(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	buf := make([]byte, file.Size)
	if _, err := io.ReadFull(src, buf); err != nil {
		return "", fmt.Errorf("failed to read uploaded file: %v", err)
	}

	encoded := base64.StdEncoding.EncodeToString(buf)
	return encoded, nil
}
