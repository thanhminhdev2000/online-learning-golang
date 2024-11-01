package utils

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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

func UploadVideo(cld *cloudinary.Cloudinary, file multipart.File) (string, int, error) {
	// Create a temporary copy of the file to get the video duration
	tempFile, err := os.CreateTemp("", "temp-video-*.mp4")
	if err != nil {
		return "", 0, fmt.Errorf("unable to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Copy the content into the temporary file
	_, err = io.Copy(tempFile, file)
	if err != nil {
		return "", 0, fmt.Errorf("unable to copy content to temporary file: %v", err)
	}
	tempFile.Close()

	// Get the video duration
	cmd := exec.Command("ffprobe", "-i", tempFile.Name(), "-show_entries", "format=duration", "-v", "quiet", "-of", "csv=p=0")
	output, err := cmd.Output()
	if err != nil {
		return "", 0, fmt.Errorf("ffprobe error: %v", err)
	}

	duration, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		return "", 0, fmt.Errorf("unable to convert duration: %v", err)
	}

	// Reset the pointer of the original file to the beginning for upload
	_, err = file.Seek(0, 0)
	if err != nil {
		return "", 0, fmt.Errorf("unable to reset file pointer: %v", err)
	}

	// Upload video
	ctx := context.Background()
	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: "videos",
	})
	if err != nil {
		return "", 0, fmt.Errorf("failed to upload video: %v", err)
	}

	return uploadResult.SecureURL, int(duration), nil
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

func GetVideoDuration(file multipart.File) (int, error) {
	tempFile, err := os.CreateTemp("", "temp-video-*.mp4")
	if err != nil {
		return 0, fmt.Errorf("không thể tạo file tạm thời: %v", err)
	}
	defer os.Remove(tempFile.Name())

	_, err = io.Copy(tempFile, file)
	if err != nil {
		return 0, fmt.Errorf("không thể sao chép nội dung vào file tạm: %v", err)
	}
	tempFile.Close()

	// Sử dụng ffprobe để lấy thông tin thời lượng video
	cmd := exec.Command("ffprobe", "-i", tempFile.Name(), "-show_entries", "format=duration", "-v", "quiet", "-of", "csv=p=0")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe lỗi: %v", err)
	}

	// Chuyển đổi thời lượng sang int
	duration, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		return 0, fmt.Errorf("không thể chuyển đổi thời lượng: %v", err)
	}

	return int(duration), nil
}
