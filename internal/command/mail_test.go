package command

import (
	"context"
	"net/mail"
	"testing"

	"github.com/Haya372/smtp-server/internal/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMail_Command(t *testing.T) {
	target := NewMailHandler(nil)
	assert.Equal(t, MAIL, target.Command())
}

func TestMail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock.NewInitializedMockLogger(ctrl)

	tests := []struct {
		name                      string
		arg                       []string
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
			name:                      "with SIZE param",
			arg:                       []string{"from:<from@example.com>", "SIZE=100"},
			expectEnvelopeFromAddress: "<from@example.com>",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := mock.NewInitializedMockSession(ctrl, mock.SessionMockParam{
				SenderDomain: "example.com",
			})

			var expect *mail.Address
			if len(test.expectEnvelopeFromAddress) != 0 {
				expect, _ = mail.ParseAddress(test.expectEnvelopeFromAddress)
			} else {
				expect = &mail.Address{}
			}

			s.EXPECT().SetEnvelopeFrom(expect).Times(1)
			s.EXPECT().Response(gomock.Eq(CodeOk), gomock.Eq(MsgOk)).Times(1)

			target := NewMailHandler(log)
			target.HandleCommand(context.TODO(), s, test.arg)
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
			name:         "param error '=' not found",
			arg:          []string{"from:<from@example.com>", "SIZE100"},
			senderDomain: "example.com",
			code:         CodeOptionParamNotRecognized,
			msg:          MsgOptionParamNotRecognized,
		},
		{
			name:         "param error SIZE value not integer",
			arg:          []string{"from:<from@example.com>", "SIZE=hoge"},
			senderDomain: "example.com",
			code:         CodeArgumentSyntaxError,
			msg:          MsgArgumentSyntaxError,
		},
		{
			name:         "message size exceed limit",
			arg:          []string{"from:<from@example.com>", "SIZE=1000000000"},
			senderDomain: "example.com",
			code:         CodeAborted,
			msg:          MsgAborted,
		},
		{
			name:         "unknown option",
			arg:          []string{"from:<from@example.com>", "UNKNOWN=hoge"},
			senderDomain: "example.com",
			code:         CodeCommandParamNotImplemented,
			msg:          MsgCommandParamNotImplemented,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var envelopeFrom string
			if test.alreadyCalled {
				envelopeFrom = "test@example.com"
			}
			s := mock.NewInitializedMockSession(ctrl, mock.SessionMockParam{
				SenderDomain: test.senderDomain,
				EnvelopeFrom: envelopeFrom,
			})

			s.EXPECT().Response(gomock.Eq(test.code), gomock.Eq(test.msg)).Times(1)

			target := NewMailHandler(log)

			target.HandleCommand(context.TODO(), s, test.arg)
		})
	}
}
