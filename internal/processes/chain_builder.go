package processes

import "github.com/AlexRojer31/sandbox/internal/config"

type IBuilder interface {
	Build(conf config.ChainConfig) IChain
}

type Builder struct{}

func (b *Builder) Build(conf config.ChainConfig) IChain {
	return NewChain(conf.Name, b.makeProcesses(conf.Name, conf.Processes))
}

func (b *Builder) makeProcesses(name string, names []string) []IProcess {
	var proc []IProcess
	pc := processesCreator{prefix: name}
	for _, n := range names {
		switch n {
		case "emitter":
			proc = append(proc, pc.getCustomEmitter())
		case "filter":
			proc = append(proc, pc.getCustomFilter())
		case "sender":
			proc = append(proc, pc.getCustomSender())
		case "reader":
			proc = append(proc, pc.getCustomReader())
		}
	}

	return proc
}
