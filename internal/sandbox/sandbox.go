package sandbox

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/AlexRojer31/sandbox/internal/chain"
	"github.com/AlexRojer31/sandbox/internal/chain_builder"
	"github.com/AlexRojer31/sandbox/internal/container"
	"github.com/AlexRojer31/sandbox/internal/observer"
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

	// chain := chain.NewChain(
	// 	"MyChain",
	// 	[]processes.IProcess{
	// 		processes.NewCustomReader(),
	// 		processes.NewFilter(func(msg dto.Data) bool {
	// 			if v, ok := msg.Value.(int); ok {
	// 				return v > 50
	// 			}
	// 			return false
	// 		}),
	// 		processes.NewSender("Super"),
	// 	},
	// )
	// chain.Run(ctx, errCh)

	builder := chain_builder.Builder{}
	chain := builder.Build(chain.ChainConfig{
		Name:      "MyTestChain",
		Processes: []string{"emitter", "sender"},
	})
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
