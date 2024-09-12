package router

import (
	"github.com/saidmithilesh/keryx/config"
	"github.com/saidmithilesh/keryx/event"
	"github.com/saidmithilesh/keryx/hub"
	"go.uber.org/zap"
)

type Router struct {
	cfg         *config.Config
	logger      *zap.Logger
	wsHub       *hub.Hub
	eventParser *event.EventParser
}

func (r *Router) onMessage(message []byte, userID string) {
	event, err := r.eventParser.Parse(message)
	if err != nil {
		r.logger.Error("Failed to parse event", zap.Error(err))
		return
	}
	r.logger.Info("Received message", zap.String("message", string(event.ID)))
	payload, err := r.eventParser.Serialize(event)
	if err != nil {
		r.logger.Error("Failed to serialize event", zap.Error(err))
		return
	}
	r.wsHub.Multicast([]string{userID}, payload)
}

func New(cfg *config.Config, logger *zap.Logger, wsHub *hub.Hub, eventParser *event.EventParser) *Router {
	router := &Router{
		cfg:         cfg,
		logger:      logger,
		wsHub:       wsHub,
		eventParser: eventParser,
	}

	wsHub.SetMessageListener(router.onMessage)
	return router
}
