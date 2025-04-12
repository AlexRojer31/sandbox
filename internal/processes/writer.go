package processes

import (
	"context"
	"math/rand"
	"time"

	"github.com/AlexRojer31/sandbox/internal/dto"
)

type writer struct {
	process
}

func NewWriter(to chan dto.Data) IProcess {
	writer := writer{process: newProcess("Writer", to, nil)}

	writer.process.runf = writer.run
	return &writer
}

func (w *writer) run(ctx context.Context, errCh chan error, from chan dto.Data, args ...any) {
	// defer recovery.Recover()
	w.process.status <- 1
	for {
		select {
		case <-ctx.Done():
			close(w.process.to)
			close(w.process.status)
			return
		default:
			time.Sleep(time.Second)
			w.process.to <- dto.Data{
				Value: rand.Intn(100),
			}
		}
	}
}
