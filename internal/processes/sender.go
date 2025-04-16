package processes

import (
	"fmt"

	"github.com/AlexRojer31/sandbox/internal/dto"
	"github.com/AlexRojer31/sandbox/internal/recovery"
)

type sender struct {
	*abstractProcess
}

func NewSender(name string) IProcess {
	sender := sender{}
	sender.abstractProcess = newProcess(name+"Sender", (Handlef)(sender.handle))

	return &sender
}

func (s *sender) handle(msg dto.Data, errCh chan<- dto.Data) {
	defer recovery.Recover()
	fmt.Println("Send to another service: ", msg.Value)
}
