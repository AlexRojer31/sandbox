package processes

import (
	"context"
	"math/rand"
	"time"

	"github.com/AlexRojer31/sandbox/internal/dto"
	"github.com/AlexRojer31/sandbox/internal/recovery"
)

type abstractEmitter struct {
	*abstractProcess
}

func newAbstractEmitter(name string, args ...any) *abstractEmitter {
	emitter := abstractEmitter{abstractProcess: newAbstractProcess(name + "Emitter")}

	emitter.runf = emitter.run

	for _, obj := range args {
		switch v := obj.(type) {
		case Runf:
			emitter.runf = v
		}
	}
	return &emitter
}

func (e *abstractEmitter) run(ctx context.Context, errCh chan<- dto.Data, from <-chan dto.Data) {
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
