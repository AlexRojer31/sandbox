package processes

import (
	"fmt"

	"github.com/AlexRojer31/sandbox/internal/dto"
	"github.com/AlexRojer31/sandbox/internal/recovery"
)

type observer struct {
	process
}

func NewObserver(name string) IProcess {
	reader := observer{}
	reader.process = newProcess(name+"Observer", nil, reader.handle)

	return &reader
}

func (p *observer) handle(msg dto.Data, errCh chan dto.Data) {
	defer recovery.Recover()
	fmt.Println(dto.ParceData[error](msg))
}
