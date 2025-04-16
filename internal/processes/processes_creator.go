package processes

type IProcessesCreator interface {
	GetCustomEmitter() IProcess
	GetCustomSender() IProcess
	GetCustomFilter() IProcess
	GetCustomReader() IProcess
}

type processesCreator struct {
	prefix string
}

func NewProcessCreator(prefix string) IProcessesCreator {
	return &processesCreator{prefix: prefix}
}

func (pc *processesCreator) GetCustomEmitter() IProcess {
	return newCustomEmitter(pc.prefix)
}

func (pc *processesCreator) GetCustomSender() IProcess {
	return newCustomSender(pc.prefix)
}

func (pc *processesCreator) GetCustomFilter() IProcess {
	return newCustomFilter(pc.prefix)
}

func (pc *processesCreator) GetCustomReader() IProcess {
	return newCustomReader(pc.prefix)
}
