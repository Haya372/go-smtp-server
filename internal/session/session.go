package session

import (
	"bufio"
	"crypto/tls"
	"net"
	"net/mail"
	"net/textproto"

	"github.com/Haya372/hlog"
	"github.com/google/uuid"
)

type Session struct {
	// connection unique ID
	Id uuid.UUID
	// when this flag is true, connection will close immediately
	ShouldClose bool
	// domain name received by HELO/EHLO
	SenderDomain string
	// sender address received by MAIL
	EnvelopeFrom *mail.Address
	// recipient addresses received by RCPT
	EnvelopeTo []mail.Address
	// raw mail data
	RawData []byte

	Conn   net.Conn
	log    hlog.Logger
	reader textproto.Reader
	writer textproto.Writer
}

func (s *Session) IP() net.IP {
	return net.IP(s.Conn.RemoteAddr().Network())
}

func (s *Session) AddEnvelopeTo(address mail.Address) {
	s.EnvelopeTo = append(s.EnvelopeTo, address)
}

func (s *Session) ReadLine() (string, error) {
	return s.reader.ReadLine()
}

func (s *Session) ReadRawData() ([]byte, error) {
	return s.reader.ReadDotBytes()
}

func (s *Session) Close() {
	s.Conn.Close()
}

func (s *Session) Response(code int, msg string) error {
	return s.writer.PrintfLine("%d %s", code, msg)
}

func (s *Session) ResponseLine(line string) error {
	return s.writer.PrintfLine(line)
}

func (s *Session) Reset() {
	s.SenderDomain = ""
	s.ShouldClose = false
	s.EnvelopeFrom = nil
	s.EnvelopeTo = make([]mail.Address, 0)
	s.RawData = make([]byte, 0)
}

func (s *Session) IsTls() bool {
	switch s.Conn.(type) {
	case *tls.Conn:
		return true
	default:
		return false
	}
}

func (s *Session) ConvertToTls(tlsConf *tls.Config) error {
	conn := tls.Server(s.Conn, tlsConf)
	if err := conn.Handshake(); err != nil {
		return err
	}
	s.Conn = conn
	s.reader = *textproto.NewReader(bufio.NewReader(conn))
	s.writer = *textproto.NewWriter(bufio.NewWriter(conn))
	return nil
}
