package processes

import (
	"context"
	"fmt"

	"github.com/AlexRojer31/sandbox/internal/dto"
	"github.com/AlexRojer31/sandbox/internal/recovery"
)

type observer struct {
	process

	observeCh chan dto.Data
}

func NewObserver(name string) (IProcess, chan dto.Data) {
	observeCh := make(chan dto.Data, 1000)
	observer := observer{
		observeCh: observeCh,
	}
	observer.process = newProcess(name+"Observer", nil)

	observer.process.runf = observer.run
	observer.process.stopf = observer.stop
	return &observer, observeCh
}

func (o *observer) run(ctx context.Context, errCh chan dto.Data, from chan dto.Data, args ...any) {
	defer recovery.Recover()
	o.status <- 1
	o.logger.Info(o.name, " started.")
	for {
		select {
		case <-ctx.Done():
			for msg := range o.observeCh {
				o.handle(msg)
			}
			close(o.status)
			return
		case msg, ok := <-o.observeCh:
			if ok {
				o.handle(msg)
			}
		}
	}
}

func (o *observer) handle(msg dto.Data) {
	defer recovery.Recover()
	fmt.Println(dto.ParceData[error](msg))
}

func (o *observer) stop(errCh chan dto.Data, args ...any) {
	o.process.logger.Warn(o.process.name, " stopping...")
	close(o.observeCh)
	for v := range o.process.status {
		if v > 0 {
			continue
		}
	}
	o.process.logger.Warn(o.process.name, " stopped")
}
