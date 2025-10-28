// utils/cloudinary.go
package utils

import (
	"context"
	"mime/multipart"
	"strings"

	"github.com/Hann-arc/task-management-backend/config"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadToCloudinary(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	uploadResult, err := config.Cld.Upload.Upload(
		context.Background(),
		src,
		uploader.UploadParams{
			Folder: "MgApp",
		},
	)
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}

func DeleteFromCloudinary(fileUrl string) error {

	parts := strings.Split(fileUrl, "/upload/")
	if len(parts) < 2 {
		return nil
	}

	publicIDWithExt := parts[1]
	publicID := strings.Split(publicIDWithExt, ".")[0]

	_, err := config.Cld.Admin.DeleteAssets(
		context.Background(),
		admin.DeleteAssetsParams{
			PublicIDs: []string{publicID},
		},
	)
	return err
}
