package command

import (
	"context"
	"errors"
	"testing"

	"github.com/Haya372/smtp-server/internal/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestData_Command(t *testing.T) {
	target := NewDataHandler(nil)

	assert.Equal(t, target.Command(), DATA)
}

func TestData_Err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock.NewInitializedMockLogger(ctrl)

	tests := []struct {
		name         string
		sessionParam mock.SessionMockParam
		arg          []string
		setupFunc    func(s *mock.MockSession)
		code         int
		msg          string
	}{
		{
			name:         "rcpt not called",
			sessionParam: mock.SessionMockParam{},
			code:         CodeBadSequence,
			msg:          MsgBadSequence,
		},
		{
			name: "read raw data err",
			sessionParam: mock.SessionMockParam{
				EnvelopeTo: []string{"to@example.com"},
			},
			setupFunc: func(s *mock.MockSession) {
				s.EXPECT().Response(CodeStartInput, MsgStartInput)
				s.EXPECT().ReadRawData().Return(nil, errors.New("test error")).Times(1)
			},
			code: CodeTransactionFail,
			msg:  MsgTransactionFail,
		},
		{
			name: "message size exceeds limit",
			sessionParam: mock.SessionMockParam{
				EnvelopeTo: []string{"to@example.com"},
			},
			setupFunc: func(s *mock.MockSession) {
				s.EXPECT().Response(CodeStartInput, MsgStartInput)
				s.EXPECT().ReadRawData().Return(make([]byte, 1000000000), nil).Times(1)
			},
			code: CodeAborted,
			msg:  MsgAborted,
		},
		{
			name: "with parameter",
			sessionParam: mock.SessionMockParam{
				EnvelopeTo: []string{"to@example.com"},
			},
			arg:  []string{"hoge"},
			code: CodeSyntaxError,
			msg:  MsgSyntaxError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := mock.NewInitializedMockSession(ctrl, test.sessionParam)

			if test.setupFunc != nil {
				test.setupFunc(s)
			}
			s.EXPECT().Response(gomock.Eq(test.code), gomock.Eq(test.msg)).Times(1)

			target := NewDataHandler(log)
			target.HandleCommand(context.TODO(), s, test.arg)
		})
	}
}

func TestData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock.NewInitializedMockLogger(ctrl)

	target := NewDataHandler(log)

	s := mock.NewInitializedMockSession(ctrl, mock.SessionMockParam{
		EnvelopeTo: []string{"to@example.com"},
	})

	s.EXPECT().Response(gomock.Eq(CodeStartInput), gomock.Eq(MsgStartInput)).Times(1)
	s.EXPECT().ReadRawData().Return([]byte("test"), nil).Times(1)
	s.EXPECT().Response(gomock.Eq(CodeOk), gomock.Eq(MsgOk)).Times(1)
	s.EXPECT().Reset().Times(1)
	target.HandleCommand(context.TODO(), s, make([]string, 0))
}
