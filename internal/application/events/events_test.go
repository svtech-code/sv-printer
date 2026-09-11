package events

import (
	"encoding/json"
	"testing"
)

func TestBusPublishSubscribe(t *testing.T) {
	b := NewBus()
	ch, cancel := b.Subscribe()
	defer cancel()

	b.Publish(NewEvent(PrintStarted, map[string]any{"job_id": "job1"}))

	ev := <-ch
	if ev.Type != PrintStarted {
		t.Errorf("Type = %q, want %q", ev.Type, PrintStarted)
	}
	if ev.Data["job_id"] != "job1" {
		t.Errorf("Data[job_id] = %v, want job1", ev.Data["job_id"])
	}
}

func TestBusMultipleSubscribers(t *testing.T) {
	b := NewBus()
	ch1, cancel1 := b.Subscribe()
	defer cancel1()
	ch2, cancel2 := b.Subscribe()
	defer cancel2()

	b.Publish(NewEvent(PrintCompleted, nil))

	if ev := <-ch1; ev.Type != PrintCompleted {
		t.Errorf("subscriber 1 got %q", ev.Type)
	}
	if ev := <-ch2; ev.Type != PrintCompleted {
		t.Errorf("subscriber 2 got %q", ev.Type)
	}
}

func TestBusCancelStopsDelivery(t *testing.T) {
	b := NewBus()
	ch, cancel := b.Subscribe()
	cancel()

	b.Publish(NewEvent(PrintStarted, nil))

	select {
	case _, ok := <-ch:
		if ok {
			t.Error("closed channel should not deliver events")
		}
	default:
	}
}

func TestEventMarshalJSON(t *testing.T) {
	ev := NewEvent(PrintStarted, map[string]any{"job_id": "job1", "printer_id": "p1"})

	data, err := json.Marshal(ev)
	if err != nil {
		t.Fatalf("MarshalJSON: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if m["event"] != "print.started" {
		t.Errorf("event = %v, want print.started", m["event"])
	}
	if m["job_id"] != "job1" {
		t.Errorf("job_id = %v, want job1", m["job_id"])
	}
}
