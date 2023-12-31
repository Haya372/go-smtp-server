package session

import (
	"bufio"
	"net"
	"net/mail"
	"net/textproto"

	"github.com/Haya372/hlog"
	"github.com/google/uuid"
)

type SessionFactory interface {
	CreateSession(conn net.Conn) *Session
}

type SessionFactoryImpl struct {
	log hlog.Logger
}

func (f *SessionFactoryImpl) CreateSession(conn net.Conn) *Session {
	return &Session{
		Id:         uuid.New(),
		EnvelopeTo: make([]mail.Address, 0),
		Conn:       conn,
		log:        f.log,
		reader:     *textproto.NewReader(bufio.NewReader(conn)),
		writer:     *textproto.NewWriter(bufio.NewWriter(conn)),
	}
}

func NewSessionFactory(log hlog.Logger) SessionFactory {
	return &SessionFactoryImpl{
		log: log,
	}
}
