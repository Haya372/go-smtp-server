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

type Session interface {
	// return unique session ID
	Id() uuid.UUID
	// get client IP address
	IP() net.IP
	// close connection after this method called
	CloseImmediately()
	IsCloseImmediately() bool
	// set envelope from address
	SetEnvelopeFrom(address *mail.Address)
	// get envelope from address
	EnvelopeFrom() *mail.Address
	// set client hostname
	SetSenderDomain(domain string)
	// get client hostname
	SenderDomain() string
	// add envelope to address
	AddEnvelopeTo(address mail.Address)
	// get envelope to address list
	EnvelopeTo() []mail.Address
	// set raw mime data
	SetRawData(rawData []byte)
	// get raw mime data
	RawData() []byte
	// read message from client
	ReadLine() (string, error)
	// read mime data
	ReadRawData() ([]byte, error)
	// close connection
	Close()
	// response with status code
	Response(code int, msg string)
	// response single line message
	ResponseLine(line string)
	// reset session state
	Reset()
	// check TLS conn
	IsTls() bool
	// update to TLS
	ConvertToTls(tlsConf *tls.Config) error
}

type sessionImpl struct {
	// connection unique ID
	id uuid.UUID
	// when this flag is true, connection will close immediately
	shouldClose bool
	// domain name received by HELO/EHLO
	senderDomain string
	// sender address received by MAIL
	envelopeFrom *mail.Address
	// recipient addresses received by RCPT
	envelopeTo []mail.Address
	// raw mail data
	rawData []byte

	conn   net.Conn
	log    hlog.Logger
	reader textproto.Reader
	writer textproto.Writer
}

func (s *sessionImpl) Id() uuid.UUID {
	return s.id
}

func (s *sessionImpl) IP() net.IP {
	return net.IP(s.conn.RemoteAddr().Network())
}

func (s *sessionImpl) CloseImmediately() {
	s.shouldClose = true
}

func (s *sessionImpl) IsCloseImmediately() bool {
	return s.shouldClose
}

func (s *sessionImpl) SetEnvelopeFrom(address *mail.Address) {
	s.envelopeFrom = address
}

func (s *sessionImpl) EnvelopeFrom() *mail.Address {
	return s.envelopeFrom
}

func (s *sessionImpl) SetSenderDomain(domain string) {
	s.senderDomain = domain
}

func (s *sessionImpl) SenderDomain() string {
	return s.senderDomain
}

func (s *sessionImpl) AddEnvelopeTo(address mail.Address) {
	s.envelopeTo = append(s.envelopeTo, address)
}

func (s *sessionImpl) EnvelopeTo() []mail.Address {
	return s.envelopeTo
}

func (s *sessionImpl) SetRawData(rawData []byte) {
	s.rawData = rawData
}

func (s *sessionImpl) RawData() []byte {
	return s.rawData
}

func (s *sessionImpl) ReadLine() (string, error) {
	return s.reader.ReadLine()
}

func (s *sessionImpl) ReadRawData() ([]byte, error) {
	return s.reader.ReadDotBytes()
}

func (s *sessionImpl) Close() {
	s.conn.Close()
}

func (s *sessionImpl) Response(code int, msg string) {
	s.writer.PrintfLine("%d %s", code, msg)
}

func (s *sessionImpl) ResponseLine(line string) {
	s.writer.PrintfLine(line)
}

func (s *sessionImpl) Reset() {
	s.senderDomain = ""
	s.shouldClose = false
	s.envelopeFrom = nil
	s.envelopeTo = make([]mail.Address, 0)
	s.rawData = make([]byte, 0)
}

func (s *sessionImpl) IsTls() bool {
	switch s.conn.(type) {
	case *tls.Conn:
		return true
	default:
		return false
	}
}

func (s *sessionImpl) ConvertToTls(tlsConf *tls.Config) error {
	conn := tls.Server(s.conn, tlsConf)
	if err := conn.Handshake(); err != nil {
		return err
	}
	s.conn = conn
	s.reader = *textproto.NewReader(bufio.NewReader(conn))
	s.writer = *textproto.NewWriter(bufio.NewWriter(conn))
	return nil
}
