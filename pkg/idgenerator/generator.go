package idgenerator

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/yokeTH/backend/pkg/kv"
)

type IDGenerator struct {
	unit   *snowflakeUnit
	lease  *lease
	leaser *slotLeaser

	stopCh chan struct{}
}

func NewServiceFromValkey(ctx context.Context, connectionString string) (*IDGenerator, error) {
	v, err := kv.NewValkeyClient(connectionString)
	if err != nil {
		return nil, err
	}
	return NewService(ctx, v)
}

func NewService(ctx context.Context, store store) (*IDGenerator, error) {
	leaser := NewSlotLeaser(store)
	holder := uuid.NewString()

	lease, err := leaser.acquire(ctx, holder)
	if err != nil {
		return nil, err
	}

	gen := NewSnowflakeGenerator(lease.Slot)

	svc := &IDGenerator{
		unit:   gen,
		lease:  lease,
		leaser: leaser,
		stopCh: make(chan struct{}),
	}

	go svc.runHeartbeat(ctx)

	return svc, nil
}

func (s *IDGenerator) runHeartbeat(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ok, err := s.leaser.heartbeat(ctx, s.lease)
			if err != nil {
				log.Err(err).Msg("[id-svc] heartbeat error: %v")
				continue
			}
			if !ok {
				log.Info().Msg("[id-svc] lost lease; stopping heartbeat")
				return
			}
		case <-s.stopCh:
			return
		case <-ctx.Done():
			return
		}
	}
}

func (s *IDGenerator) Shutdown(ctx context.Context) error {
	close(s.stopCh)
	_ = s.leaser.release(ctx, s.lease)
	return nil
}

func (s *IDGenerator) Generate() uint64 {
	return s.unit.nextID()
}
