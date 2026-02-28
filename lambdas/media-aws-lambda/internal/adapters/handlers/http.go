package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type HTTPHandler struct {
	client     *s3.Client
	presignCl  *s3.PresignClient
	bucketName string
}

func NewHTTPHandler(client *s3.Client, bucketName string) *HTTPHandler {
	presignClient := s3.NewPresignClient(client)
	return &HTTPHandler{
		client:     client,
		presignCl:  presignClient,
		bucketName: bucketName,
	}
}

func (h *HTTPHandler) HandleHTTPRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	eventId := req.PathParameters["eventId"]

	if eventId == "" {
		return h.jsonResponse(http.StatusBadRequest, `{"error":"Missing eventId"}`)
	}

	// Calculate S3 object key based on business rule (matches seed script)
	objectKey := fmt.Sprintf("events/%s/cover.jpg", eventId)

	// Create Presigned URL (valid for 15 minutes)
	reqInput := &s3.GetObjectInput{
		Bucket: aws.String(h.bucketName),
		Key:    aws.String(objectKey),
	}

	presignedReq, err := h.presignCl.PresignGetObject(context.TODO(), reqInput, func(opts *s3.PresignOptions) {
		opts.Expires = 15 * time.Minute
	})

	if err != nil {
		return h.jsonResponse(http.StatusInternalServerError, `{"error":"Failed to sign URL"}`)
	}

	redirectUrl := presignedReq.URL

	// Because of LocalStack docker networking, we need to map the internal docker host (localhost/4566)
	// back out for the browser. The browser cannot resolve host.docker.internal properly.
	importUrl := strings.Replace(redirectUrl, "host.docker.internal", "localhost", 1)

	// Return the presigned URL as a 302 redirect so <img src="..."> works properly
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusFound,
		Headers: map[string]string{
			"Location":                     importUrl,
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET, OPTIONS",
		},
	}, nil
}

func (h *HTTPHandler) jsonResponse(status int, body string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       body,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}
