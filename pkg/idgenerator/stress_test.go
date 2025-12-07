package idgenerator_test

import (
	"context"
	"errors"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/valkey-io/valkey-go"
	"github.com/yokeTH/backend/pkg/idgenerator"
)

const (
	numServices = 100
	idsPerSrv   = 100_000
	totalIDs    = numServices * idsPerSrv
)

func TestSnowflakeStress_100Services_10MIDs(t *testing.T) {
	t.Logf("starting stress test: %d services * %d ids = %d total",
		numServices, idsPerSrv, totalIDs)

	const valkeyURL = "redis://127.0.0.1:6379"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	ids := make([]uint64, totalIDs)

	var wg sync.WaitGroup
	errCh := make(chan error, numServices)

	for s := range numServices {
		srvIndex := s

		wg.Go(func() {

			svc, err := idgenerator.NewServiceFromValkey(ctx, valkeyURL)
			if err != nil {
				errCh <- err
				return
			}
			defer func() {
				if err := svc.Shutdown(ctx); err != nil {
					t.Errorf("service shutdown error: %v", err)
				}
			}()

			base := srvIndex * idsPerSrv
			for i := range idsPerSrv {
				ids[base+i] = svc.Generate()
			}
		})
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			t.Fatalf("service error: %v", err)
		}
	}

	if len(ids) != totalIDs {
		t.Fatalf("expected %d ids, got %d", totalIDs, len(ids))
	}

	t.Log("sorting IDs to check uniqueness...")
	slices.Sort(ids)

	for i := 1; i < len(ids); i++ {
		if ids[i] == ids[i-1] {
			t.Fatalf("found duplicate id at index %d: %d", i, ids[i])
		}
	}

	t.Log("stress test passed: 10M unique ids generated across 100 services")
}

func TestNewServiceFromValkey_HitsWorkerLimit(t *testing.T) {
	const valkeyURL = "redis://127.0.0.1:6379"

	opt, err := valkey.ParseURL(valkeyURL)
	if err != nil {
		t.Fatalf("parse valkey url: %v", err)
	}
	client, err := valkey.NewClient(opt)
	if err != nil {
		t.Fatalf("new valkey client: %v", err)
	}
	defer client.Close()

	if err := client.
		Do(context.Background(), client.B().Arbitrary("FLUSHDB").Build()).
		Error(); err != nil {
		t.Fatalf("flushdb: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	var services []*idgenerator.IDGenerator

	for i := range 1024 {
		svc, err := idgenerator.NewServiceFromValkey(ctx, valkeyURL)
		if err != nil {
			t.Fatalf("unexpected error creating service %d: %v", i, err)
		}
		services = append(services, svc)
	}

	defer func() {
		for _, svc := range services {
			_ = svc.Shutdown(context.Background())
		}
	}()

	if _, err := idgenerator.NewServiceFromValkey(ctx, valkeyURL); err == nil {
		t.Fatalf("expected error when exceeding worker-id limit, got nil")
	} else if !errors.Is(err, idgenerator.ErrNoFreeSlot) {
		t.Fatalf("expected ErrNoFreeSlot, got %v", err)
	}
}
