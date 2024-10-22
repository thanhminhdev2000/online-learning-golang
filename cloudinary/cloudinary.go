package cloudinary

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func SetupCloudinary() (*cloudinary.Cloudinary, error) {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary: %v", err)
		return nil, err
	}
	return cld, nil
}

func UploadAvatar(cld *cloudinary.Cloudinary, fileContent io.Reader) (string, error) {
	uploadResult, err := cld.Upload.Upload(context.Background(), fileContent, uploader.UploadParams{})
	if err != nil {
		log.Fatalf("Failed to upload image: %v", err)
		return "", err
	}
	return uploadResult.SecureURL, nil
}
