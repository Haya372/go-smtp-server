package mock

import (
	mail "net/mail"

	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func NewAddressMatcher(address mail.Address) gomock.Matcher {
	return &addressMatcher{address: address}
}

type addressMatcher struct {
	address mail.Address
}

func (m *addressMatcher) Matches(x interface{}) bool {
	// ptr
	pArg, pOk := x.(*mail.Address)
	if pOk {
		return m.address.String() == pArg.String()
	}
	// struct
	sArg, sOk := x.(mail.Address)
	if sOk {
		return m.address.String() == sArg.String()
	}

	return false
}

func (m *addressMatcher) String() string {
	return "Matcher:" + m.address.String()
}

func NewInitializedMockLogger(ctrl *gomock.Controller) *MockLogger {
	log := NewMockLogger(ctrl)

	log.EXPECT().WithError(gomock.Any()).Return(log).AnyTimes()
	log.EXPECT().Trace(gomock.Any(), gomock.Any()).AnyTimes()
	log.EXPECT().Tracef(gomock.Any(), gomock.Any()).AnyTimes()
	log.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	log.EXPECT().Debugf(gomock.Any(), gomock.Any()).AnyTimes()
	log.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	log.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
	log.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
	log.EXPECT().Warnf(gomock.Any(), gomock.Any()).AnyTimes()
	log.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	log.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
	log.EXPECT().Fatal(gomock.Any(), gomock.Any()).AnyTimes()
	log.EXPECT().Fatalf(gomock.Any(), gomock.Any()).AnyTimes()

	return log
}

type SessionMockParam struct {
	IsCloseImmediately bool
	SenderDomain       string
	EnvelopeFrom       string
	EnvelopeTo         []string
	RawData            []byte
}

func NewInitializedMockSession(ctrl *gomock.Controller, p SessionMockParam) *MockSession {
	s := NewMockSession(ctrl)

	s.EXPECT().Id().Return(uuid.New()).AnyTimes()
	s.EXPECT().SenderDomain().Return(p.SenderDomain).AnyTimes()
	if len(p.EnvelopeFrom) == 0 {
		s.EXPECT().EnvelopeFrom().Return(nil).AnyTimes()
	} else {
		address, _ := mail.ParseAddress(p.EnvelopeFrom)
		s.EXPECT().EnvelopeFrom().Return(address).AnyTimes()
	}
	envelopeTo := make([]mail.Address, 0)
	for _, to := range p.EnvelopeTo {
		address, _ := mail.ParseAddress(to)
		envelopeTo = append(envelopeTo, *address)
	}
	s.EXPECT().EnvelopeTo().Return(envelopeTo).AnyTimes()
	s.EXPECT().RawData().Return(p.RawData).AnyTimes()

	return s
}

func NewInitializedMockCommandHandler(ctrl *gomock.Controller, command string) *MockCommandHandler {
	h := NewMockCommandHandler(ctrl)

	h.EXPECT().Command().Return(command).AnyTimes()

	return h
}
