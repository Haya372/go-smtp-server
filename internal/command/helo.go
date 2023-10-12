package command

import (
	"context"

	"github.com/Haya372/hlog"
	"github.com/Haya372/smtp-server/internal/session"
)

type heloHandler struct {
	log hlog.Logger
}

func (h *heloHandler) Command() string {
	return HELO
}

func (h *heloHandler) HandleCommand(ctx context.Context, s session.Session, arg []string) error {
	if len(arg) == 0 {
		s.Response(CodeSyntaxError, MsgSyntaxError)
		return nil
	}

	// when helo command is called, session state should be initialized
	s.Reset()

	s.SetSenderDomain(arg[0])

	// TODO: ホスト名に変更する
	s.Response(CodeOk, "localhost")
	return nil
}

func NewHeloHandler(log hlog.Logger) CommandHandler {
	return &heloHandler{
		log: log,
	}
}
