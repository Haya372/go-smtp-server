package command

import (
	"context"
	"fmt"
	"testing"

	"github.com/Haya372/smtp-server/internal/mock"
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

	s := mock.NewMockSession(ctrl)

	s.EXPECT().ResponseLine(fmt.Sprintf("%d-%s", CodeHelp, MsgHelp))
	s.EXPECT().Response(gomock.Eq(CodeHelp), gomock.Any()).Times(1)

	target.HandleCommand(context.TODO(), s, make([]string, 0))
}
