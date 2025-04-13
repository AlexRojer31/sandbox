package processes

import (
	"context"
	"math/rand"
	"time"

	"github.com/AlexRojer31/sandbox/internal/dto"
	"github.com/AlexRojer31/sandbox/internal/recovery"
)

type emitter struct {
	*process
}

func NewEriter(to chan dto.Data) IProcess {
	emitter := emitter{process: newProcess("Emitter", to)}

	emitter.runf = emitter.run
	return &emitter
}

func (e *emitter) run(ctx context.Context, errCh chan dto.Data, from chan dto.Data, args ...any) {
	defer recovery.Recover()
	e.status <- 1
	e.logger.Info(e.name, " started.")
	for {
		select {
		case <-ctx.Done():
			close(e.to)
			close(e.status)
			return
		default:
			time.Sleep(time.Second)
			e.to <- dto.Data{
				Value: rand.Intn(100),
			}
		}
	}
}
