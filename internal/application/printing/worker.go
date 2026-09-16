package printing

import (
	"context"
	"time"

	"sv-printer/internal/application/events"
	"sv-printer/internal/domain/job"
	"sv-printer/internal/infrastructure/transport"
)

type TransportFactory func(printerID string) (transport.PrinterTransport, error)

type Worker struct {
	queue            PrintQueue
	transportFactory TransportFactory
	stopCh           chan struct{}
	bus              *events.Bus
}

func NewWorker(q PrintQueue, tf TransportFactory, bus ...*events.Bus) *Worker {
	w := &Worker{
		queue:            q,
		transportFactory: tf,
		stopCh:           make(chan struct{}),
	}
	if len(bus) > 0 {
		w.bus = bus[0]
	}
	return w
}

func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.stopCh:
			return
		case <-ticker.C:
			w.processNext(ctx)
		}
	}
}

func (w *Worker) Stop() {
	close(w.stopCh)
}

func (w *Worker) processNext(ctx context.Context) {
	j := w.queue.Dequeue()
	if j == nil {
		return // Nothing to process
	}

	_ = w.queue.UpdateStatus(j.ID, job.StatusPrinting)
	w.publish(events.NewEvent(events.PrintStarted, map[string]any{
		"job_id":     j.ID,
		"printer_id": j.PrinterID,
	}))

	t, err := w.transportFactory(j.PrinterID)
	if err != nil {
		w.fail(j.ID, j.PrinterID, err)
		return
	}

	defer t.Close()

	if err := t.Open(ctx); err != nil {
		w.fail(j.ID, j.PrinterID, err)
		return
	}

	if err := t.Write(ctx, j.Payload); err != nil {
		w.fail(j.ID, j.PrinterID, err)
		return
	}

	_ = w.queue.UpdateStatus(j.ID, job.StatusCompleted)
	w.publish(events.NewEvent(events.PrintCompleted, map[string]any{
		"job_id":     j.ID,
		"printer_id": j.PrinterID,
	}))
}

func (w *Worker) fail(jobID, printerID string, err error) {
	_ = w.queue.UpdateStatus(jobID, job.StatusFailed)
	w.publish(events.NewEvent(events.PrintFailed, map[string]any{
		"job_id":     jobID,
		"printer_id": printerID,
		"error":      err.Error(),
	}))
}

func (w *Worker) publish(e events.Event) {
	if w.bus != nil {
		w.bus.Publish(e)
	}
}
