package command

import (
	"context"
	"testing"

	"github.com/Haya372/smtp-server/internal/mock"
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
			s := mock.NewInitializedMockSession(ctrl, mock.SessionMockParam{})

			s.EXPECT().Response(gomock.Eq(test.code), gomock.Eq(test.msg)).Times(1)

			target := NewHeloHandler(log)
			target.HandleCommand(context.TODO(), s, test.arg)
		})
	}
}

func TestHelo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock.NewInitializedMockLogger(ctrl)

	target := NewHeloHandler(log)

	s := mock.NewMockSession(ctrl)

	arg := []string{"test"}

	s.EXPECT().Reset().Times(1)
	s.EXPECT().SetSenderDomain("test")
	// TODO: 第２引数をホスト名に修正する
	s.EXPECT().Response(gomock.Eq(CodeOk), gomock.Any())

	target.HandleCommand(context.TODO(), s, arg)
}
