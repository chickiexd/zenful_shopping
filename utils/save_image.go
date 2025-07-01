package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/chickiexd/zenful_shopping/internal/env"
	"github.com/chickiexd/zenful_shopping/internal/logger"
	"github.com/google/uuid"
)

func SaveImageLocally(file multipart.File, header *multipart.FileHeader) (string, error) {
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	file_name := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	file_storage := env.GetString("FILE_STORAGE_PATH", "images")
	save_path := filepath.Join(file_storage, "recipes", file_name)
	logger.Log.Debugf("Saving image to %s", save_path)
	if err := os.MkdirAll(filepath.Dir(save_path), os.ModePerm); err != nil {
		logger.Log.Errorf("Failed to create directory %s: %v", filepath.Dir(save_path), err)
		return "", err
	}
	dst, err := os.Create(save_path)
	if err != nil {
		logger.Log.Errorf("Failed to create file %s: %v", save_path, err)
		return "", err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		logger.Log.Errorf("Failed to copy file content to %s: %v", save_path, err)
		return "", err
	}
	logger.Log.Infof("Image saved successfully: %s", save_path)
	return file_name, nil
}
