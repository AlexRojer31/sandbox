package observer

import (
	"context"

	"github.com/AlexRojer31/sandbox/internal/dto"
	"github.com/AlexRojer31/sandbox/internal/recovery"
	"github.com/sirupsen/logrus"
)

type IObserve interface {
	GetName() string
	GetChannel() chan dto.Data
	Observe(ctx context.Context)
	Stop()
}

type Namef func() string
type Handlef func(msg dto.Data)
type Channelf func() chan dto.Data
type Observf func(ctx context.Context)
type Stopf func()

type observer[T any] struct {
	name      string
	observeCh chan dto.Data

	status chan int
	logger *logrus.Logger

	handlef  Handlef
	namef    Namef
	channelf Channelf
	observf  Observf
	stopf    Stopf
}

func newObserver[T any](name string, args ...any) *observer[T] {
	observeCh := make(chan dto.Data, 1000)
	observer := observer[T]{
		name:      name,
		observeCh: observeCh,
	}

	for _, obj := range args {
		switch v := obj.(type) {
		case func(msg dto.Data):
			observer.handlef = v
		}
	}

	if observer.handlef == nil {
		observer.handlef = func(msg dto.Data) {}
	}

	observer.namef = func() string {
		return observer.name
	}
	observer.channelf = func() chan dto.Data {
		return observer.observeCh
	}
	observer.observf = observer.observe
	observer.stopf = func() {
		observer.logger.Warn(observer.name, " stopping...")
		close(observer.observeCh)
		for v := range observer.status {
			if v > 0 {
				continue
			}
		}
		observer.logger.Warn(observer.name, " stopped")
	}

	return &observer
}

func (o *observer[T]) GetName() string {
	return o.namef()
}

func (o *observer[T]) GetChannel() chan dto.Data {
	return o.channelf()
}

func (o *observer[T]) Observe(ctx context.Context) {
	go o.observf(ctx)
}

func (o *observer[T]) Stop() {
	o.stopf()
}

func (o *observer[T]) observe(ctx context.Context) {
	defer recovery.Recover()
	o.status <- 1
	o.logger.Info(o.name, " started.")
	for {
		select {
		case <-ctx.Done():
			for msg := range o.observeCh {
				o.handlef(msg)
			}
			close(o.status)
			return
		case msg, ok := <-o.observeCh:
			if ok {
				o.handlef(msg)
			}
		}
	}
}
