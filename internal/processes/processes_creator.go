package processes

type processesCreator struct {
	prefix string
}

func (pc *processesCreator) getAbstractEmitter() IProcess {
	return newEmitter(pc.prefix)
}

func (pc *processesCreator) getAbstractSender() IProcess {
	return newSender(pc.prefix)
}

func (pc *processesCreator) getCustomFilter() IProcess {
	return newCustomFilter(pc.prefix)
}
