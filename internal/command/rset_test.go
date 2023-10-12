package command

import (
	"context"
	"testing"

	"github.com/Haya372/smtp-server/internal/mock"
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

	s := mock.NewMockSession(ctrl)
	s.EXPECT().Response(gomock.Eq(CodeOk), gomock.Eq(MsgOk)).Times(1)
	s.EXPECT().Reset().Times(1)

	target.HandleCommand(context.TODO(), s, make([]string, 0))
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
			s := mock.NewInitializedMockSession(ctrl, mock.SessionMockParam{})

			s.EXPECT().Response(test.code, test.msg).Times(1)
			s.EXPECT().Reset().Times(0)

			target := NewRsetHandler(log)
			target.HandleCommand(context.TODO(), s, test.arg)
		})
	}
}
