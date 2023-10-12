package connection

import (
	"context"
	"errors"
	"io"
	"net"
	"strings"

	"github.com/Haya372/hlog"
	"github.com/Haya372/smtp-server/internal/command"
	"github.com/Haya372/smtp-server/internal/session"
)

type SessionHandler struct {
	log             hlog.Logger
	commandHandlers map[string]command.CommandHandler
}

func (h *SessionHandler) handleCommand(ctx context.Context, s session.Session, line string) {
	cmd := strings.ToLower(strings.Fields(line)[0])
	cmdHandler := h.commandHandlers[cmd]

	if cmdHandler != nil {
		cmdHandler.HandleCommand(ctx, s, strings.Fields(line)[1:])
	} else {
		h.log.Errorf("[%s] receive illegal command %s.", s.Id(), cmd)
		s.Response(command.CodeCommandNotImplemented, command.MsgCommandNotImplemented)
	}
}

func (h *SessionHandler) HandleSession(ctx context.Context, s session.Session) {
	h.log.Debugf("[%s] receive connection", s.Id())
	s.Response(command.CodeGreet, command.MsgGreet)
	defer s.Close()

	for {
		line, err := s.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
				// closed by client
				h.log.Infof("[%s] connection closed.", s.Id())
			} else {
				h.log.WithError(err).Errorf("[%s] could not read line. %v", s.Id(), err)
				s.Response(command.CodeServiceNotAvailable, command.MsgServiceNotAvailable)
			}
			return
		}
		h.log.Debugf("[%s] received line: %s", s.Id(), line)

		h.handleCommand(ctx, s, line)

		if s.IsCloseImmediately() {
			break
		}
	}
}

func NewSessionHandler(log hlog.Logger, cmdHandlers []command.CommandHandler) SessionHandler {
	commandHandlers := make(map[string]command.CommandHandler, 0)

	for _, cmdHandler := range cmdHandlers {
		commandHandlers[cmdHandler.Command()] = cmdHandler
	}

	return SessionHandler{
		log:             log,
		commandHandlers: commandHandlers,
	}
}
