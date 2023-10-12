package command

import (
	"context"
	"net/mail"
	"testing"

	"github.com/Haya372/smtp-server/internal/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRcpt_Command(t *testing.T) {
	target := NewRcptHandler(nil)
	assert.Equal(t, RCPT, target.Command())
}

func TestRcpt_Err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock.NewInitializedMockLogger(ctrl)

	tests := []struct {
		name         string
		arg          []string
		envelopeFrom string
		code         int
		msg          string
	}{
		{
			name: "mail not called",
			arg:  []string{"to:<to@example.com>"},
			code: CodeBadSequence,
			msg:  MsgBadSequence,
		},
		{
			name:         "argument is empty",
			envelopeFrom: "from@example.com",
			code:         CodeSyntaxError,
			msg:          MsgSyntaxError,
		},
		{
			name:         "invalid to address",
			envelopeFrom: "from@example.com",
			arg:          []string{"to:to@example.com>"},
			code:         CodeSyntaxError,
			msg:          MsgSyntaxError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := mock.NewInitializedMockSession(ctrl, mock.SessionMockParam{
				EnvelopeFrom: test.envelopeFrom,
			})

			s.EXPECT().Response(gomock.Eq(test.code), gomock.Eq(test.msg))

			target := NewRcptHandler(log)
			target.HandleCommand(context.TODO(), s, test.arg)
		})
	}
}

func TestRcpt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock.NewInitializedMockLogger(ctrl)

	tests := []struct {
		name               string
		arg                []string
		expectedEnvelopeTo string
	}{
		{
			name:               "no param",
			arg:                []string{"to:<to@example.com>"},
			expectedEnvelopeTo: "<to@example.com>",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := mock.NewInitializedMockSession(ctrl, mock.SessionMockParam{
				EnvelopeFrom: "from@example.com",
			})

			expect, _ := mail.ParseAddress(test.expectedEnvelopeTo)
			s.EXPECT().AddEnvelopeTo(mock.NewAddressMatcher(*expect)).Times(1)
			s.EXPECT().Response(gomock.Eq(CodeOk), gomock.Eq(MsgOk)).Times(1)

			target := NewRcptHandler(log)
			target.HandleCommand(context.TODO(), s, test.arg)
		})
	}
}
