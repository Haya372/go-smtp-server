package command

import (
	"context"
	"errors"

	"github.com/Haya372/hlog"
	"github.com/Haya372/smtp-server/internal/service"
	"github.com/Haya372/smtp-server/internal/session"
)

type dataHandler struct {
	log         hlog.Logger
	authService service.AuthService
}

func (h *dataHandler) Command() string {
	return DATA
}

func (h *dataHandler) HandleCommand(ctx context.Context, s session.Session, arg []string) error {
	// DATA is not permit parameters
	if len(arg) > 0 {
		s.Response(CodeSyntaxError, MsgSyntaxError)
		return nil
	}

	// rcpt command should be called
	if len(s.EnvelopeTo()) == 0 {
		s.Response(CodeBadSequence, MsgBadSequence)
		return nil
	}

	s.Response(CodeStartInput, MsgStartInput)
	rawData, err := s.ReadRawData()
	if err != nil {
		h.log.WithError(err).Errorf("[%s] data reading error.", s.Id())
		s.Response(CodeTransactionFail, MsgTransactionFail)
		return err
	}

	// TODO: 設定から最大のサイズを取得する
	if len(rawData) > 1048576 {
		s.Response(CodeAborted, MsgAborted)
		return errors.New("message size exceed limit")
	}

	h.log.Debugf("[%s] mail data received.\n----------\n%s----------", s.Id(), string(rawData))

	s.Response(CodeOk, MsgOk)
	s.Reset()
	return nil
}

func NewDataHandler(log hlog.Logger, authService service.AuthService) CommandHandler {
	return &dataHandler{
		log:         log,
		authService: authService,
	}
}
