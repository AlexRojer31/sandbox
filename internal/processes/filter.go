package processes

import (
	"github.com/AlexRojer31/sandbox/internal/dto"
	"github.com/AlexRojer31/sandbox/internal/recovery"
)

type Filterf func(msg dto.Data) bool

type filter struct {
	process

	filterf Filterf
}

func NewFilter(to chan dto.Data, filterf Filterf) IProcess {
	filter := filter{}
	filter.process = newProcess("Filter", to, filter.handle)

	if filterf == nil {
		filter.filterf = filter.filter
	} else {
		filter.filterf = filterf
	}
	return &filter
}

func (f *filter) handle(msg dto.Data, errCh chan dto.Data) {
	defer recovery.Recover()
	if f.filterf(msg) {
		f.to <- msg
	}
}

func (f *filter) filter(msg dto.Data) bool {
	return true
}
