package uploading

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/go-faster/errors"
)

func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return errors.Wrap(err, "SaveUploadedFile file.Open:")
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return errors.Wrap(err, "SaveUploadedFile os.MkdirAll:")
	}

	out, err := os.Create(dst)
	if err != nil {
		return errors.Wrap(err, "SaveUploadedFile os.Create:")
	}
	defer out.Close()
	_, err = io.Copy(out, src)
	return errors.Wrap(err, "SaveUploadedFile io.Copy:")
}
