package command

import (
	"context"
	"net/mail"
	"strings"

	"github.com/Haya372/hlog"
	"github.com/Haya372/smtp-server/internal/session"
)

type rcptHandler struct {
	log hlog.Logger
}

func (h *rcptHandler) Command() string {
	return RCPT
}

func (h *rcptHandler) HandleCommand(ctx context.Context, s session.Session, arg []string) error {
	// mail command should be called
	if s.EnvelopeFrom() == nil {
		s.Response(CodeBadSequence, MsgBadSequence)
		return nil
	}

	if len(arg) == 0 {
		s.Response(CodeSyntaxError, MsgSyntaxError)
		return nil
	}

	addr := strings.Replace(arg[0], "to:", "", 1)
	addr = strings.Replace(addr, "TO:", "", 1)

	// TODO: rcpt コマンドのオプションチェック(ESMTP)

	address, err := mail.ParseAddress(addr)
	if err != nil {
		h.log.WithError(err).Debugf("[%s] failed to parse address %s", s.Id(), arg[0])
		s.Response(CodeSyntaxError, MsgSyntaxError)
		return nil
	}

	s.AddEnvelopeTo(*address)

	s.Response(CodeOk, MsgOk)
	return nil
}

func NewRcptHandler(log hlog.Logger) CommandHandler {
	return &rcptHandler{
		log: log,
	}
}
