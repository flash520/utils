package handler

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

type Handler struct {
}

func (h *Handler) Receive(msg string) error {
	log.Infof("RabbitMQ Consumer ReceiveData:\n %s\n", msg)
	return errors.New("failed")
	//return nil
}
