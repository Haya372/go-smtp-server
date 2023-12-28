package command

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Haya372/smtp-server/internal/config"
	"github.com/Haya372/smtp-server/internal/mock"
	"github.com/Haya372/smtp-server/internal/session"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestEhlo_Command(t *testing.T) {
	conf := &config.SmtpConfig{}
	target := NewEhloHandler(nil, conf)

	assert.Equal(t, EHLO, target.Command())
}

func TestEhlo_Err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock.NewInitializedMockLogger(ctrl)
	conf := &config.SmtpConfig{}

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
			s := session.NewMockSession(ctrl)

			s.ExpectResponse(test.code, test.msg)

			target := NewEhloHandler(log, conf)
			target.HandleCommand(context.TODO(), s.Session, test.arg)
		})
	}
}

func TestEhlo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock.NewInitializedMockLogger(ctrl)

	tests := []struct {
		name       string
		conf       *config.SmtpConfig
		setup      func(s *session.MockSession)
		alreadyTls bool
	}{
		{
			name: "no extension",
			conf: &config.SmtpConfig{},
			setup: func(s *session.MockSession) {
				hostname, _ := os.Hostname()
				s.ExpectResponseLine(CodeOk, fmt.Sprintf("%s greets %s", hostname, "test"))
			},
		},
		{
			name: "enable all",
			conf: &config.SmtpConfig{
				EnablePipelining: true,
				Enable8BitMime:   true,
				EnableSize:       true,
				EnableStartTls:   true,
				MaxMailSize:      1,
			},
			setup: func(s *session.MockSession) {
				hostname, _ := os.Hostname()
				s.ExpectResponseLine(CodeOk, fmt.Sprintf("%s greets %s", hostname, "test"))
				s.ExpectResponseLine(CodeOk, "PIPELINING")
				s.ExpectResponseLine(CodeOk, "8BITMIME")
				s.ExpectResponseLine(CodeOk, "SIZE 1")
				s.ExpectResponseLine(CodeOk, "STARTTLS")
			},
		},
		{
			name: "already tls",
			conf: &config.SmtpConfig{
				EnablePipelining: true,
				Enable8BitMime:   true,
				EnableSize:       true,
				EnableStartTls:   true,
				MaxMailSize:      1,
			},
			setup: func(s *session.MockSession) {
				hostname, _ := os.Hostname()
				s.ExpectResponseLine(CodeOk, fmt.Sprintf("%s greets %s", hostname, "test"))
				s.ExpectResponseLine(CodeOk, "PIPELINING")
				s.ExpectResponseLine(CodeOk, "8BITMIME")
				s.ExpectResponseLine(CodeOk, "SIZE 1")
			},
			alreadyTls: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			target := NewEhloHandler(log, test.conf)

			s := session.NewMockSession(ctrl)

			test.setup(s)
			s.ExpectResponse(CodeOk, strings.ToUpper(HELP))
			if test.alreadyTls {
				s.Session.Conn = &tls.Conn{}
			}

			target.HandleCommand(context.TODO(), s.Session, []string{"test"})
			assert.Equal(t, "test", s.Session.SenderDomain)
			assert.Nil(t, s.Session.EnvelopeFrom)
			assert.Empty(t, s.Session.EnvelopeTo)
			assert.Empty(t, s.Session.RawData)
		})
	}
}
