package chain

import (
	"context"

	"github.com/AlexRojer31/sandbox/internal/dto"
	"github.com/AlexRojer31/sandbox/internal/processes"
)

type IHandler interface {
	Next(ctx context.Context, errCh chan<- dto.Data, from <-chan dto.Data)
}

type Handler struct {
	process processes.IProcess
	next    IHandler
}

func NewHandler(process processes.IProcess, next IHandler) IHandler {
	return &Handler{
		process: process,
		next:    next,
	}
}

func (h *Handler) Next(ctx context.Context, errCh chan<- dto.Data, from <-chan dto.Data) {
	fromCh := h.process.Run(ctx, errCh, from)
	if h.next != nil {
		h.next.Next(ctx, errCh, fromCh)
	}
}
