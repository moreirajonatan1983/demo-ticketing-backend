package ports

import "github.com/demoticketing/tickets/internal/core/domain"

// EventPublisher is the output port to publish domain events to a message queue.
// Implementations: SQSEventPublisher (LocalStack/AWS SQS)
type EventPublisher interface {
	PublishTicketPurchased(ticket domain.Ticket) error
}
