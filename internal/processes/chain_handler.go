package processes

import (
	"context"

	"github.com/AlexRojer31/sandbox/internal/dto"
)

type IHandler interface {
	Next(ctx context.Context, errCh chan<- dto.Data, from <-chan dto.Data)
}

type Handler struct {
	process IProcess
	next    IHandler
}

func NewHandler(process IProcess, next IHandler) IHandler {
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
