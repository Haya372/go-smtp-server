package server

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/Haya372/hlog"
	"github.com/Haya372/smtp-server/internal/command"
	"github.com/Haya372/smtp-server/internal/connection"
	"github.com/Haya372/smtp-server/internal/session"
	"golang.org/x/sync/semaphore"
)

type Server struct {
	Port              string
	ConnectionTimeout time.Duration

	s       *semaphore.Weighted
	ln      net.Listener
	log     hlog.Logger
	wg      sync.WaitGroup
	factory session.SessionFactory
	handler connection.SessionHandler
}

func (s *Server) ListenSmtp() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", s.Port)
	if err != nil {
		return err
	}

	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	s.ln = ln
	s.waitConnection(context.Background())

	return nil
}

func (s *Server) waitConnection(parentCtx context.Context) {
	if s.ln == nil {
		s.log.Fatal("TCPLister is not defined.", nil)
	}
	defer s.ln.Close()

	for {
		conn, err := s.ln.Accept()
		if err != nil {
			s.log.WithError(err).Error("could not accept session.", nil)
			continue
		}

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			smtpSession := s.factory.CreateSession(conn)
			// TODO: セマフォを取得できなければ抜ける
			ctx, cancel := context.WithTimeout(parentCtx, 10*time.Millisecond)
			defer cancel()

			err = s.s.Acquire(ctx, 1)
			if err != nil {
				s.log.WithError(err).Error("could not get semaphore.", nil)
				smtpSession.Response(command.CodeTransactionFail, command.MsgBadSequence)
				conn.Close()
				return
			}
			defer s.s.Release(1)

			s.handler.HandleSession(ctx, smtpSession)
		}()
	}
}

// TODO: 設定値を与える
func NewServer(log hlog.Logger, connFactory session.SessionFactory, handler connection.SessionHandler) Server {
	return Server{
		Port:              ":25",
		ConnectionTimeout: 30 * time.Second,
		log:               log,
		factory:           connFactory,
		s:                 semaphore.NewWeighted(1),
		handler:           handler,
	}
}
