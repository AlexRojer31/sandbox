package sandbox

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/AlexRojer31/sandbox/internal/container"
	"github.com/AlexRojer31/sandbox/internal/dto"
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

	errors := make(chan error, 10000)
	ctx, ctxCancel := context.WithCancel(context.Background())

	data := make(chan dto.Data, 1)
	writer := processes.NewWriter(data)
	reader := processes.NewReader()

	writer.Run(ctx, errors, nil)
	reader.Run(ctx, errors, data)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	if sig, ok := <-interrupt; ok {
		sandbox.container.Logger.Info("Catch signal ", sig.String())
		ctxCancel()
		writer.Stop(errors)
		reader.Stop(errors)

		return exitcodes.Success
	}
	ctxCancel()

	return exitcodes.Failure
}
