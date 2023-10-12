package command

import (
	"context"
	"errors"
	"net/mail"
	"strconv"
	"strings"

	"github.com/Haya372/hlog"
	"github.com/Haya372/smtp-server/internal/session"
)

type mailHandler struct {
	log hlog.Logger
}

func (h *mailHandler) Command() string {
	return MAIL
}

func (h *mailHandler) HandleCommand(ctx context.Context, s session.Session, arg []string) error {
	// helo or ehlo command should be called
	if len(s.SenderDomain()) == 0 {
		s.Response(CodeBadSequence, MsgBadSequence)
		return nil
	}

	if s.EnvelopeFrom() != nil {
		s.Response(CodeBadSequence, MsgBadSequence)
		return nil
	}

	if len(arg) == 0 {
		s.Response(CodeSyntaxError, MsgSyntaxError)
		return nil
	}

	addr := strings.Replace(arg[0], "from:", "", 1)
	addr = strings.Replace(addr, "FROM:", "", 1)

	// check ESMTP arguments
	for _, line := range arg[1:] {
		keyVal := strings.Split(line, "=")
		if len(keyVal) != 2 {
			h.log.Errorf("[%s] failed to recognized option %s", s.Id(), line)
			s.Response(CodeOptionParamNotRecognized, MsgOptionParamNotRecognized)
			return nil
		}

		opt := keyVal[0]
		val := keyVal[1]
		var err error
		switch opt {
		case "SIZE":
			err = h.handleSizeOption(ctx, s, val)
		default:
			err = errors.New("option not implemented")
			s.Response(CodeCommandParamNotImplemented, MsgCommandParamNotImplemented)
		}
		if err != nil {
			h.log.WithError(err).Errorf("[%s] failed to handle option %s", s.Id(), opt)
			return err
		}
	}

	if addr == "<>" {
		// Envelope From is null
		s.SetEnvelopeFrom(&mail.Address{})
	} else {
		address, err := mail.ParseAddress(addr)
		if err != nil {
			h.log.WithError(err).Debugf("[%s] failed to parse address %s", s.Id(), arg[0])
			s.Response(CodeSyntaxError, MsgSyntaxError)
			return nil
		}

		s.SetEnvelopeFrom(address)
	}

	s.Response(CodeOk, MsgOk)
	return nil
}

func (h *mailHandler) handleSizeOption(ctx context.Context, s session.Session, arg string) error {
	size, err := strconv.Atoi(arg)
	if err != nil {
		s.Response(CodeArgumentSyntaxError, MsgArgumentSyntaxError)
		return err
	}
	if size > 1048576 {
		s.Response(CodeAborted, MsgAborted)
		return errors.New("message size exceed limit")
	}
	return nil
}

func NewMailHandler(log hlog.Logger) CommandHandler {
	return &mailHandler{
		log: log,
	}
}
