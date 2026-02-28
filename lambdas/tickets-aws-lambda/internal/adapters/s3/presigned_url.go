package s3

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3PresignedURLGenerator generates presigned GET URLs for ticket PDFs stored in S3.
// PDF key convention: tickets/{ticketId}.pdf  (written by ticket-worker Java service)
type S3PresignedURLGenerator struct {
	presignClient *s3.PresignClient
	bucket        string
}

func NewS3PresignedURLGenerator(client *s3.Client) *S3PresignedURLGenerator {
	bucket := os.Getenv("TICKETS_BUCKET")
	if bucket == "" {
		bucket = "ticketera-tickets-local"
	}
	return &S3PresignedURLGenerator{
		presignClient: s3.NewPresignClient(client),
		bucket:        bucket,
	}
}

func (g *S3PresignedURLGenerator) GetPresignedDownloadURL(ticketId string) (string, error) {
	key := fmt.Sprintf("tickets/%s.pdf", ticketId)

	req, err := g.presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket:                     aws.String(g.bucket),
		Key:                        aws.String(key),
		ResponseContentDisposition: aws.String(fmt.Sprintf("attachment; filename=\"ticket-%s.pdf\"", ticketId)),
		ResponseContentType:        aws.String("application/pdf"),
	}, func(o *s3.PresignOptions) {
		o.Expires = 15 * time.Minute
	})

	// Suppress unused import for types
	_ = types.ObjectCannedACLPrivate

	if err != nil {
		return "", fmt.Errorf("failed to presign URL for ticket %s: %w", ticketId, err)
	}
	return req.URL, nil
}
