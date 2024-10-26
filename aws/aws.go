package aws

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func UploadPDF(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")
	region := os.Getenv("AWS_REGION")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKey,
			secretKey,
			"",
		),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create AWS session: %w", err)
	}

	uploader := s3.New(sess)
	buffer := make([]byte, fileHeader.Size)
	file.Read(buffer)
	fileName := fileHeader.Filename

	_, err = uploader.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String("pdfs/" + fileName),
		Body:          bytes.NewReader(buffer),
		ContentLength: aws.Int64(fileHeader.Size),
		ContentType:   aws.String("application/pdf"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload PDF to S3: %w", err)
	}

	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/pdfs/%s", bucketName, region, fileName)
	return fileURL, nil
}
