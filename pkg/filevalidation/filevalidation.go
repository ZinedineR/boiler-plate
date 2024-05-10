package filevalidation

import (
	"bufio"
	"encoding/base64"
	"errors"
	"github.com/gabriel-vasile/mimetype"
	"io"
	"log/slog"
	"os"
)

func ValidateImage(base64String string, maxsize int64) (string, error) {
	// Decode the base64 string to bytes
	data, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return "", err
	}

	// Get the size of the decoded data
	fileSize := int64(len(data))

	// Check if the file size exceeds the maximum allowed size
	if fileSize > maxsize {
		return "", errors.New("file size exceeds the maximum allowed size")
	}

	// Detect mimeType
	mimeType := mimetype.Detect(data)

	// Extension or mimeType checker, if not .png/.jpeg/.webp/.bmp return error
	if mimeType.String() != "image/png" && mimeType.String() != "image/jpeg" && mimeType.String() != "image/vnd.mozilla.apng" && mimeType.String() != "image/webp" && mimeType.String() != "image/bmp" {
		return "", errors.New("extension invalid: " + mimeType.String())
	}

	return mimeType.Extension(), nil
}

func LoadImage64(filePath string, filename string) (string, error) {
	// Open file on disk.
	file, err := os.Open(filePath + filename)
	if err != nil {
		slog.Error(err.Error())
		return "", err
	}
	defer file.Close()

	// Read entire JPG into byte slice.
	reader := bufio.NewReader(file)
	content, err := io.ReadAll(reader)
	if err != nil {
		slog.Error(err.Error())
		return "", err
	}

	// Encode as base64.
	encoded := base64.StdEncoding.EncodeToString(content)
	return encoded, nil
}
