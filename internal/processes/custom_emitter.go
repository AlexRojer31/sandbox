package processes

import (
	"context"
	"math/rand"
	"time"

	"github.com/AlexRojer31/sandbox/internal/dto"
	"github.com/AlexRojer31/sandbox/internal/recovery"
)

type customEmitter struct {
	*abstractEmitter
}

func newCustomEmitter(name string) IProcess {
	emitter := customEmitter{abstractEmitter: newAbstractEmitter(name + "Custom")}

	emitter.runf = emitter.run
	return &emitter
}

func (e *customEmitter) run(ctx context.Context, errCh chan<- dto.Data, from <-chan dto.Data) {
	defer recovery.Recover()
	e.status <- 1
	e.logger.Info(e.name, " started.")
	for {
		select {
		case <-ctx.Done():
			close(e.to)
			e.status <- -1
			return
		default:
			time.Sleep(time.Second)
			e.to <- dto.Data{
				Value: rand.Intn(100),
			}
		}
	}
}
