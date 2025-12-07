package idgenerator

import (
	"context"
	"fmt"
	"time"
)

type store interface {
	SetNX(ctx context.Context, key, val string, ttl time.Duration) (bool, error)
	Get(ctx context.Context, key string) (string, error)
	Expire(ctx context.Context, key string, ttl time.Duration) (bool, error)
	CasDel(ctx context.Context, key, expected string) (int, error)
}

const DefaultTTLSecs int64 = 30

type lease struct {
	Key     string
	Slot    uint16
	Holder  string
	TTLSecs int64
}

type slotLeaser struct {
	Store   store
	Prefix  string
	TTLSecs int64
	MaxSlot uint16
}

func NewSlotLeaser(store store) *slotLeaser {
	return &slotLeaser{
		Store:   store,
		Prefix:  "snowflake:wid:",
		TTLSecs: DefaultTTLSecs,
		MaxSlot: 1023,
	}
}

func (s *slotLeaser) keyFor(slot uint16) string {
	return s.Prefix + strconvU16(slot)
}

func (s *slotLeaser) acquire(ctx context.Context, holder string) (*lease, error) {
	nanos := time.Now().UnixNano()
	start := uint16(uint32(nanos) % uint32(s.MaxSlot+1))

	for i := uint16(0); i <= s.MaxSlot; i++ {
		slot := (start + i) & s.MaxSlot
		key := s.keyFor(slot)

		ok, err := s.Store.SetNX(ctx, key, holder, time.Duration(s.TTLSecs)*time.Second)
		if err != nil {
			return nil, err
		}
		if ok {
			return &lease{
				Key:     key,
				Slot:    slot,
				Holder:  holder,
				TTLSecs: s.TTLSecs,
			}, nil
		}
	}
	return nil, ErrNoFreeSlot
}

var ErrNoFreeSlot = fmt.Errorf("no free slot available")

// heartbeat extends lease TTL if we still own the key.
func (s *slotLeaser) heartbeat(ctx context.Context, lease *lease) (bool, error) {
	val, err := s.Store.Get(ctx, lease.Key)
	if err != nil {
		return false, err
	}
	if val != lease.Holder {
		return false, nil
	}
	return s.Store.Expire(ctx, lease.Key, time.Duration(lease.TTLSecs)*time.Second)
}

// release performs CAS delete: DEL key only if value matches holder.
func (s *slotLeaser) release(ctx context.Context, lease *lease) error {
	_, err := s.Store.CasDel(ctx, lease.Key, lease.Holder)
	return err
}
