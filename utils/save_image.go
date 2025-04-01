package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func SaveImageLocally(file multipart.File, header *multipart.FileHeader) (string, error) {
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	savePath := filepath.Join("uploads", "recipes", fileName)

	if err := os.MkdirAll(filepath.Dir(savePath), os.ModePerm); err != nil {
		return "", err
	}
	dst, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return savePath, nil
}
