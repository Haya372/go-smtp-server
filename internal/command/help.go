package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/Haya372/hlog"
	"github.com/Haya372/smtp-server/internal/session"
)

type helpHandler struct {
	log hlog.Logger
}

func (h *helpHandler) Command() string {
	return HELP
}

func (h *helpHandler) HandleCommand(ctx context.Context, s *session.Session, arg []string) error {
	s.ResponseLine(fmt.Sprintf("%d-%s", CodeHelp, MsgHelp))
	supportCommands := []string{
		HELO, EHLO, MAIL, RCPT, DATA, QUIT, RSET, NOOP, HELP,
	}

	respStr := strings.ToUpper(strings.Join(supportCommands, " "))
	s.Response(CodeHelp, respStr)
	return nil
}

func NewHelpHandler(log hlog.Logger) CommandHandler {
	return &helpHandler{
		log: log,
	}
}
