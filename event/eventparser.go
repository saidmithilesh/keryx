package event

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/saidmithilesh/keryx/config"
	"go.uber.org/zap"
)

type EventParser struct {
	cfg    *config.Config
	logger *zap.Logger
}

func (ep *EventParser) Parse(payload []byte) (*Event, error) {
	var event Event
	err := json.Unmarshal(payload, &event)
	if err != nil {
		ep.logger.Error("Failed to parse event", zap.Error(err))
		return nil, err
	}
	event.ID = uuid.New().String()
	event.ServerTimestamp = time.Now().Unix()
	return &event, nil
}

func (ep *EventParser) Serialize(event *Event) ([]byte, error) {
	payload, err := json.Marshal(event)
	if err != nil {
		ep.logger.Error("Failed to serialize event", zap.Error(err))
		return nil, err
	}
	return payload, nil
}

func NewParser(cfg *config.Config, logger *zap.Logger) *EventParser {
	return &EventParser{
		cfg:    cfg,
		logger: logger,
	}
}
