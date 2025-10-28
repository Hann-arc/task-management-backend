package services

type EmailService interface {
	SendInvitation(to, token, projectName string) error
}

// NoopEmailService: dummy implementation for development
type NoopEmailService struct{}

func (n *NoopEmailService) SendInvitation(to, token, projectName string) error {
	// No operation performed
	return nil
}
