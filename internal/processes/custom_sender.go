package processes

import (
	"fmt"

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
	fmt.Println("Custom Send: ", msg.Value)
}
