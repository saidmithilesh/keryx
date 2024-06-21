package connections

import (
	"keryx/utils"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	RegPrefix = "reg"
)

type HubRegistry struct {
	hubID  uuid.UUID
	config *utils.Config
	logger *zap.Logger

	stop chan struct{}
}

func (hr *HubRegistry) set() {
	now := time.Now().Unix()
	hr.logger.Info(
		"hub registry heartbeat",
		zap.String("hubId", h.ID.String()),
		zap.Int64("beatTimestamp", now),
	)
}

func (hr *HubRegistry) unset() {
	hr.logger.Info(
		"hub registry unset",
		zap.String("hubId", h.ID.String()),
	)
}

func (hr *HubRegistry) start() {
	interval := hr.config.RegistryInterval.Seconds()
	hr.logger.Info(
		"starting hub registry heartbeat",
		zap.String("hubId", h.ID.String()),
		zap.Int("registryInterval", int(interval)),
	)

	ticker := time.NewTicker(hr.config.RegistryInterval)
	for {
		select {
		case <-ticker.C:
			hr.set()
			continue
		case <-hr.stop:
			hr.unset()
			return
		}
	}
}

func NewHubRegistry(
	config *utils.Config,
	logger *zap.Logger,
	hubID uuid.UUID,
) *HubRegistry {
	return &HubRegistry{
		hubID:  hubID,
		config: config,
		logger: logger,
		stop:   make(chan struct{}),
	}
}
