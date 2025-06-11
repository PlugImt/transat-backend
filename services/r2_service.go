package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type R2Service struct {
	client     *s3.Client
	bucketName string
	publicURL  string
}

func NewR2Service() (*R2Service, error) {
	log.Println("Initializing R2 (storage) service...")
	accountID := os.Getenv("R2_ACCOUNT_ID")
	accessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("R2_ACCESS_KEY_SECRET")
	bucketName := os.Getenv("R2_BUCKET_NAME")
	publicURL := os.Getenv("R2_PUBLIC_DOMAIN")

	if accountID == "" || accessKeyID == "" || accessKeySecret == "" || bucketName == "" || publicURL == "" {
		missingVars := []string{}
		if accountID == "" {
			missingVars = append(missingVars, "R2_ACCOUNT_ID")
		}
		if accessKeyID == "" {
			missingVars = append(missingVars, "R2_ACCESS_KEY_ID")
		}
		if accessKeySecret == "" {
			missingVars = append(missingVars, "R2_ACCESS_KEY_SECRET")
		}
		if bucketName == "" {
			missingVars = append(missingVars, "R2_BUCKET_NAME")
		}
		if publicURL == "" {
			missingVars = append(missingVars, "R2_PUBLIC_DOMAIN")
		}
		return nil, fmt.Errorf("missing required R2 configuration: %v", missingVars)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, accessKeySecret, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	// eu.r2 car on est dans l'UE (RGPD etc)
	endpoint := fmt.Sprintf("https://%s.eu.r2.cloudflarestorage.com", accountID)
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	log.Println("✅ R2 service initialized successfully")
	return &R2Service{
		client:     client,
		bucketName: bucketName,
		publicURL:  publicURL,
	}, nil
}

func (s *R2Service) UploadFile(key string, reader io.Reader, contentType string) (string, error) {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String(contentType),
	}

	_, err := s.client.PutObject(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return s.GetPublicURL(key), nil
}

func (s *R2Service) DeleteFile(key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}

	_, err := s.client.DeleteObject(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func (s *R2Service) GetPublicURL(key string) string {
	publicURL := strings.TrimRight(s.publicURL, "/")
	return fmt.Sprintf("%s/%s", publicURL, key)
}

func (s *R2Service) GetPresignedURL(key string, duration time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s.client)

	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}

	request, err := presignClient.PresignGetObject(context.TODO(), input, s3.WithPresignExpires(duration))
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return request.URL, nil
}

func (s *R2Service) GetObject(key string) (io.ReadCloser, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}

	result, err := s.client.GetObject(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	return result.Body, nil
}
