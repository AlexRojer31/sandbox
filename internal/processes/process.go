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

type Handlef func(msg dto.Data, errCh chan error)

type process struct {
	name string
	to   chan dto.Data

	status chan int
	logger *logrus.Logger

	namef   func() string
	runf    func(ctx context.Context, errCh chan error, from chan dto.Data, args ...any)
	stopf   func(errCh chan error, args ...any)
	handlef Handlef
}

func newProcess(name string, to chan dto.Data, args ...any) process {
	process := process{
		name: name,
		to:   to,
	}
	process.status = make(chan int, 1)
	process.logger = container.GetInstance().Logger

	for _, obj := range args {
		switch v := obj.(type) {
		case func(msg dto.Data, errCh chan error):
			process.handlef = v
		}
	}

	if process.handlef == nil {
		process.handlef = func(msg dto.Data, errCh chan error) {
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
	process.stopf = func(errCh chan error, args ...any) {
		process.logger.Warn(process.name, " stopping...")
		for v := range process.status {
			if v > 0 {
				continue
			}
		}
		process.logger.Warn(process.name, " stopped")
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
