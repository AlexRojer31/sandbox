package processes

import (
	"context"

	"github.com/AlexRojer31/sandbox/internal/common"
	"github.com/AlexRojer31/sandbox/internal/container"
	"github.com/AlexRojer31/sandbox/internal/dto"
	"github.com/AlexRojer31/sandbox/internal/recovery"
	"github.com/sirupsen/logrus"
)

type IRun interface {
	Run(ctx context.Context, errCh chan dto.Data, from chan dto.Data, args ...any)
	Stop(errCh chan dto.Data, args ...any)
}

type IProcess interface {
	common.INamed
	IRun
}

type Namef func() string
type Runf func(ctx context.Context, errCh chan dto.Data, from chan dto.Data, args ...any)
type Stopf func(errCh chan dto.Data, args ...any)
type Handlef func(msg dto.Data, errCh chan dto.Data)

type process struct {
	name string
	to   chan dto.Data

	status chan int
	logger *logrus.Logger

	namef   Namef
	runf    Runf
	stopf   Stopf
	handlef Handlef
}

func newProcess(name string, to chan dto.Data, args ...any) *process {
	process := process{
		name: name,
		to:   to,
	}
	process.status = make(chan int, 1)
	process.logger = container.GetInstance().Logger

	for _, obj := range args {
		switch v := obj.(type) {
		case func(msg dto.Data, errCh chan dto.Data):
			process.handlef = v
		}
	}

	if process.handlef == nil {
		process.handlef = func(msg dto.Data, errCh chan dto.Data) {
			defer recovery.Recover()
			if process.to != nil {
				process.to <- msg
			}
		}
	}

	process.runf = process.run
	process.namef = func() string {
		return process.name
	}
	process.stopf = func(errCh chan dto.Data, args ...any) {
		process.logger.Warn(process.name, " stopping...")
		for v := range process.status {
			if v > 0 {
				continue
			}
		}
		process.logger.Warn(process.name, " stopped")
	}
	return &process
}

func (p *process) GetName() string {
	return p.namef()
}

func (p *process) Run(ctx context.Context, errCh chan dto.Data, from chan dto.Data, args ...any) {
	go p.runf(ctx, errCh, from, args...)
}

func (p *process) Stop(errCh chan dto.Data, args ...any) {
	p.stopf(errCh, args...)
}

func (p *process) run(ctx context.Context, errCh chan dto.Data, from chan dto.Data, args ...any) {
	defer recovery.Recover()
	p.status <- 1
	p.logger.Info(p.name, " started.")
	for {
		select {
		case <-ctx.Done():
			for msg := range from {
				p.handlef(msg, errCh)
			}
			if p.to != nil {
				close(p.to)
			}
			close(p.status)
			return
		case msg, ok := <-from:
			if ok {
				p.handlef(msg, errCh)
			}
		}
	}
}
