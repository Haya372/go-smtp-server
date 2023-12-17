package command

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Haya372/smtp-server/internal/config"
	"github.com/Haya372/smtp-server/internal/mock"
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
			s := mock.NewInitializedMockSession(ctrl, mock.SessionMockParam{})

			s.EXPECT().Response(gomock.Eq(test.code), gomock.Eq(test.msg)).Times(1)

			target := NewEhloHandler(log, conf)
			target.HandleCommand(context.TODO(), s, test.arg)
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
		setup      func(s *mock.MockSession)
		alreadyTls bool
	}{
		{
			name: "no extension",
			conf: &config.SmtpConfig{},
			setup: func(s *mock.MockSession) {
				s.EXPECT().ResponseLine(gomock.Any()).Times(1)
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
			setup: func(s *mock.MockSession) {
				s.EXPECT().ResponseLine(gomock.Any()).Times(1)
				s.EXPECT().ResponseLine(gomock.Eq(fmt.Sprintf("%d-PIPELINING", CodeOk))).Times(1)
				s.EXPECT().ResponseLine(gomock.Eq(fmt.Sprintf("%d-8BITMIME", CodeOk))).Times(1)
				s.EXPECT().ResponseLine(gomock.Eq(fmt.Sprintf("%d-SIZE %d", CodeOk, 1))).Times(1)
				s.EXPECT().ResponseLine(gomock.Eq(fmt.Sprintf("%d-STARTTLS", CodeOk))).Times(1)
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
			setup: func(s *mock.MockSession) {
				s.EXPECT().ResponseLine(gomock.Any()).Times(1)
				s.EXPECT().ResponseLine(gomock.Eq(fmt.Sprintf("%d-PIPELINING", CodeOk))).Times(1)
				s.EXPECT().ResponseLine(gomock.Eq(fmt.Sprintf("%d-8BITMIME", CodeOk))).Times(1)
				s.EXPECT().ResponseLine(gomock.Eq(fmt.Sprintf("%d-SIZE %d", CodeOk, 1))).Times(1)
			},
			alreadyTls: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			target := NewEhloHandler(log, test.conf)

			s := mock.NewMockSession(ctrl)
			s.EXPECT().Reset().Times(1)
			s.EXPECT().SetSenderDomain("test")
			test.setup(s)
			s.EXPECT().Response(gomock.Eq(CodeOk), gomock.Eq(strings.ToUpper(HELP))).Times(1)
			s.EXPECT().IsTls().AnyTimes().Return(test.alreadyTls)

			target.HandleCommand(context.TODO(), s, []string{"test"})
		})
	}
}
