package sandbox

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/AlexRojer31/sandbox/internal/container"
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

	builder := processes.Builder{}
	chain1 := builder.Build(processes.ChainConfig{
		Name:      "MyTestChain",
		Processes: []string{"emitter", "filter", "sender"},
	})
	chain1.Run(ctx, errCh)

	chain2 := builder.Build(processes.ChainConfig{
		Name:      "MyAnotherChain",
		Processes: []string{"reader", "filter", "sender"},
	})
	chain2.Run(ctx, errCh)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	if sig, ok := <-interrupt; ok {
		sandbox.container.Logger.Info("Catch signal ", sig.String())
		ctxCancel()

		chain1.Stop(errCh)
		chain2.Stop(errCh)

		errorObserver.Stop()
		return exitcodes.Success
	}
	ctxCancel()

	return exitcodes.Failure
}
