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
	chains    []processes.IChain
}

func Run(args []string) int {
	defer recovery.Recover()
	sandbox := Sandbox{
		container: container.GetInstance(args),
	}

	ctx, ctxCancel := context.WithCancel(context.Background())
	errorObserver := observer.NewErrorObserver()
	errorObserver.Observe(ctx)
	errCh := errorObserver.GetChannel()

	go sandbox.Run(ctx, errCh)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	if sig, ok := <-interrupt; ok {
		sandbox.container.Logger.Info("Catch signal ", sig.String())
		ctxCancel()

		sandbox.Stop(errCh)

		errorObserver.Stop()
		return exitcodes.Success
	}
	ctxCancel()

	return exitcodes.Failure
}

func (sandbox *Sandbox) Run(ctx context.Context, errCh chan<- dto.Data) {
	builder := processes.Builder{}
	for _, c := range sandbox.container.Env.Config.Chains {
		sandbox.chains = append(sandbox.chains, builder.Build(c))
	}
	for _, chain := range sandbox.chains {
		chain.Run(ctx, errCh)
	}
}

func (sandbox *Sandbox) Stop(errCh chan<- dto.Data) {
	for _, chain := range sandbox.chains {
		chain.Stop(errCh)
	}
}
