package processes

type IBuilder interface {
	Build(ChainConfig) IChain
}

type Builder struct{}

func (b *Builder) Build(conf ChainConfig) IChain {
	return NewChain(conf.Name, b.makeProcesses(conf.Name, conf.Processes))
}

func (b *Builder) makeProcesses(name string, names []string) []IProcess {
	var proc []IProcess
	for _, n := range names {
		switch n {
		case "emitter":
			proc = append(proc, NewEmitter())
		case "sender":
			proc = append(proc, NewSender(name))
		}
	}

	return proc
}
