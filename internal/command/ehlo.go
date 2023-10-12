package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/Haya372/hlog"
	"github.com/Haya372/smtp-server/internal/session"
)

type ehloHandler struct {
	log hlog.Logger
}

func (h *ehloHandler) Command() string {
	return EHLO
}

func (h *ehloHandler) HandleCommand(ctx context.Context, s session.Session, arg []string) error {
	if len(arg) == 0 {
		s.Response(CodeSyntaxError, MsgSyntaxError)
		return nil
	}

	// when ehlo command is called, session state should be initialized
	s.Reset()

	s.SetSenderDomain(arg[0])

	// TODO: ESMTPのレスポンス定義
	// TODO: 設定から読み込む
	s.ResponseLine(fmt.Sprintf("250-%s greets %s", "localhost", arg[0]))
	s.ResponseLine("250-PIPELINING")
	s.ResponseLine("250-8BITMIME")
	s.ResponseLine(fmt.Sprintf("250-SIZE %d", 1048576)) // 1MB
	s.Response(CodeOk, strings.ToUpper(HELP))
	return nil
}

func NewEhloHandler(log hlog.Logger) CommandHandler {
	return &ehloHandler{
		log: log,
	}
}
