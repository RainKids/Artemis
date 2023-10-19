package recovery

import (
	"context"
	"go.uber.org/zap"
	"log"
)

// Handler is a function that recovers from the panic `p` by returning an `error`.
// The context can be used to extract request scoped metadata and context values.
type Handler interface {
	Handle(ctx context.Context, p interface{}) error
}

// NewZapRecoveryHandler todo
func NewZapRecoveryHandler() *ZapRecoveryHandler {
	return &ZapRecoveryHandler{}
}

// ZapRecoveryHandler todo
type ZapRecoveryHandler struct {
	log *zap.Logger
}

// SetLogger todo
func (h *ZapRecoveryHandler) SetLogger(l *zap.Logger) *ZapRecoveryHandler {
	h.log = l
	return h
}

// Handle todo
func (h *ZapRecoveryHandler) Handle(ctx context.Context, p interface{}) error {
	stack := zap.Stack("stack").String

	if h.log != nil {
		h.log.Error(RecoveryExplanation, zap.Any("panic", p), zap.Any("stack", stack))
		return nil
	}

	log.Println(RecoveryExplanation, p, stack)
	return nil
}
