package processes

import (
	"github.com/AlexRojer31/sandbox/internal/dto"
)

type customFilter struct {
	*abstractFilter

	filterValue int
}

func newCustomFilter(name string) IProcess {
	filter := customFilter{}
	filter.abstractFilter = newAbstractFilter(name+"Custom", (Filterf)(filter.filter))
	filter.to = make(chan dto.Data, filter.settings.CustomFilterSetting.Size)
	filter.filterValue = filter.settings.CustomFilterSetting.MinValue

	return &filter
}

func (f *customFilter) filter(msg dto.Data) bool {
	v, ok := msg.Value.(int)
	if ok {
		return v > f.filterValue
	}
	return false
}
