package command

import (
	"context"
	"errors"
	"net/mail"
	"testing"

	"github.com/Haya372/smtp-server/internal/config"
	"github.com/Haya372/smtp-server/internal/mock"
	"github.com/Haya372/smtp-server/internal/session"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestData_Command(t *testing.T) {
	conf := &config.SmtpConfig{}
	target := NewDataHandler(nil, conf)

	assert.Equal(t, target.Command(), DATA)
}

func TestData_Err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock.NewInitializedMockLogger(ctrl)

	conf := &config.SmtpConfig{
		MaxMailSize: 1000,
	}

	tests := []struct {
		name      string
		arg       []string
		setupFunc func(s *session.MockSession)
		code      int
		msg       string
	}{
		{
			name: "rcpt not called",
			code: CodeBadSequence,
			msg:  MsgBadSequence,
		},
		{
			name: "read raw data err",
			setupFunc: func(s *session.MockSession) {
				s.Session.EnvelopeTo = []mail.Address{{Address: "to@example.com"}}
				s.ExpectResponse(CodeStartInput, MsgStartInput)
				s.ExpectReadLine("", errors.New("test error"))

			},
			code: CodeTransactionFail,
			msg:  MsgTransactionFail,
		},
		{
			name: "message size exceeds limit",
			setupFunc: func(s *session.MockSession) {
				s.Session.EnvelopeTo = []mail.Address{{Address: "to@example.com"}}
				s.ExpectResponse(CodeStartInput, MsgStartInput)
				data := "Subject: test\r\n\r\n"
				for i := 0; i < conf.MaxMailSize; i++ {
					data += "a"
				}
				data += "\r\n.\r\n"
				s.ExpectReadLine(data, nil)
			},
			code: CodeAborted,
			msg:  MsgAborted,
		},
		{
			name: "with parameter",
			setupFunc: func(s *session.MockSession) {
				s.Session.EnvelopeTo = []mail.Address{{Address: "to@example.com"}}
			},
			arg:  []string{"hoge"},
			code: CodeSyntaxError,
			msg:  MsgSyntaxError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := session.NewMockSession(ctrl)

			if test.setupFunc != nil {
				test.setupFunc(s)
			}
			s.ExpectResponse(test.code, test.msg)

			target := NewDataHandler(log, conf)
			target.HandleCommand(context.TODO(), s.Session, test.arg)
		})
	}
}

func TestData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock.NewInitializedMockLogger(ctrl)

	conf := &config.SmtpConfig{
		MaxMailSize: 1000,
	}

	target := NewDataHandler(log, conf)

	s := session.NewMockSession(ctrl)
	s.Session.EnvelopeTo = make([]mail.Address, 1)

	s.ExpectResponse(CodeStartInput, MsgStartInput)
	s.ExpectReadLine("Subject: test\r\n\r\n.\r\n", nil)
	s.ExpectResponse(CodeOk, MsgOk)

	target.HandleCommand(context.TODO(), s.Session, make([]string, 0))
	assert.Empty(t, s.Session.SenderDomain)
	assert.Nil(t, s.Session.EnvelopeFrom)
	assert.Empty(t, s.Session.EnvelopeTo)
	assert.Empty(t, s.Session.RawData)
}
