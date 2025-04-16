package processes

import (
	"fmt"

	"github.com/AlexRojer31/sandbox/internal/dto"
	"github.com/AlexRojer31/sandbox/internal/recovery"
)

type abstractSender struct {
	*abstractProcess
}

func newAbstractSender(name string, args ...any) *abstractSender {
	sender := abstractSender{}
	sender.abstractProcess = newAbstractProcess(name+"Sender", (Handlef)(sender.handle))

	for _, obj := range args {
		switch v := obj.(type) {
		case Handlef:
			sender.handlef = v
		}
	}

	return &sender
}

func (s *abstractSender) handle(msg dto.Data, errCh chan<- dto.Data) {
	defer recovery.Recover()
	fmt.Println("Send to another service: ", msg.Value)
}
