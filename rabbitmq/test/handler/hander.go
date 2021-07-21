package handler

import (
	"errors"
	"fmt"
)

type Handler struct {
}

func (h *Handler) Receive(msg string) error {
	fmt.Printf("\rRabbitMQ Consumer ReceiveData: %s", msg)
	return errors.New("failed")
	//return nil
}
