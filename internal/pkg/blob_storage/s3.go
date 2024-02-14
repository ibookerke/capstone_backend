package blob_storage

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type BlobStorage interface {
	Upload(ctx context.Context, filename, contentType string, file []byte) (string, error)
}

type BlobStorageS3 struct {
	bucketName string
	cfg        aws.Config
	uploader   *manager.Uploader
}

func NewBlobStorageS3(bucketName, awsRegion string) (*BlobStorageS3, error) {
	svc := &BlobStorageS3{
		bucketName: bucketName,
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(awsRegion))

	if err != nil {
		return nil, fmt.Errorf("service error during loading aws config: %w", err)
	}

	svc.cfg = cfg
	svc.uploader = manager.NewUploader(s3.NewFromConfig(cfg))

	return svc, nil
}

func (s *BlobStorageS3) Upload(ctx context.Context, filename, contentType string, file []byte) (string, error) {
	awsFile := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(filename),
		Body:        bytes.NewReader(file),
		ContentType: aws.String(contentType),
	}

	object, err := s.uploader.Upload(ctx, awsFile)

	if err != nil {
		return "", fmt.Errorf("uploading to s3: %w", err)
	}

	if object == nil {
		return "", fmt.Errorf("uploading to s3Service: object is nil")
	}

	return object.Location, nil
}
