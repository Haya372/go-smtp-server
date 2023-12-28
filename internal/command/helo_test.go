package command

import (
	"context"
	"os"
	"testing"

	"github.com/Haya372/smtp-server/internal/mock"
	"github.com/Haya372/smtp-server/internal/session"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHelo_Command(t *testing.T) {
	target := NewHeloHandler(nil)
	assert.Equal(t, HELO, target.Command())
}

func TestHelo_Err(t *testing.T) {
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
			name: "empty argument",
			arg:  []string{},
			code: CodeSyntaxError,
			msg:  MsgSyntaxError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := session.NewMockSession(ctrl)

			s.ExpectResponse(test.code, test.msg)

			target := NewHeloHandler(log)
			target.HandleCommand(context.TODO(), s.Session, test.arg)
		})
	}
}

func TestHelo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock.NewInitializedMockLogger(ctrl)

	target := NewHeloHandler(log)

	s := session.NewMockSession(ctrl)

	arg := []string{"test"}

	hostname, _ := os.Hostname()
	s.ExpectResponse(CodeOk, hostname)

	target.HandleCommand(context.TODO(), s.Session, arg)
	assert.Equal(t, "test", s.Session.SenderDomain)
	assert.Nil(t, s.Session.EnvelopeFrom)
	assert.Empty(t, s.Session.EnvelopeTo)
	assert.Empty(t, s.Session.RawData)
}
