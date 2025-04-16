package processes

import (
	"context"
	"time"

	"github.com/AlexRojer31/sandbox/internal/dto"
	"github.com/AlexRojer31/sandbox/internal/recovery"
)

type Readerf func(ctx context.Context, errCh chan<- dto.Data)
type Fetchf func() (dto.Data, error)
type Commitf func(msg dto.Data, errCh chan<- dto.Data)

type reader struct {
	*abstractProcess

	commitCh   chan dto.Data
	fetchf     Readerf
	commitf    Readerf
	fetchMsgf  Fetchf
	commitMsgf Commitf
}

func newReader(name string, args ...any) *reader {
	reader := reader{
		commitCh: make(chan dto.Data, 10000),
	}
	reader.abstractProcess = newProcess(name + "Reader")

	reader.runf = reader.run
	reader.fetchf = reader.fetch
	reader.fetchMsgf = func() (dto.Data, error) {
		defer recovery.Recover()
		time.Sleep(time.Second)
		return dto.Data{
			Value: 51,
		}, nil
	}
	reader.commitf = reader.commit
	reader.commitMsgf = func(msg dto.Data, errCh chan<- dto.Data) {}

	for _, obj := range args {
		switch v := obj.(type) {
		case Fetchf:
			reader.fetchMsgf = v
		case Commitf:
			reader.commitMsgf = v
		}
	}
	return &reader
}

func (r *reader) run(ctx context.Context, errCh chan<- dto.Data, from <-chan dto.Data) {
	go r.fetchf(ctx, errCh)
	go r.commitf(ctx, errCh)
}

func (r *reader) fetch(ctx context.Context, errCh chan<- dto.Data) {
	defer recovery.Recover()
	r.status <- 1
	r.logger.Info(r.name, " fetch started.")
	for {
		select {
		case <-ctx.Done():
			close(r.to)
			close(r.commitCh)
			return
		default:
			msg, err := r.fetchMsgf()
			if err != nil {
				errCh <- dto.Data{
					Value: err,
				}
			} else {
				r.to <- msg
				r.commitCh <- msg
			}
		}
	}
}

func (r *reader) commit(ctx context.Context, errCh chan<- dto.Data) {
	defer recovery.Recover()
	r.logger.Info(r.name, " commit started.")
	for {
		select {
		case <-ctx.Done():
			for msg := range r.commitCh {
				r.commitMsgf(msg, errCh)
			}
			close(r.status)
			return
		case msg, ok := <-r.commitCh:
			if ok {
				r.commitMsgf(msg, errCh)
			}
		}
	}
}
