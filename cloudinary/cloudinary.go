package cloudinary

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
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

func UploadImage(cld *cloudinary.Cloudinary, fileContent io.Reader) (string, error) {
	uploadParams := uploader.UploadParams{
		Folder:   "images",
		PublicID: "",
	}

	uploadResult, err := cld.Upload.Upload(context.Background(), fileContent, uploadParams)
	if err != nil {
		log.Fatalf("Failed to upload image: %v", err)
		return "", err
	}
	return uploadResult.SecureURL, nil
}

func UploadVideo(cld *cloudinary.Cloudinary, fileContent io.Reader) (string, error) {
	buffer := make([]byte, 512)
	if _, err := fileContent.Read(buffer); err != nil {
		log.Println("Failed to read file")
		return "", err
	}

	filetype := http.DetectContentType(buffer)
	if filetype != "video/mp4" && filetype != "video/avi" && filetype != "video/mkv" {
		log.Println("File type is not supported. Please upload a video file.")
		return "", fmt.Errorf("unsupported file type: %s", filetype)
	}

	if seeker, ok := fileContent.(io.Seeker); ok {
		seeker.Seek(0, io.SeekStart)
	} else {
		return "", fmt.Errorf("file content is not seekable")
	}

	uploadParams := uploader.UploadParams{
		Folder:       "videos",
		PublicID:     "",
		ResourceType: "video",
	}

	uploadResult, err := cld.Upload.Upload(context.Background(), fileContent, uploadParams)
	if err != nil {
		log.Fatalf("Failed to upload video: %v", err)
		return "", err
	}

	return uploadResult.SecureURL, nil
}
