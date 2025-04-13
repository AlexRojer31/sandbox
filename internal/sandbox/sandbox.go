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

	ctx, ctxCancel := context.WithCancel(context.Background())
	errorObserver, errCh := processes.NewObserver[error]("errors")
	errorObserver.Run(ctx, nil, nil)

	emitter2filter := make(chan dto.Data, 1000)
	filter2sender := make(chan dto.Data, 1000)
	emitter := processes.NewEriter(emitter2filter)
	filter := processes.NewFilter(filter2sender, func(msg dto.Data) bool {
		return dto.ParceData[int](msg) > 50
	})
	sender := processes.NewSender("blaBlaBla")

	emitter.Run(ctx, errCh, nil)
	filter.Run(ctx, errCh, emitter2filter)
	sender.Run(ctx, errCh, filter2sender)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	if sig, ok := <-interrupt; ok {
		sandbox.container.Logger.Info("Catch signal ", sig.String())
		ctxCancel()
		emitter.Stop(errCh)
		filter.Stop(errCh)
		sender.Stop(errCh)

		errorObserver.Stop(nil)
		return exitcodes.Success
	}
	ctxCancel()

	return exitcodes.Failure
}
