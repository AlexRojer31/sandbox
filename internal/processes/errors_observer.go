package processes

import (
	"context"
	"fmt"

	"github.com/AlexRojer31/sandbox/internal/dto"
	"github.com/AlexRojer31/sandbox/internal/recovery"
)

type errorsObserver struct {
	process

	errors chan dto.Data
}

func NewErrorsObserver(errors chan dto.Data) IProcess {
	observer := errorsObserver{
		errors: errors,
	}
	observer.process = newProcess("ErrorObserver", nil)

	observer.process.runf = observer.run
	observer.process.stopf = observer.stop
	return &observer
}

func (o *errorsObserver) run(ctx context.Context, errCh chan dto.Data, from chan dto.Data, args ...any) {
	defer recovery.Recover()
	o.status <- 1
	o.logger.Info(o.name, " started.")
	for {
		select {
		case <-ctx.Done():
			for msg := range o.errors {
				o.handle(msg)
			}
			close(o.status)
			return
		case msg, ok := <-o.errors:
			if ok {
				o.handle(msg)
			}
		}
	}
}

func (o *errorsObserver) handle(msg dto.Data) {
	defer recovery.Recover()
	fmt.Println(dto.ParceData[error](msg))
}

func (o *errorsObserver) stop(errCh chan dto.Data, args ...any) {
	o.process.logger.Warn(o.process.name, " stopping...")
	close(o.errors)
	for v := range o.process.status {
		if v > 0 {
			continue
		}
	}
	o.process.logger.Warn(o.process.name, " stopped")
}
