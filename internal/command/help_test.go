package command

import (
	"context"
	"strings"
	"testing"

	"github.com/Haya372/smtp-server/internal/mock"
	"github.com/Haya372/smtp-server/internal/session"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCommand_Help(t *testing.T) {
	target := NewHelpHandler(nil)
	assert.Equal(t, HELP, target.Command())
}

func TestHelp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock.NewInitializedMockLogger(ctrl)

	target := NewHelpHandler(log)

	s := session.NewMockSession(ctrl)

	s.ExpectResponseLine(CodeHelp, MsgHelp)
	supportCommands := []string{
		HELO, EHLO, MAIL, RCPT, DATA, QUIT, RSET, NOOP, HELP,
	}
	respStr := strings.ToUpper(strings.Join(supportCommands, " "))
	s.ExpectResponse(CodeHelp, respStr)

	target.HandleCommand(context.TODO(), s.Session, make([]string, 0))
}
