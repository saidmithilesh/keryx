package comms

import (
	"fmt"
	"keryx/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	StatusCodeSucess           = 90000
	StatusCodeInvalidPacket    = 90001
	StatusCodeInvalidSender    = 90002
	StatusCodeInvalidEventType = 90003
)

type RoutingStatus struct {
	EventID    uuid.UUID `json:"eventId"`
	SenderID   uuid.UUID `json:"senderId"`
	Success    bool      `json:"success"`
	StatusCode int       `json:"statusCode"`
	Reason     string    `json:"reason"`
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

func routingSuccess(event *Event) *RoutingStatus {
	return &RoutingStatus{
		EventID:    event.ID,
		SenderID:   event.SenderID,
		Success:    true,
		StatusCode: StatusCodeSucess,
		Reason:     "sent successfully",
	}
}

type Router struct {
	config *utils.Config
	logger *zap.Logger

	Output chan *RoutingStatus
}

func (router *Router) Route(packet []byte) {
	event, err := EventFromBytes(packet)
	if err != nil {
		router.logger.Error(
			"failed to deserialize data on channel",
			zap.Error(err),
		)
		router.Output <- invalidPacketError(event)
		return
	}

	err = event.validate()
	if err != nil {
		router.logger.Error("event validation failed", zap.Error(err))
		router.Output <- invalidEventTypeError(event)
	}

	router.Output <- routingSuccess(event)
}

func NewRouter(
	config *utils.Config,
	logger *zap.Logger,
) (router *Router) {
	return &Router{
		config: config,
		logger: logger,
		Output: make(chan *RoutingStatus, 256),
	}
}
