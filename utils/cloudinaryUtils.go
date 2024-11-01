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

func UploadImage(cld *cloudinary.Cloudinary, file multipart.File) (string, error) {
	ctx := context.Background()
	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: "images",
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

func UploadVideo(cld *cloudinary.Cloudinary, file multipart.File) (string, error) {
	ctx := context.Background()
	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: "videos",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload video: %v", err)
	}
	return uploadResult.SecureURL, nil
}

func DeleteVideo(cld *cloudinary.Cloudinary, videoURL string) error {
	if videoURL == "" {
		return nil
	}

	parts := strings.Split(videoURL, "/")
	if len(parts) < 2 {
		return fmt.Errorf("invalid video URL format")
	}
	filename := parts[len(parts)-1]
	publicID := "videos/" + strings.TrimSuffix(filename, filepath.Ext(filename))

	ctx := context.Background()
	_, err := cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID})
	if err != nil {
		return fmt.Errorf("failed to delete video: %v", err)
	}

	return nil
}
