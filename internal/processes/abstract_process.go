package processes

import (
	"context"

	"github.com/AlexRojer31/sandbox/internal/config"
	"github.com/AlexRojer31/sandbox/internal/container"
	"github.com/AlexRojer31/sandbox/internal/dto"
	"github.com/AlexRojer31/sandbox/internal/recovery"
	"github.com/sirupsen/logrus"
)

type INamed interface {
	GetName() string
}

type IRun interface {
	Run(ctx context.Context, errCh chan<- dto.Data, from <-chan dto.Data) <-chan dto.Data
	Stop(errCh chan<- dto.Data)
}

type IProcess interface {
	INamed
	IRun
}

type Namef func() string
type Runf func(ctx context.Context, errCh chan<- dto.Data, from <-chan dto.Data)
type Stopf func(errCh chan<- dto.Data)
type Handlef func(msg dto.Data, errCh chan<- dto.Data)

type abstractProcess struct {
	name string
	to   chan dto.Data

	settings config.ProcessesSettings
	status   chan int
	logger   *logrus.Logger

	namef   Namef
	runf    Runf
	stopf   Stopf
	handlef Handlef
}

func newAbstractProcess(name string, args ...any) *abstractProcess {
	settings := container.GetInstance().Env.Config.ProcessesSettings
	process := abstractProcess{
		name:     name,
		settings: settings,
		to:       make(chan dto.Data, settings.Common.Size),
	}
	process.status = make(chan int, 1)
	process.logger = container.GetInstance().Logger

	process.runf = process.run
	process.namef = func() string {
		return process.name
	}
	process.stopf = func(errCh chan<- dto.Data) {
		process.logger.Warn(process.name, " stopping...")
		for v := range process.status {
			if v > 0 {
				continue
			}
		}
		process.logger.Warn(process.name, " stopped")
	}
	process.handlef = func(msg dto.Data, errCh chan<- dto.Data) {
		defer recovery.Recover()
		if process.to != nil {
			process.to <- msg
		}
	}

	for _, obj := range args {
		switch v := obj.(type) {
		case Handlef:
			process.handlef = v
		}
	}

	return &process
}

func (p *abstractProcess) GetName() string {
	return p.namef()
}

func (p *abstractProcess) Run(ctx context.Context, errCh chan<- dto.Data, from <-chan dto.Data) <-chan dto.Data {
	go p.runf(ctx, errCh, from)
	return p.to
}

func (p *abstractProcess) Stop(errCh chan<- dto.Data) {
	p.stopf(errCh)
}

func (p *abstractProcess) run(ctx context.Context, errCh chan<- dto.Data, from <-chan dto.Data) {
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
