package comms

import (
	"keryx/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	StatusCodeInvalidPacket = 90001
	StatusCodeInvalidSender = 90002
	StatusCodeInvalidRoomID = 90003
)

type RoutingStatus struct {
	EventID    uuid.UUID
	Success    bool
	StatusCode int
	Reason     string
}

func invalidPacketError(event *Event) *RoutingStatus {
	return &RoutingStatus{
		EventID:    event.ID,
		Success:    false,
		StatusCode: StatusCodeInvalidPacket,
		Reason:     "Failed to parse invalid packet",
	}
}

type Router struct {
	config *utils.Config
	logger *zap.Logger

	inputChan  chan []byte
	outputChan chan *RoutingStatus
}

func (router *Router) Listen() {
	for {
		packet := <-router.inputChan
		event, err := EventFromBytes(packet)
		if err != nil {
			router.logger.Error("failed to deserialize data on channel")
			router.outputChan <- invalidPacketError(event)
		}
	}
}

func NewRouter(
	config *utils.Config,
	logger *zap.Logger,
	inputChan chan []byte,
	outputChan chan *RoutingStatus,
) (router *Router) {
	return &Router{
		config:     config,
		logger:     logger,
		inputChan:  inputChan,
		outputChan: outputChan,
	}
}
