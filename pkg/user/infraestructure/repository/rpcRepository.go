package repository

import (
	"errors"
	"net/rpc"

	"github.com/jeffleon/oauth-microservice/pkg/user/domain"
	"github.com/sirupsen/logrus"
)

type rpcRepository struct {
	clientRPC *rpc.Client
}

func NewRPCRepository(clientRPC *rpc.Client) domain.RPCRepository {
	return &rpcRepository{
		clientRPC,
	}
}

func (r *rpcRepository) SendEmail(rpcPayload domain.RPCPayload) (string, error) {
	var result string
	if r.clientRPC == nil {
		return "", errors.New("rpc client is not allowed to send emails")
	}

	if err := r.clientRPC.Call("Mailer.RPCSendEmail", rpcPayload, &result); err != nil {
		logrus.Errorf("Error, send to Mailer.RPCSendEmail %s", err)
		return "", err
	}

	return result, nil
}
