package processes

import (
	"fmt"

	"github.com/AlexRojer31/sandbox/internal/dto"
	"github.com/AlexRojer31/sandbox/internal/recovery"
)

type reader struct {
	process
}

func NewReader() IProcess {
	reader := reader{}
	reader.process = newProcess("Reader", nil, reader.handle)

	return &reader
}

func (p *reader) handle(msg dto.Data, errCh chan error) {
	recovery.Recover()
	fmt.Println(dto.ParceData[string](msg))
}
