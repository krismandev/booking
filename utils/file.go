package utils

import (
	"errors"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/sirupsen/logrus"
)

func SaveImageFile(sourceFile *multipart.FileHeader, filename string, path string) (string, error) {
	var err error
	var filePath string

	npwpSrc, err := sourceFile.Open()
	if err != nil {
		logrus.Errorf("failed to open image: %v", err)
		return filePath, err
	}

	isImage := isImage(sourceFile)
	if !isImage {
		return filePath, &BadRequestError{Code: 400, Message: "not a valid image"}
	}

	filePath = filepath.Join(path, filename+GenerateRandomString(30)+".jpg")

	npwpImg, err := imaging.Decode(npwpSrc)
	if err != nil {
		logrus.Errorf("failed to decode image: %v", err)
		return filePath, errors.New("failed to save " + filename + " file ")
	}
	resizedNpwpImg := imaging.Resize(npwpImg, 800, 600, imaging.Lanczos)

	err = imaging.Save(resizedNpwpImg, filePath)
	if err != nil {
		logrus.Errorf("failed to save image: %v", err)
		return filePath, errors.New("failed to save " + filename + " file ")
	}

	return filePath, err
}

func isImage(fileHeader *multipart.FileHeader) bool {
	if fileHeader == nil {
		return false
	}

	// Open file
	file, err := fileHeader.Open()
	if err != nil {
		return false
	}
	defer file.Close()

	// Read first 512 bytes to detect MIME type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return false
	}

	// Get MIME type
	mimeType := http.DetectContentType(buffer)

	// List of allowed image MIME types
	allowedTypes := map[string]bool{
		"image/png":  true,
		"image/jpeg": true,
		"image/jpg":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	return allowedTypes[mimeType]
}
