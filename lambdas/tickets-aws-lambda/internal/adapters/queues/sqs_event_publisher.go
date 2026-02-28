package queues

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/demoticketing/tickets/internal/core/domain"
)

type SQSEventPublisher struct {
	client   *sqs.Client
	queueURL string
}

func NewSQSEventPublisher(client *sqs.Client) *SQSEventPublisher {
	queueURL := os.Getenv("TICKET_QUEUE_URL")
	if queueURL == "" {
		queueURL = "http://localhost:4566/000000000000/ticket-purchased-queue"
	}
	return &SQSEventPublisher{client: client, queueURL: queueURL}
}

func (p *SQSEventPublisher) PublishTicketPurchased(ticket domain.Ticket) error {
	body, err := json.Marshal(ticket)
	if err != nil {
		return fmt.Errorf("failed to marshal ticket event: %w", err)
	}

	_, err = p.client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(p.queueURL),
		MessageBody: aws.String(string(body)),
	})
	if err != nil {
		return fmt.Errorf("failed to publish TicketPurchased to SQS: %w", err)
	}

	log.Printf("[SQS-PUBLISHER] TicketPurchased event published for ticket %s", ticket.ID)
	return nil
}
