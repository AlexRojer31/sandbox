package observer

import (
	"github.com/AlexRojer31/sandbox/internal/dto"
	"github.com/AlexRojer31/sandbox/internal/recovery"
)

type errorObserver struct {
	*observer[error]
}

func NewErrorObserver() IObserve {
	errorObserver := errorObserver{}
	errorObserver.observer = (*observer[error])(newObserver[error]("Error", errorObserver.handle))

	return &errorObserver
}

func (o *errorObserver) handle(msg dto.Data) {
	recovery.Recover()
	o.logger.Error(msg.Value)
}
