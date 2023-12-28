package connection

import (
	"context"
	"errors"
	"io"
	"net"
	"testing"

	"github.com/Haya372/smtp-server/internal/command"
	"github.com/Haya372/smtp-server/internal/mock"
	"github.com/Haya372/smtp-server/internal/mock/oss"
	"github.com/Haya372/smtp-server/internal/session"
	"github.com/golang/mock/gomock"
)

func TestSessionHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock.NewInitializedMockLogger(ctrl)

	tests := []struct {
		name  string
		setup func(s *session.MockSession, h *mock.MockCommandHandler)
		close bool
	}{
		{
			name: "no error",
			setup: func(s *session.MockSession, h *mock.MockCommandHandler) {
				h.EXPECT().HandleCommand(gomock.Any(), gomock.Any(), []string{"example.com"}).Times(1)
				s.ExpectReadLine("helo example.com", nil)
			},
		},
		{
			name: "read line error (io.EOF)",
			setup: func(s *session.MockSession, h *mock.MockCommandHandler) {
				s.ExpectReadLine("", io.EOF)
			},
		},
		{
			name: "read line error (net.ErrClosed)",
			setup: func(s *session.MockSession, h *mock.MockCommandHandler) {
				s.ExpectReadLine("", net.ErrClosed)
			},
		},
		{
			name: "read line error (others)",
			setup: func(s *session.MockSession, h *mock.MockCommandHandler) {
				s.ExpectReadLine("", errors.New("test error"))
				s.ExpectResponse(command.CodeServiceNotAvailable, command.MsgServiceNotAvailable)
			},
			close: true,
		},
		{
			name: "command not implemented",
			setup: func(s *session.MockSession, h *mock.MockCommandHandler) {
				s.ExpectReadLine("test", nil)
				s.ExpectResponse(command.CodeCommandNotImplemented, command.MsgCommandNotImplemented)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := session.NewMockSession(ctrl)
			h := mock.NewInitializedMockCommandHandler(ctrl, command.HELO)

			s.ExpectResponse(command.CodeGreet, command.MsgGreet)

			conn := oss.NewMockConn(ctrl)
			conn.EXPECT().Close().Times(1)
			s.Session.Conn = conn
			target := NewSessionHandler(log, []command.CommandHandler{h})

			if test.setup != nil {
				test.setup(s, h)
			}

			target.HandleSession(context.TODO(), s.Session)
		})
	}
}
