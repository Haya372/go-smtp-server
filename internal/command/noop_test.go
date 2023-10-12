package command

import (
	"context"
	"testing"

	"github.com/Haya372/smtp-server/internal/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCommand_Noop(t *testing.T) {
	target := NewNoopHandler(nil)
	assert.Equal(t, NOOP, target.Command())
}

func TestNoop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock.NewInitializedMockLogger(ctrl)

	target := NewNoopHandler(log)

	s := mock.NewMockSession(ctrl)
	s.EXPECT().Response(gomock.Eq(CodeOk), gomock.Eq(MsgOk)).Times(1)

	target.HandleCommand(context.TODO(), s, make([]string, 0))
}
