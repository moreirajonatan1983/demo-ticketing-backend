package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/demoticketing/media/internal/adapters/handlers"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	var client *s3.Client
	endpoint := os.Getenv("AWS_ENDPOINT_URL")
	if endpoint != "" {
		client = s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(endpoint)
			o.UsePathStyle = true
		})
	} else {
		client = s3.NewFromConfig(cfg)
	}

	bucketName := os.Getenv("MEDIA_BUCKET_NAME")
	if bucketName == "" {
		bucketName = "ticketera-images-local"
	}

	handler := handlers.NewHTTPHandler(client, bucketName)
	lambda.Start(handler.HandleHTTPRequest)
}
