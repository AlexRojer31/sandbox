package processes

import (
	"github.com/AlexRojer31/sandbox/internal/dto"
	"github.com/AlexRojer31/sandbox/internal/recovery"
)

type customSender struct {
	*abstractSender
}

func newCustomSender(name string) IProcess {
	sender := customSender{}
	sender.abstractSender = newAbstractSender(name+"Custom", (Handlef)(sender.handle))

	return &sender
}

func (s *customSender) handle(msg dto.Data, errCh chan<- dto.Data) {
	defer recovery.Recover()
	s.logger.Info(s.name, " Custom Send: ", msg.Value)
}
