package sandbox

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/AlexRojer31/sandbox/internal/app"
	"github.com/AlexRojer31/sandbox/internal/dto"
	"github.com/AlexRojer31/sandbox/internal/processes"
	"github.com/AlexRojer31/sandbox/internal/recovery"
	"github.com/golangci/golangci-lint/pkg/exitcodes"
)

func Run(args []string) int {
	defer recovery.Recover()
	sandbox := app.GetInstance(args)

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
		sandbox.Logger.Info("Catch signal " + sig.String())
		ctxCancel()
		writer.Stop(errors)
		reader.Stop(errors)

		return exitcodes.Success
	}
	ctxCancel()

	return exitcodes.Failure
}
