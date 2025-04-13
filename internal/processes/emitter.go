package processes

import (
	"context"
	"math/rand"
	"time"

	"github.com/AlexRojer31/sandbox/internal/dto"
	"github.com/AlexRojer31/sandbox/internal/recovery"
)

type emitter struct {
	process
}

func NewWriter(to chan dto.Data) IProcess {
	emitter := emitter{process: newProcess("Emitter", to)}

	emitter.process.runf = emitter.run
	return &emitter
}

func (e *emitter) run(ctx context.Context, errCh chan error, from chan dto.Data, args ...any) {
	defer recovery.Recover()
	e.process.status <- 1
	e.logger.Info(e.name, " started.")
	for {
		select {
		case <-ctx.Done():
			close(e.process.to)
			close(e.process.status)
			return
		default:
			time.Sleep(time.Second)
			e.process.to <- dto.Data{
				Value: rand.Intn(100),
			}
		}
	}
}
