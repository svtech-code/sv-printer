package events

import (
	"encoding/json"
	"sync"
)

const (
	PrinterDiscovered    = "printer.discovered"
	PrinterRemoved       = "printer.removed"
	PrinterStatusChanged = "printer.status_changed"
	PrintStarted         = "print.started"
	PrintCompleted       = "print.completed"
	PrintFailed          = "print.failed"
	AgentStatus          = "agent.status"
)

type Event struct {
	Type string
	Data map[string]any
}

func NewEvent(t string, data map[string]any) Event {
	return Event{Type: t, Data: data}
}

func (e Event) MarshalJSON() ([]byte, error) {
	m := make(map[string]any, len(e.Data)+1)
	for k, v := range e.Data {
		m[k] = v
	}
	m["event"] = e.Type
	return json.Marshal(m)
}

type Bus struct {
	mu   sync.RWMutex
	subs map[chan Event]struct{}
}

func NewBus() *Bus {
	return &Bus{subs: make(map[chan Event]struct{})}
}

func (b *Bus) Subscribe() (<-chan Event, func()) {
	ch := make(chan Event, 64)
	b.mu.Lock()
	b.subs[ch] = struct{}{}
	b.mu.Unlock()

	cancel := func() {
		b.mu.Lock()
		if _, ok := b.subs[ch]; ok {
			delete(b.subs, ch)
			close(ch)
		}
		b.mu.Unlock()
	}
	return ch, cancel
}

func (b *Bus) Publish(e Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for ch := range b.subs {
		select {
		case ch <- e:
		default:
		}
	}
}
