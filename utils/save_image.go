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
	file_name := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	save_path := filepath.Join("uploads", "recipes", file_name)

	if err := os.MkdirAll(filepath.Dir(save_path), os.ModePerm); err != nil {
		return "", err
	}
	dst, err := os.Create(save_path)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return file_name, nil
}
