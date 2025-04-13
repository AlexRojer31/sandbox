package processes

import (
	"context"
	"fmt"

	"github.com/AlexRojer31/sandbox/internal/dto"
	"github.com/AlexRojer31/sandbox/internal/recovery"
)

type observer[T any] struct {
	process

	observeCh chan dto.Data
}

func NewObserver[T any](name string) (IProcess, chan dto.Data) {
	observeCh := make(chan dto.Data, 1000)
	observer := observer[T]{
		observeCh: observeCh,
	}
	observer.process = newProcess(name+"Observer", nil)

	observer.runf = observer.run
	observer.stopf = observer.stop
	return &observer, observeCh
}

func (o *observer[T]) run(ctx context.Context, errCh chan dto.Data, from chan dto.Data, args ...any) {
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

func (o *observer[T]) handle(msg dto.Data) {
	defer recovery.Recover()
	fmt.Println(dto.ParceData[T](msg))
}

func (o *observer[T]) stop(errCh chan dto.Data, args ...any) {
	o.logger.Warn(o.name, " stopping...")
	close(o.observeCh)
	for v := range o.status {
		if v > 0 {
			continue
		}
	}
	o.logger.Warn(o.name, " stopped")
}
