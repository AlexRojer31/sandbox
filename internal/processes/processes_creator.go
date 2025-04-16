package processes

type processesCreator struct {
	prefix string
}

func (pc *processesCreator) getCustomEmitter() IProcess {
	return newCustomEmitter(pc.prefix)
}

func (pc *processesCreator) getCustomSender() IProcess {
	return newCustomSender(pc.prefix)
}

func (pc *processesCreator) getCustomFilter() IProcess {
	return newCustomFilter(pc.prefix)
}

func (pc *processesCreator) getCustomReader() IProcess {
	return newCustomReader(pc.prefix)
}
