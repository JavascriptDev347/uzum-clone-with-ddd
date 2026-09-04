package media

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Uploader struct {
	client *s3.Client
	bucket string
	region string
	folder string // masalan: "uzum-clone" - barcha fayllar shu papkada saqlanadi
}

type S3Config struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	Bucket          string
	Folder          string
}

func NewS3Uploader(cfg S3Config) (*S3Uploader, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("s3 init error: %w", err)
	}

	return &S3Uploader{
		client: s3.NewFromConfig(awsCfg),
		bucket: cfg.Bucket,
		region: cfg.Region,
		folder: cfg.Folder,
	}, nil
}

func (u *S3Uploader) key(fileName string) string {
	if u.folder == "" {
		return fileName
	}
	return u.folder + "/" + fileName
}

func (u *S3Uploader) Upload(ctx context.Context, input UploadInput) (UploadResult, error) {
	key := u.key(input.FileName)

	_, err := u.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(u.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(input.Data),
		ContentType: aws.String(input.ContentType),
	})
	if err != nil {
		return UploadResult{}, fmt.Errorf("%w: %v", ErrUploadFailed, err)
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", u.bucket, u.region, key)

	return UploadResult{
		URL:      url,
		PublicID: key,
	}, nil
}

func (u *S3Uploader) Delete(ctx context.Context, publicID string) error {
	_, err := u.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(u.bucket),
		Key:    aws.String(publicID),
	})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteFailed, err)
	}
	return nil
}
