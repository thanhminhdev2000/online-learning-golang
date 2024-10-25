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
	uploadParams := uploader.UploadParams{
		Folder:   "avatar",
		PublicID: "",
	}

	uploadResult, err := cld.Upload.Upload(context.Background(), fileContent, uploadParams)
	if err != nil {
		log.Fatalf("Failed to upload image: %v", err)
		return "", err
	}
	return uploadResult.SecureURL, nil
}

func UploadDocument(cld *cloudinary.Cloudinary, fileContent io.Reader) (string, error) {
	uploadParams := uploader.UploadParams{
		Folder:   "documentation",
		PublicID: "",
	}

	uploadResult, err := cld.Upload.Upload(context.Background(), fileContent, uploadParams)
	if err != nil {
		log.Printf("Failed to upload document: %v", err)
		return "", err
	}

	return uploadResult.SecureURL, nil
}
