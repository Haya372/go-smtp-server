package connection

import (
	"context"
	"errors"
	"io"
	"net"
	"testing"

	"github.com/Haya372/smtp-server/internal/command"
	"github.com/Haya372/smtp-server/internal/mock"
	"github.com/golang/mock/gomock"
)

func TestSessionHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock.NewInitializedMockLogger(ctrl)

	tests := []struct {
		name  string
		setup func(s *mock.MockSession, h *mock.MockCommandHandler)
	}{
		{
			name: "no error",
			setup: func(s *mock.MockSession, h *mock.MockCommandHandler) {
				h.EXPECT().HandleCommand(gomock.Any(), gomock.Any(), []string{"example.com"}).Times(1)
				s.EXPECT().ReadLine().Return("helo example.com", nil).Times(1)
				s.EXPECT().IsCloseImmediately().Return(true).Times(1)
			},
		},
		{
			name: "read line error (io.EOF)",
			setup: func(s *mock.MockSession, h *mock.MockCommandHandler) {
				s.EXPECT().ReadLine().Return("", io.EOF).Times(1)
			},
		},
		{
			name: "read line error (net.ErrClosed)",
			setup: func(s *mock.MockSession, h *mock.MockCommandHandler) {
				s.EXPECT().ReadLine().Return("", net.ErrClosed).Times(1)
			},
		},
		{
			name: "read line error (others)",
			setup: func(s *mock.MockSession, h *mock.MockCommandHandler) {
				s.EXPECT().ReadLine().Return("", errors.New("test error")).Times(1)
				s.EXPECT().Response(gomock.Eq(command.CodeServiceNotAvailable), gomock.Eq(command.MsgServiceNotAvailable)).Times(1)
			},
		},
		{
			name: "command not implemented",
			setup: func(s *mock.MockSession, h *mock.MockCommandHandler) {
				s.EXPECT().ReadLine().Return("test", nil).Times(1)
				s.EXPECT().Response(gomock.Eq(command.CodeCommandNotImplemented), gomock.Eq(command.MsgCommandNotImplemented)).Times(1)
				s.EXPECT().IsCloseImmediately().Return(true).Times(1)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := mock.NewInitializedMockSession(ctrl, mock.SessionMockParam{})
			h := mock.NewInitializedMockCommandHandler(ctrl, command.HELO)

			s.EXPECT().Response(gomock.Eq(command.CodeGreet), gomock.Eq(command.MsgGreet)).Times(1)
			s.EXPECT().Close().Times(1)
			target := NewSessionHandler(log, []command.CommandHandler{h})

			if test.setup != nil {
				test.setup(s, h)
			}

			target.HandleSession(context.TODO(), s)
		})
	}
}
