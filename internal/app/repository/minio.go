package repository

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func (r *Repository) UploadToMinIO(fileHeader *multipart.FileHeader, filePrefix string) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		ext = ".bin"
	}

	latinFilename := fmt.Sprintf("%s_%s%s", filePrefix, uuid.New().String()[:8], ext)

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения заголовочных байтов файла: %w", err)
	}
	contentType := http.DetectContentType(buffer)

	_, err = file.Seek(0, 0)
	if err != nil {
		return "", fmt.Errorf("ошибка сброса указателя файла: %w", err)
	}

	_, err = r.minio.PutObject(
		context.Background(),
		r.minioBucket,
		latinFilename,
		file,
		fileHeader.Size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return "", fmt.Errorf("не удалось сохранить файл в MinIO: %w", err)
	}

	fileURL := fmt.Sprintf("http://127.0.0.1:9000/%s/%s", r.minioBucket, latinFilename)
	return fileURL, nil
}
