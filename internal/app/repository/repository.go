package repository

import (
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/gorm"
)

type RepositorySettings struct {
	DB              *gorm.DB
	MinioEndpoint   string
	MinioAccessKey  string
	MinioSecretKey  string
	MinioBucketName string
}

type Repository struct {
	db          *gorm.DB
	minio       *minio.Client
	minioBucket string
}

func NewRepository(settings *RepositorySettings) (*Repository, error) {
	minioClient, err := minio.New(settings.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(settings.MinioAccessKey, settings.MinioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к MinIO: %w", err)
	}

	return &Repository{
		db:          settings.DB,
		minio:       minioClient,
		minioBucket: settings.MinioBucketName,
	}, nil
}
