package command

import (
	"context"

	"github.com/Haya372/hlog"
	"github.com/Haya372/smtp-server/internal/session"
)

type quitHandler struct {
	log hlog.Logger
}

func (h *quitHandler) Command() string {
	return QUIT
}

func (h *quitHandler) HandleCommand(ctx context.Context, s *session.Session, arg []string) error {
	// QUIT is not permit parameters
	if len(arg) > 0 {
		s.Response(CodeSyntaxError, MsgSyntaxError)
		return nil
	}

	s.Response(CodeQuit, MsgQuit)
	s.ShouldClose = true
	return nil
}

func NewQuitHandler(log hlog.Logger) CommandHandler {
	return &quitHandler{
		log: log,
	}
}
