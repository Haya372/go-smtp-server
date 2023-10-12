package command

import (
	"context"

	"github.com/Haya372/hlog"
	"github.com/Haya372/smtp-server/internal/session"
)

type noopHandler struct {
	log hlog.Logger
}

func (h *noopHandler) Command() string {
	return NOOP
}

func (h *noopHandler) HandleCommand(ctx context.Context, s session.Session, arg []string) error {
	s.Response(CodeOk, MsgOk)
	return nil
}

func NewNoopHandler(log hlog.Logger) CommandHandler {
	return &noopHandler{
		log: log,
	}
}
