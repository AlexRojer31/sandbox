package chain

import (
	"context"

	"github.com/AlexRojer31/sandbox/internal/dto"
	"github.com/AlexRojer31/sandbox/internal/processes"
)

type ChainConfig struct {
	Name      string
	Processes []string
}

type IChain interface {
	Run(ctx context.Context, errCh chan<- dto.Data)
	Stop(errCh chan<- dto.Data)
}

type Chain struct {
	name      string
	handlers  []IHandler
	processes []processes.IProcess
}

func NewChain(name string, processes []processes.IProcess) IChain {
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
