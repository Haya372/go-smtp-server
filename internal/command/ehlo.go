package command

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Haya372/hlog"
	"github.com/Haya372/smtp-server/internal/config"
	"github.com/Haya372/smtp-server/internal/session"
)

type ehloHandler struct {
	log  hlog.Logger
	conf *config.SmtpConfig
}

func (h *ehloHandler) Command() string {
	return EHLO
}

func (h *ehloHandler) HandleCommand(ctx context.Context, s *session.Session, arg []string) error {
	if len(arg) == 0 {
		s.Response(CodeSyntaxError, MsgSyntaxError)
		return nil
	}

	// when ehlo command is called, session state should be initialized
	s.Reset()

	s.SenderDomain = arg[0]

	// TODO: ESMTPのレスポンス定義
	hostname, _ := os.Hostname()
	s.ResponseLine(fmt.Sprintf("%d-%s greets %s", CodeOk, hostname, arg[0]))
	if h.conf.EnablePipelining {
		s.ResponseLine(fmt.Sprintf("%d-PIPELINING", CodeOk))
	}
	if h.conf.Enable8BitMime {
		s.ResponseLine(fmt.Sprintf("%d-8BITMIME", CodeOk))
	}
	if h.conf.EnableSize {
		s.ResponseLine(fmt.Sprintf("%d-SIZE %d", CodeOk, h.conf.MaxMailSize))
	}
	if h.conf.EnableStartTls && !s.IsTls() {
		s.ResponseLine(fmt.Sprintf("%d-STARTTLS", CodeOk))
	}
	s.Response(CodeOk, strings.ToUpper(HELP))
	return nil
}

func NewEhloHandler(log hlog.Logger, conf *config.SmtpConfig) CommandHandler {
	return &ehloHandler{
		log:  log,
		conf: conf,
	}
}
