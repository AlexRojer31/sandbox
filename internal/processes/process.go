package processes

import (
	"context"

	"github.com/AlexRojer31/sandbox/internal/container"
	"github.com/AlexRojer31/sandbox/internal/dto"
	"github.com/AlexRojer31/sandbox/internal/recovery"
	"github.com/sirupsen/logrus"
)

type Status int

type INamed interface {
	GetName() string
}

type IRun interface {
	Run(ctx context.Context, errCh chan error, from chan dto.Data, args ...any)
	Stop(errCh chan error, args ...any)
}

type IProcess interface {
	INamed
	IRun
}

type handle func(msg dto.Data, errCh chan error)

type process struct {
	name string
	to   chan dto.Data

	status chan int
	logger *logrus.Logger

	namef   func() string
	runf    func(ctx context.Context, errCh chan error, from chan dto.Data, args ...any)
	stopf   func(errCh chan error, args ...any)
	handlef handle
}

func newProcess(name string, to chan dto.Data, handle handle) process {
	process := process{
		name: name,
		to:   to,
	}
	if to == nil {
		process.to = make(chan dto.Data, 1)
	}
	process.status = make(chan int, 1)
	process.logger = container.GetInstance().Logger

	process.runf = process.run
	process.namef = func() string {
		return process.namef()
	}
	process.stopf = func(errCh chan error, args ...any) {
		process.logger.Warn(process.name, " stopping...")
		for v := range process.status {
			if v > 0 {
				continue
			}
		}
		process.logger.Warn(process.name, " stopped")
	}

	if handle == nil {
		process.handlef = func(msg dto.Data, errCh chan error) {
			// recovery.Recover()
			process.to <- msg
		}
	} else {
		process.handlef = handle
	}
	return process
}

func (p *process) GetName() string {
	return p.namef()
}

func (p *process) Run(ctx context.Context, errCh chan error, from chan dto.Data, args ...any) {
	go p.runf(ctx, errCh, from, args...)
}

func (p *process) Stop(errCh chan error, args ...any) {
	p.stopf(errCh, args...)
}

func (p *process) run(ctx context.Context, errCh chan error, from chan dto.Data, args ...any) {
	defer recovery.Recover()
	p.status <- 1
	for {
		select {
		case <-ctx.Done():
			for msg := range from {
				p.handlef(msg, errCh)
			}
			close(p.to)
			close(p.status)
			return
		case msg, ok := <-from:
			if ok {
				p.handlef(msg, errCh)
			}
		}
	}
}
