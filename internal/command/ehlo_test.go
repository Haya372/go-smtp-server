package command

import (
	"context"
	"testing"

	"github.com/Haya372/smtp-server/internal/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestEhlo_Command(t *testing.T) {
	target := NewEhloHandler(nil)

	assert.Equal(t, EHLO, target.Command())
}

func TestEhlo_Err(t *testing.T) {
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
			s := mock.NewInitializedMockSession(ctrl, mock.SessionMockParam{})

			s.EXPECT().Response(gomock.Eq(test.code), gomock.Eq(test.msg)).Times(1)

			target := NewEhloHandler(log)
			target.HandleCommand(context.TODO(), s, test.arg)
		})
	}
}

func TestEhlo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock.NewInitializedMockLogger(ctrl)

	target := NewEhloHandler(log)

	s := mock.NewMockSession(ctrl)

	arg := []string{"test"}

	s.EXPECT().Reset().Times(1)
	s.EXPECT().SetSenderDomain("test")
	s.EXPECT().ResponseLine(gomock.Any()).AnyTimes()
	s.EXPECT().Response(gomock.Eq(CodeOk), gomock.Any())

	target.HandleCommand(context.TODO(), s, arg)
}
