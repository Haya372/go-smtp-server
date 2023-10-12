package command

import (
	"context"

	"github.com/Haya372/hlog"
	"github.com/Haya372/smtp-server/internal/session"
)

type rsetHandler struct {
	log hlog.Logger
}

func (h *rsetHandler) Command() string {
	return RSET
}

func (h *rsetHandler) HandleCommand(ctx context.Context, s session.Session, arg []string) error {
	// RSET is not permit parameters
	if len(arg) > 0 {
		s.Response(CodeSyntaxError, MsgSyntaxError)
		return nil
	}

	s.Reset()
	s.Response(CodeOk, MsgOk)
	return nil
}

func NewRsetHandler(log hlog.Logger) CommandHandler {
	return &rsetHandler{
		log: log,
	}
}
