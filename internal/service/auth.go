package service

import (
	"context"
	"net"
	"net/mail"
	"strings"

	"blitiri.com.ar/go/spf"
	"github.com/Haya372/hlog"
	"github.com/Haya372/smtp-server/internal/data"
)

type AuthService interface {
	Auth(ctx context.Context, mime data.MimeData) (data.AuthResult, error)
}

type authServiceImpl struct {
	log hlog.Logger
}

func (s *authServiceImpl) Auth(ctx context.Context, mime data.MimeData) (data.AuthResult, error) {
	spf, err := s.spf(ctx, mime.Ip, getDomain(*mime.EnvelopeFrom), mime.SenderDomain)
	if err != nil {
		s.log.WithError(err).Errorf("spf error")
	}
	return data.AuthResult{Spf: spf}, nil
}

func (s *authServiceImpl) spf(ctx context.Context, ip net.IP, helo string, sender string) (spf.Result, error) {
	return spf.CheckHostWithSender(ip, helo, sender, spf.WithContext(ctx))
}

func NewAuthService(log hlog.Logger) AuthService {
	return &authServiceImpl{
		log: log,
	}
}

func getDomain(address mail.Address) string {
	list := strings.Split(address.Address, "@")
	return list[1]
}
