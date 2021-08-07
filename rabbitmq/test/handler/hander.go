package handler

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

type Handler struct {
}

func (h *Handler) Receive(msg, msgType string) error {
	log.Infof("RabbitMQ Consumer ReceiveData:\n消息内容: %s\n消息类型: %s\n", msg, msgType)
	return errors.New("failed")
	//return nil
}
