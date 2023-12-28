package command

import (
	"context"
	"testing"

	"github.com/Haya372/smtp-server/internal/mock"
	"github.com/Haya372/smtp-server/internal/session"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestQuit_Command(t *testing.T) {
	target := NewQuitHandler(nil)
	assert.Equal(t, QUIT, target.Command())
}

func TestQuit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock.NewInitializedMockLogger(ctrl)

	target := NewQuitHandler(log)

	s := session.NewMockSession(ctrl)
	s.ExpectResponse(CodeQuit, MsgQuit)

	target.HandleCommand(context.TODO(), s.Session, make([]string, 0))
	assert.True(t, s.Session.ShouldClose)
}

func TestQuit_Err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock.NewInitializedMockLogger(ctrl)

	tests := []struct {
		name string
		arg  []string
		code int
		msg  string
	}{
		{
			name: "with parameter",
			arg:  []string{"hoge"},
			code: CodeSyntaxError,
			msg:  MsgSyntaxError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := session.NewMockSession(ctrl)

			s.ExpectResponse(test.code, test.msg)

			target := NewQuitHandler(log)
			target.HandleCommand(context.TODO(), s.Session, test.arg)
			assert.False(t, s.Session.ShouldClose)
		})
	}
}
