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

	emitter2filter := make(chan dto.Data, 1000)
	filter2reader := make(chan dto.Data, 1000)
	emitter := processes.NewEriter(emitter2filter)
	filter := processes.NewFilter(filter2reader, func(msg dto.Data) bool {
		return dto.ParceData[int](msg) > 50
	})
	reader := processes.NewReader()

	emitter.Run(ctx, errors, nil)
	filter.Run(ctx, errors, emitter2filter)
	reader.Run(ctx, errors, filter2reader)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	if sig, ok := <-interrupt; ok {
		sandbox.container.Logger.Info("Catch signal ", sig.String())
		ctxCancel()
		emitter.Stop(errors)
		filter.Stop(errors)
		reader.Stop(errors)

		return exitcodes.Success
	}
	ctxCancel()

	return exitcodes.Failure
}
