package config

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DatabaseURL string
	onceDB      sync.Once
}

type MinIO struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	client    *s3.Client
	onceMinio sync.Once
}

func NewConfig() (*Config, *MinIO, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=postgres user=postgres password=secret dbname=userapp port=5432 sslmode=disable"
	}

	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	if minioEndpoint == "" {
		minioEndpoint = "http://minio:9000"
	}
	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	if minioAccessKey == "" {
		minioAccessKey = "minioadmin"
	}
	minioSecretKey := os.Getenv("MINIO_SECRET_KEY")
	if minioSecretKey == "" {
		minioSecretKey = "minioadmin"
	}
	minioBucket := os.Getenv("MINIO_BUCKET")
	if minioBucket == "" {
		minioBucket = "recipe-images"
	}

	return &Config{DatabaseURL: dsn}, &MinIO{
		Endpoint:  minioEndpoint,
		AccessKey: minioAccessKey,
		SecretKey: minioSecretKey,
		Bucket:    minioBucket,
	}, nil
}

func NewDB(cfg *Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	cfg.onceDB.Do(func() {
		db, err = gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (m *MinIO) GetClient() (*s3.Client, error) {
	var err error
	m.onceMinio.Do(func() {
		cfg, errCfg := config.LoadDefaultConfig(context.Background(),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(m.AccessKey, m.SecretKey, "")),
			config.WithEndpointResolverWithOptions(
				aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:           m.Endpoint,
						SigningRegion: "us-east-1",
						Source:        aws.EndpointSourceCustom,
					}, nil
				}),
			),
			config.WithRegion("us-east-1"),
		)
		if errCfg != nil {
			err = errCfg
			return
		}

		m.client = s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.UsePathStyle = true // Required for MinIO
		})
	})

	if err != nil {
		return nil, err
	}
	if m.client == nil {
		return nil, fmt.Errorf("failed to initialize MinIO client")
	}
	return m.client, nil
}
