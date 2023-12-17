package command

import (
	"context"
	"crypto/tls"
	"errors"
	"testing"

	"github.com/Haya372/smtp-server/internal/config"
	"github.com/Haya372/smtp-server/internal/mock"
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
		setup     func(s *mock.MockSession)
		expectErr bool
	}{
		{
			name: "Already TLS",
			setup: func(s *mock.MockSession) {
				s.EXPECT().IsTls().Return(true).Times(1)
				s.EXPECT().Response(gomock.Eq(CodeBadSequence), gomock.Eq(MsgAlreadyTls))
			},
		},
		{
			name: "TLS Error",
			setup: func(s *mock.MockSession) {
				s.EXPECT().IsTls().Return(false).Times(1)
				s.EXPECT().Response(gomock.Eq(CodeGreet), gomock.Eq(MsgGoAhead))
				s.EXPECT().ConvertToTls(gomock.Eq(conf)).Return(errors.New("error")).Times(1)
				s.EXPECT().Response(gomock.Eq(CodeTransactionFail), gomock.Eq(MsgTransactionFail)).Times(1)
			},
			expectErr: true,
		},
		{
			name: "Success",
			setup: func(s *mock.MockSession) {
				s.EXPECT().IsTls().Return(false).Times(1)
				s.EXPECT().Response(gomock.Eq(CodeGreet), gomock.Eq(MsgGoAhead))
				s.EXPECT().ConvertToTls(gomock.Any()).Return(nil).Times(1)
				s.EXPECT().Reset().Times(1)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := mock.NewInitializedMockSession(ctrl, mock.SessionMockParam{})
			test.setup(s)

			target := NewStartTlsHandler(log, tlsConf)
			err := target.HandleCommand(context.TODO(), s, []string{})

			if test.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
