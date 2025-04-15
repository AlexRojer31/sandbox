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

	reader := processes.NewCustomReader()
	filter := processes.NewFilter(func(msg dto.Data) bool {
		return dto.ParceData[int](msg) > 50
	})
	sender := processes.NewSender("Super")

	reader2filter := reader.Run(ctx, errCh, nil)
	filter2sender := filter.Run(ctx, errCh, reader2filter)
	sender.Run(ctx, errCh, filter2sender)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	if sig, ok := <-interrupt; ok {
		sandbox.container.Logger.Info("Catch signal ", sig.String())
		ctxCancel()

		sender.Stop(errCh)
		filter.Stop(errCh)
		reader.Stop(errCh)

		errorObserver.Stop()
		return exitcodes.Success
	}
	ctxCancel()

	return exitcodes.Failure
}
