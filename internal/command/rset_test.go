package command

import (
	"context"
	"testing"

	"github.com/Haya372/smtp-server/internal/mock"
	"github.com/Haya372/smtp-server/internal/session"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRset_Command(t *testing.T) {
	target := NewRsetHandler(nil)
	assert.Equal(t, RSET, target.Command())
}

func TestRset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock.NewInitializedMockLogger(ctrl)

	target := NewRsetHandler(log)

	s := session.NewMockSession(ctrl)
	s.ExpectResponse(CodeOk, MsgOk)

	target.HandleCommand(context.TODO(), s.Session, make([]string, 0))
	assert.Empty(t, s.Session.SenderDomain)
	assert.Nil(t, s.Session.EnvelopeFrom)
	assert.Empty(t, s.Session.EnvelopeTo)
	assert.Empty(t, s.Session.RawData)
}

func TestRset_Err(t *testing.T) {
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
			s.Session.SenderDomain = "test"

			s.ExpectResponse(test.code, test.msg)

			target := NewRsetHandler(log)
			target.HandleCommand(context.TODO(), s.Session, test.arg)
			assert.NotEmpty(t, s.Session.SenderDomain)
		})
	}
}
