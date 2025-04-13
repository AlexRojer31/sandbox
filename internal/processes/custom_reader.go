package processes

import (
	"math/rand"
	"time"

	"github.com/AlexRojer31/sandbox/internal/dto"
	"github.com/AlexRojer31/sandbox/internal/recovery"
)

type customReader struct {
	reader
}

func NewCustomReader(to chan dto.Data, args ...any) IProcess {
	custom := customReader{}
	custom.reader = newReader("Custom", to, custom.fetchMsg, custom.commitMsg)

	return &custom
}

func (r *customReader) fetchMsg() (dto.Data, error) {
	defer recovery.Recover()
	time.Sleep(time.Second)
	return dto.Data{
		Value: rand.Intn(100),
	}, nil
}

func (r *customReader) commitMsg(msg dto.Data, errCh chan dto.Data) {
	r.logger.Info("I COMMIT!")
}
