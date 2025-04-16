package processes

import (
	"github.com/AlexRojer31/sandbox/internal/dto"
)

type customFilter struct {
	*abstractFilter
}

func newCustomFilter(name string) IProcess {
	filter := customFilter{}
	filter.abstractFilter = newAbstractFilter(name+"Custom", (Filterf)(filter.filter))

	return &filter
}

func (f *customFilter) filter(msg dto.Data) bool {
	v, ok := msg.Value.(int)
	if ok {
		return v > 10
	}
	return false
}
