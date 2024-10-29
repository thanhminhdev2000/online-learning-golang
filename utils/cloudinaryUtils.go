package utils

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var (
	allowedImageTypes = map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}
)

func SetupCloudinary() (*cloudinary.Cloudinary, error) {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("missing required Cloudinary environment variables")
	}

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Cloudinary: %v", err)
	}

	return cld, nil
}

func ValidateImageType(file *multipart.FileHeader) error {
	if !allowedImageTypes[file.Header.Get("Content-Type")] {
		return fmt.Errorf("invalid file type. Only JPEG, PNG and GIF are allowed")
	}
	return nil
}

func UploadImage(cld *cloudinary.Cloudinary, file multipart.File) (string, error) {
	ctx := context.Background()
	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: "courses",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %v", err)
	}
	return uploadResult.SecureURL, nil
}

func DeleteImage(cld *cloudinary.Cloudinary, imageURL string) error {
	if imageURL == "" {
		return nil
	}

	// Extract public ID from URL
	parts := strings.Split(imageURL, "/")
	if len(parts) < 2 {
		return fmt.Errorf("invalid image URL format")
	}
	filename := parts[len(parts)-1]
	publicID := "courses/" + strings.TrimSuffix(filename, filepath.Ext(filename))

	ctx := context.Background()
	_, err := cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID})
	if err != nil {
		return fmt.Errorf("failed to delete image: %v", err)
	}

	return nil
} 