package uploading

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/go-faster/errors"
	"github.com/minio/minio-go/v7"
)

type UserUploadUsecase struct {
	s3Client   *minio.Client
	bucketName string
}

func New(s3Client *minio.Client, bucketName string) *UserUploadUsecase {
	return &UserUploadUsecase{
		s3Client:   s3Client,
		bucketName: bucketName,
	}
}

func (u *UserUploadUsecase) Upload(ctx context.Context, userID string, file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return errors.Wrap(err, "Upload file.Open")
	}
	objectName := fmt.Sprintf("user_%s/%s", userID, file.Filename)
	_, err = u.s3Client.PutObject(
		ctx,
		u.bucketName,
		objectName,
		src,
		file.Size,
		minio.PutObjectOptions{
			ContentType: file.Header.Get("Content-Type"),
		},
	)
	if err != nil {
		return errors.Wrap(err, "Uplaod s3Client.PutObject")
	}
	return nil
}
