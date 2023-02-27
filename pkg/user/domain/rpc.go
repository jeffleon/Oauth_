package domain

import "github.com/jeffleon/oauth-microservice/internal/config"

type RPCRepository interface {
	SendEmail(rpcPayload RPCPayload) (string, error)
}

type RPCPayload struct {
	From     string `json:"from"`
	FromName string `json:"from_name"`
	To       string `json:"to"`
	Subject  string `json:"subject"`
	Message  string `json:"message"`
}

func NewRPCPayloadEmailWelcome(email string) RPCPayload {
	return RPCPayload{
		From:     config.Config.EmailFrom,
		FromName: config.Config.EmailFromName,
		To:       email,
		Subject:  config.EmailWelcomeSubject,
		Message:  config.EmailWelcomeMsg,
	}
}
