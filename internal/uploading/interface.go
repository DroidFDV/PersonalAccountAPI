package uploading

import (
	"context"
	"mime/multipart"
)

type UploadingProvider interface {
	Upload(ctx context.Context, subDir string, file *multipart.FileHeader) error
}
