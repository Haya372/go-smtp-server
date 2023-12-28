package command

import (
	"context"
	"net/mail"
	"testing"

	"github.com/Haya372/smtp-server/internal/config"
	"github.com/Haya372/smtp-server/internal/mock"
	"github.com/Haya372/smtp-server/internal/session"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMail_Command(t *testing.T) {
	conf := &config.SmtpConfig{}
	target := NewMailHandler(nil, conf)
	assert.Equal(t, MAIL, target.Command())
}

func TestMail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock.NewInitializedMockLogger(ctrl)

	tests := []struct {
		name                      string
		arg                       []string
		conf                      *config.SmtpConfig
		expectEnvelopeFromAddress string
	}{
		{
			name:                      "no param",
			arg:                       []string{"from:<from@example.com>"},
			expectEnvelopeFromAddress: "<from@example.com>",
		},
		{
			name: "empty from address",
			arg:  []string{"from:<>"},
		},
		{
			name: "with SIZE param",
			arg:  []string{"from:<from@example.com>", "SIZE=100"},
			conf: &config.SmtpConfig{
				EnableSize:  true,
				MaxMailSize: 1000,
			},
			expectEnvelopeFromAddress: "<from@example.com>",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := session.NewMockSession(ctrl)
			s.Session.SenderDomain = "example.com"
			s.ExpectResponse(CodeOk, MsgOk)

			target := NewMailHandler(log, test.conf)
			target.HandleCommand(context.TODO(), s.Session, test.arg)
			var expect *mail.Address
			if len(test.expectEnvelopeFromAddress) != 0 {
				expect, _ = mail.ParseAddress(test.expectEnvelopeFromAddress)
			} else {
				expect = &mail.Address{}
			}
			assert.Equal(t, s.Session.EnvelopeFrom, expect)
		})
	}
}

func TestMail_Err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock.NewInitializedMockLogger(ctrl)

	tests := []struct {
		name          string
		arg           []string
		conf          *config.SmtpConfig
		senderDomain  string
		alreadyCalled bool
		code          int
		msg           string
	}{
		{
			name: "hello not called",
			arg:  []string{"from:<from@example.com>"},
			code: CodeBadSequence,
			msg:  MsgBadSequence,
		},
		{
			name:          "mail already called",
			arg:           []string{"from:<from@example.com>"},
			senderDomain:  "example.com",
			alreadyCalled: true,
			code:          CodeBadSequence,
			msg:           MsgBadSequence,
		},
		{
			name:         "argument is empty",
			senderDomain: "example.com",
			code:         CodeSyntaxError,
			msg:          MsgSyntaxError,
		},
		{
			name:         "invalid from address",
			senderDomain: "example.com",
			arg:          []string{"from:from@example.com>"},
			code:         CodeSyntaxError,
			msg:          MsgSyntaxError,
		},
		{
			name: "param error '=' not found",
			arg:  []string{"from:<from@example.com>", "SIZE100"},
			conf: &config.SmtpConfig{
				EnableSize:  true,
				MaxMailSize: 1000,
			},
			senderDomain: "example.com",
			code:         CodeOptionParamNotRecognized,
			msg:          MsgOptionParamNotRecognized,
		},
		{
			name: "param error SIZE value not integer",
			arg:  []string{"from:<from@example.com>", "SIZE=hoge"},
			conf: &config.SmtpConfig{
				EnableSize:  true,
				MaxMailSize: 1000,
			},
			senderDomain: "example.com",
			code:         CodeArgumentSyntaxError,
			msg:          MsgArgumentSyntaxError,
		},
		{
			name: "message size exceed limit",
			arg:  []string{"from:<from@example.com>", "SIZE=1000000000"},
			conf: &config.SmtpConfig{
				EnableSize:  true,
				MaxMailSize: 1000,
			},
			senderDomain: "example.com",
			code:         CodeAborted,
			msg:          MsgAborted,
		},
		{
			name: "size option disabled",
			arg:  []string{"from:<from@example.com>", "SIZE=1000000000"},
			conf: &config.SmtpConfig{
				EnableSize: false,
			},
			senderDomain: "example.com",
			code:         CodeCommandParamNotImplemented,
			msg:          MsgCommandParamNotImplemented,
		},
		{
			name: "unknown option",
			arg:  []string{"from:<from@example.com>", "UNKNOWN=hoge"},
			conf: &config.SmtpConfig{
				EnableSize:  true,
				MaxMailSize: 1000,
			},
			senderDomain: "example.com",
			code:         CodeCommandParamNotImplemented,
			msg:          MsgCommandParamNotImplemented,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := session.NewMockSession(ctrl)
			s.Session.SenderDomain = test.senderDomain
			if test.alreadyCalled {
				s.Session.EnvelopeFrom = &mail.Address{Address: "test@example.com"}
			}

			s.ExpectResponse(test.code, test.msg)

			target := NewMailHandler(log, test.conf)

			target.HandleCommand(context.TODO(), s.Session, test.arg)
		})
	}
}
