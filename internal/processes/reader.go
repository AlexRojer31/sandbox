package processes

import (
	"context"
	"time"

	"github.com/AlexRojer31/sandbox/internal/dto"
	"github.com/AlexRojer31/sandbox/internal/recovery"
)

type Readerf func(ctx context.Context, errCh chan dto.Data, args ...any)
type Fetchf func() (dto.Data, error)
type Commitf func(msg dto.Data, errCh chan dto.Data)

type reader struct {
	process

	commitCh   chan dto.Data
	fetchf     Readerf
	commitf    Readerf
	fetchMsgf  Fetchf
	commitMsgf Commitf
}

func newReader(name string, to chan dto.Data, args ...any) reader {
	reader := reader{
		commitCh: make(chan dto.Data, 1000),
	}
	reader.process = newProcess(name+"Reader", to)

	for _, obj := range args {
		switch v := obj.(type) {
		case func() (dto.Data, error):
			reader.fetchMsgf = v
		case func(msg dto.Data, errCh chan dto.Data):
			reader.commitMsgf = v
		}
	}

	if reader.fetchMsgf == nil {
		reader.fetchMsgf = func() (dto.Data, error) {
			defer recovery.Recover()
			time.Sleep(time.Second)
			return dto.Data{
				Value: 51,
			}, nil
		}
	}

	if reader.commitMsgf == nil {
		reader.commitMsgf = func(msg dto.Data, errCh chan dto.Data) {}
	}

	reader.fetchf = reader.fetch
	reader.commitf = reader.commit
	reader.runf = reader.run
	return reader
}

func (r *reader) run(ctx context.Context, errCh chan dto.Data, from chan dto.Data, args ...any) {
	go r.fetchf(ctx, errCh, args...)
	go r.commitf(ctx, errCh, args...)
}

func (r *reader) fetch(ctx context.Context, errCh chan dto.Data, args ...any) {
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

func (r *reader) commit(ctx context.Context, errCh chan dto.Data, args ...any) {
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
