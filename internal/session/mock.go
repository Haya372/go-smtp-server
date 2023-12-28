package session

import (
	"bufio"
	"fmt"
	"net/textproto"
	"strings"

	"github.com/Haya372/smtp-server/internal/mock/oss"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

type MockSession struct {
	Session *Session

	ctrl   *gomock.Controller
	Writer *oss.MockWriter
}

func NewMockSession(ctrl *gomock.Controller) *MockSession {
	writer := oss.NewMockWriter(ctrl)

	return &MockSession{
		Session: &Session{
			Id:     uuid.New(),
			writer: *textproto.NewWriter(bufio.NewWriter(writer)),
		},
		ctrl:   ctrl,
		Writer: writer,
	}
}

func (s *MockSession) ExpectResponse(code int, msg string) {
	s.expectResponseStr(fmt.Sprintf("%d %s\r\n", code, msg))
}

func (s *MockSession) ExpectResponseLine(code int, msg string) {
	s.expectResponseStr(fmt.Sprintf("%d-%s\r\n", code, msg))
}

func (s *MockSession) expectResponseStr(msg string) {
	s.Writer.EXPECT().Write([]byte(msg)).Return(len([]byte(msg)), nil)
}

func (s *MockSession) ExpectReadLine(line string, err error) {
	if err == nil {
		s.Session.reader = *textproto.NewReader(bufio.NewReader(strings.NewReader(line)))
	} else {
		reader := oss.NewMockReader(s.ctrl)
		s.Session.reader = *textproto.NewReader(bufio.NewReader(reader))

		reader.EXPECT().Read(gomock.Any()).Return(0, err)
	}
}
