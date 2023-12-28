package command

import (
	"context"
	"crypto/tls"
	"testing"

	"github.com/Haya372/smtp-server/internal/config"
	"github.com/Haya372/smtp-server/internal/mock"
	"github.com/Haya372/smtp-server/internal/session"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestStartTls_Command(t *testing.T) {
	target := NewStartTlsHandler(nil, nil)
	assert.Equal(t, STARTTLS, target.Command())
}

func TestStartTls(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock.NewInitializedMockLogger(ctrl)

	conf := &tls.Config{}

	tlsConf := &config.TlsConfig{
		TlsConfig: conf,
	}

	tests := []struct {
		name      string
		setup     func(s *session.MockSession)
		expectErr bool
	}{
		{
			name: "Already TLS",
			setup: func(s *session.MockSession) {
				s.Session.Conn = &tls.Conn{}
				s.ExpectResponse(CodeBadSequence, MsgAlreadyTls)
			},
		},
		// NOTE: モックだとテストが難しいため後回し
		// TODO: テストコード修正
		// {
		// 	name: "TLS Error",
		// 	setup: func(s *session.MockSession) {
		// 		s.ExpectResponse(CodeGreet, MsgGoAhead)
		// 		cer, _ := tls.LoadX509KeyPair("./testdata/server.crt", "server.key")
		// 		tlsConf.TlsConfig.Certificates = []tls.Certificate{cer}
		// 		s.ExpectResponse(CodeTransactionFail, MsgTransactionFail)
		// 	},
		// 	expectErr: true,
		// },
		// {
		// 	name: "Success",
		// 	setup: func(s *session.MockSession) {
		// 		s.ExpectResponse(CodeGreet, MsgGoAhead)
		// 		cer, err := tls.LoadX509KeyPair("../../testdata/server.crt", "../../testdata/server.key")
		// 		dir, _ := os.Getwd()
		// 		t.Log(dir)
		// 		if err != nil {
		// 			t.Log(err)
		// 			t.Fail()
		// 		}
		// 		tlsConf.TlsConfig.Certificates = []tls.Certificate{cer}
		// 	},
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := session.NewMockSession(ctrl)
			test.setup(s)

			target := NewStartTlsHandler(log, tlsConf)
			err := target.HandleCommand(context.TODO(), s.Session, []string{})

			if test.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
