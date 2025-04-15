package sandbox

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/AlexRojer31/sandbox/internal/container"
	"github.com/AlexRojer31/sandbox/internal/dto"
	"github.com/AlexRojer31/sandbox/internal/observer"
	"github.com/AlexRojer31/sandbox/internal/processes"
	"github.com/AlexRojer31/sandbox/internal/recovery"
	"github.com/golangci/golangci-lint/pkg/exitcodes"
)

type Sandbox struct {
	container *container.Container
}

func Run(args []string) int {
	defer recovery.Recover()
	sandbox := Sandbox{
		container: container.GetInstance(args),
	}

	ctx, ctxCancel := context.WithCancel(context.Background())
	errorObserver := observer.NewErrorObserver()
	errCh := errorObserver.GetChannel()
	errorObserver.Observe(ctx)

	chain := NewChain(
		"MyChain",
		[]processes.IProcess{
			processes.NewCustomReader(),
			processes.NewFilter(func(msg dto.Data) bool {
				if v, ok := msg.Value.(int); ok {
					return v > 50
				}
				return false
			}),
			processes.NewSender("Super"),
		},
	)
	chain.Run(ctx, errCh)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	if sig, ok := <-interrupt; ok {
		sandbox.container.Logger.Info("Catch signal ", sig.String())
		ctxCancel()

		chain.Stop(errCh)

		errorObserver.Stop()
		return exitcodes.Success
	}
	ctxCancel()

	return exitcodes.Failure
}

type Chain struct {
	name      string
	handlers  []IHandler
	processes []processes.IProcess
}

func NewChain(name string, processes []processes.IProcess) *Chain {
	len := len(processes)
	chain := Chain{
		name:     name,
		handlers: make([]IHandler, len),
	}
	chain.processes = append(chain.processes, processes...)

	lastElem := len - 1
	for i := 0; i < len; i++ {
		elemNum := len - 1 - i
		if elemNum == lastElem {
			chain.handlers[elemNum] = NewHandler(processes[elemNum], nil)
		} else if elemNum < lastElem {
			chain.handlers[elemNum] = NewHandler(processes[elemNum], chain.handlers[elemNum+1])
		}
	}

	return &chain
}

func (c *Chain) Run(ctx context.Context, errCh chan<- dto.Data) {
	c.handlers[0].Next(ctx, errCh, nil)
}

func (c *Chain) Stop(errCh chan<- dto.Data) {
	for _, p := range c.processes {
		p.Stop(errCh)
	}
}

type Handler struct {
	process processes.IProcess
	next    IHandler
}

func NewHandler(process processes.IProcess, next IHandler) IHandler {
	return &Handler{
		process: process,
		next:    next,
	}
}

type IHandler interface {
	Next(ctx context.Context, errCh chan<- dto.Data, from <-chan dto.Data)
}

func (h *Handler) Next(ctx context.Context, errCh chan<- dto.Data, from <-chan dto.Data) {
	fromCh := h.process.Run(ctx, errCh, from)
	if h.next != nil {
		h.next.Next(ctx, errCh, fromCh)
	}
}
