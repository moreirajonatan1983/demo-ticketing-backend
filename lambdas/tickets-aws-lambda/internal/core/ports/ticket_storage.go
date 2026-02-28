package ports

// TicketStorage is the output port to generate presigned download URLs for ticket PDFs.
// Implementations: S3PresignedURLGenerator (LocalStack/AWS S3)
type TicketStorage interface {
	GetPresignedDownloadURL(ticketId string) (string, error)
}
