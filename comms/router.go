package comms

import (
	"fmt"
	"keryx/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	StatusCodeInvalidPacket    = 90001
	StatusCodeInvalidSender    = 90002
	StatusCodeInvalidEventType = 90003
)

type RoutingStatus struct {
	EventID    uuid.UUID
	SenderID   uuid.UUID
	Success    bool
	StatusCode int
	Reason     string
}

func invalidPacketError(event *Event) *RoutingStatus {
	return &RoutingStatus{
		EventID:    event.ID,
		SenderID:   event.SenderID,
		Success:    false,
		StatusCode: StatusCodeInvalidPacket,
		Reason:     "Failed to parse invalid packet",
	}
}

func invalidEventTypeError(event *Event) *RoutingStatus {
	return &RoutingStatus{
		EventID:    event.ID,
		SenderID:   event.SenderID,
		Success:    false,
		StatusCode: StatusCodeInvalidEventType,
		Reason:     fmt.Sprintf("Invalid event type %s", event.Type),
	}
}

type Router struct {
	config *utils.Config
	logger *zap.Logger

	outputChan chan *RoutingStatus
}

func (router *Router) Route(packet []byte) {
	event, err := EventFromBytes(packet)
	if err != nil {
		router.logger.Error("failed to deserialize data on channel")
		router.outputChan <- invalidPacketError(event)
		return
	}

	err = event.validate()
	if err != nil {
		router.logger.Error("event validation failed", zap.Error(err))
		router.outputChan <- invalidEventTypeError(event)
	}
}

func NewRouter(
	config *utils.Config,
	logger *zap.Logger,
	outputChan chan *RoutingStatus,
) (router *Router) {
	return &Router{
		config:     config,
		logger:     logger,
		outputChan: outputChan,
	}
}
