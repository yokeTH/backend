package idgenerator

import (
	"sync"
	"time"
)

const (
	epochMs    int64  = 1_044_230_400_000 // 2003-02-03T00:00:00Z
	workerBits uint64 = 10
	seqBits    uint64 = 12

	maxWorker uint16 = uint16((1 << workerBits) - 1) // 1023
	maxSeq    uint16 = uint16((1 << seqBits) - 1)    // 4095

	workerShift uint64 = seqBits
	tsShift     uint64 = workerBits + seqBits
)

type snowflakeState struct {
	lastTs   int64
	seq      uint16
	workerID uint16
}

type snowflakeUnit struct {
	mu    sync.Mutex
	state snowflakeState
}

func NewSnowflakeGenerator(workerID uint16) *snowflakeUnit {
	if workerID > maxWorker {
		panic("workerID out of range")
	}
	return &snowflakeUnit{
		state: snowflakeState{
			lastTs:   -1,
			seq:      0,
			workerID: workerID,
		},
	}
}

func nowMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func waitNextMs(last int64) int64 {
	for {
		n := nowMs()
		if n > last {
			return n
		}
	}
}

func (g *snowflakeUnit) nextID() uint64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	wall := nowMs()
	if wall < g.state.lastTs {
		wall = waitNextMs(g.state.lastTs)
	}

	if wall == g.state.lastTs {
		if g.state.seq == maxSeq {
			next := waitNextMs(g.state.lastTs)
			g.state.lastTs = next
			g.state.seq = 0
		} else {
			g.state.seq++
		}
	} else {
		g.state.seq = 0
		g.state.lastTs = wall
	}

	ts := uint64(g.state.lastTs - epochMs)
	return (ts << tsShift) |
		(uint64(g.state.workerID) << workerShift) |
		uint64(g.state.seq)
}

//lint:ignore U1000 exported for external use
func (g *snowflakeUnit) workerID() uint16 {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.state.workerID
}
