package data

import (
	"net"
	"net/mail"

	"github.com/Haya372/smtp-server/internal/session"
)

type MimeData struct {
	// client IP address
	Ip net.IP
	// envelope from address
	EnvelopeFrom *mail.Address
	// envelope to address
	EnvelopeTo []mail.Address
	// ehlo domain
	SenderDomain string
	// raw mime data
	RawData []byte
	// authentication result
	AuthResult AuthResult

	// private field
	// headers which this system append
	header map[string]string
}

func (m *MimeData) AddHeader(key, val string) {
	m.header[key] = val
}

func NewMimeData(session session.Session) *MimeData {
	return &MimeData{
		EnvelopeFrom: session.EnvelopeFrom,
		EnvelopeTo:   session.EnvelopeTo,
		SenderDomain: session.SenderDomain,
		RawData:      session.RawData,
	}
}
