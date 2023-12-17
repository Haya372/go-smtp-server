package command

import (
	"context"

	"github.com/Haya372/hlog"
	"github.com/Haya372/smtp-server/internal/config"
	"github.com/Haya372/smtp-server/internal/session"
)

type startTlsHandler struct {
	log  hlog.Logger
	conf *config.TlsConfig
}

func (h *startTlsHandler) Command() string {
	return STARTTLS
}

func (h startTlsHandler) HandleCommand(ctx context.Context, s session.Session, arg []string) error {
	if s.IsTls() {
		s.Response(CodeBadSequence, MsgAlreadyTls)
		return nil
	}

	s.Response(CodeGreet, MsgGoAhead)
	if err := s.ConvertToTls(h.conf.TlsConfig); err != nil {
		h.log.Errorf("[%d] tls error, err=%v", s.Id(), err)
		s.Response(CodeTransactionFail, MsgTransactionFail)
		return err
	}

	s.Reset()
	return nil
}

func NewStartTlsHandler(log hlog.Logger, conf *config.TlsConfig) CommandHandler {
	return &startTlsHandler{
		log:  log,
		conf: conf,
	}
}
