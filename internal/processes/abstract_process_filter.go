package processes

import (
	"github.com/AlexRojer31/sandbox/internal/dto"
	"github.com/AlexRojer31/sandbox/internal/recovery"
)

type Filterf func(msg dto.Data) bool

type abstractFilter struct {
	*abstractProcess

	filterf Filterf
}

func newAbstractFilter(name string, args ...any) *abstractFilter {
	filter := abstractFilter{}
	filter.abstractProcess = newAbstractProcess("Filter", (Handlef)(filter.handle))

	filter.filterf = filter.filter

	for _, obj := range args {
		switch v := obj.(type) {
		case Filterf:
			filter.filterf = v
		}
	}

	return &filter
}

func (f *abstractFilter) handle(msg dto.Data, errCh chan<- dto.Data) {
	defer recovery.Recover()
	if f.filterf(msg) {
		f.to <- msg
	}
}

func (f *abstractFilter) filter(msg dto.Data) bool {
	return true
}
