package cloudinary

import (
	"context"
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func SetupCloudinary() (*cloudinary.Cloudinary, error) {
	cld, err := cloudinary.NewFromParams("your-cloud-name", "your-api-key", "your-api-secret")
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary: %v", err)
		return nil, err
	}
	return cld, nil
}

func UploadAvatar(cld *cloudinary.Cloudinary, filePath string) (string, error) {
	uploadResult, err := cld.Upload.Upload(context.Background(), filePath, uploader.UploadParams{})
	if err != nil {
		log.Fatalf("Failed to upload image: %v", err)
		return "", err
	}
	return uploadResult.SecureURL, nil
}
