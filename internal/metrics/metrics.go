package metrics

import (
	"sync/atomic"
	"time"
)

type Metrics struct {
	StartTime time.Time

	RequestsTotal  atomic.Uint64
	ConnectTotal   atomic.Uint64
	ConnectErrors  atomic.Uint64
	UpstreamErrors atomic.Uint64

	BytesUp   atomic.Uint64
	BytesDown atomic.Uint64

	LastUpstreamRTTMs atomic.Int64
}

func New() *Metrics {
	return &Metrics{StartTime: time.Now()}
}

func (m *Metrics) NoteRTT(d time.Duration) {
	if d < 0 {
		return
	}
	m.LastUpstreamRTTMs.Store(d.Milliseconds())
}
