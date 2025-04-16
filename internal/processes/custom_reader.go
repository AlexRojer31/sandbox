package processes

import (
	"math/rand"
	"time"

	"github.com/AlexRojer31/sandbox/internal/dto"
	"github.com/AlexRojer31/sandbox/internal/recovery"
)

type customReader struct {
	*reader
}

func newCustomReader(name string) IProcess {
	custom := customReader{}
	custom.reader = newReader(name+"Custom", (Fetchf)(custom.fetchMsg), (Commitf)(custom.commitMsg))

	return &custom
}

func (r *customReader) fetchMsg() (dto.Data, error) {
	defer recovery.Recover()
	time.Sleep(time.Second)
	r.logger.Info("Fetching")
	return dto.Data{
		Value: rand.Intn(100),
	}, nil
}

func (r *customReader) commitMsg(msg dto.Data, errCh chan<- dto.Data) {
	r.logger.Info("Commiting")
}
