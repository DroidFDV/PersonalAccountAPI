package storage

import (
	"PersonalAccountAPI/internal/models"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
)

func GetConnectDB(connString string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, errors.Wrap(err, "GetConnectDB pgx.Connect")
	}
	return conn, err
}

// NOTE: возможно стоит перенести в другую директорию
func InitS3Client(config models.S3Config) (*minio.Client, error) {
	minioClient, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, errors.Wrap(err, "InitS3Client minio.New")
	}

	exists, err := minioClient.BucketExists(context.Background(), config.BucketName)
	if err != nil {
		return nil, errors.Wrap(err, "InitS3Client minioClient.BucketExists")
	}
	if !exists {
		if err := minioClient.MakeBucket(context.Background(), config.BucketName, minio.MakeBucketOptions{}); err != nil {
			return nil, errors.Wrap(err, "InitS3Client minioClient.MakeBucket")
		}
	}

	return minioClient, nil
}
