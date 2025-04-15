package chain_builder

import (
	"github.com/AlexRojer31/sandbox/internal/chain"
	"github.com/AlexRojer31/sandbox/internal/processes"
)

type IBuilder interface {
	Build(chain.ChainConfig) chain.IChain
}

type Builder struct{}

func (b *Builder) Build(conf chain.ChainConfig) chain.IChain {
	return chain.NewChain(conf.Name, b.makeProcesses(conf.Name, conf.Processes))
}

func (b *Builder) makeProcesses(name string, names []string) []processes.IProcess {
	var proc []processes.IProcess
	for _, n := range names {
		switch n {
		case "emitter":
			proc = append(proc, processes.NewEmitter())
		case "sender":
			proc = append(proc, processes.NewSender(name))
		}
	}

	return proc
}
