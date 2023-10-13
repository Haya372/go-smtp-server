package service

import (
	"bytes"
	"context"
	"errors"
	"net"
	"net/mail"
	"strings"

	"blitiri.com.ar/go/spf"
	"github.com/Haya372/hlog"
	"github.com/Haya372/smtp-server/internal/data"
	"github.com/emersion/go-msgauth/authres"
	"github.com/emersion/go-msgauth/dkim"
	"github.com/emersion/go-msgauth/dmarc"
)

type AuthService interface {
	Auth(ctx context.Context, mime data.MimeData) *data.AuthResult
}

type authServiceImpl struct {
	log hlog.Logger
}

func (s *authServiceImpl) Auth(ctx context.Context, mime data.MimeData) *data.AuthResult {
	res := &data.AuthResult{
		Spf:  s.spf(ctx, mime.Ip, mime.EnvelopeFrom.Address, mime.SenderDomain),
		Dkim: s.dkim(ctx, mime),
	}
	// TODO: Envelope-FromからHeader-Fromに直す
	res.Dmarc = s.dmarc(ctx, res, *mime.EnvelopeFrom)
	return res
}

func (s *authServiceImpl) spf(ctx context.Context, ip net.IP, helo string, sender string) authres.SPFResult {
	res, err := spf.CheckHostWithSender(ip, helo, sender, spf.WithContext(ctx))
	if err != nil {
		s.log.WithError(err).Errorf("spf error")
	}
	spfRes := authres.SPFResult{
		From: sender,
		Helo: helo,
	}
	switch res {
	case spf.None:
		spfRes.Value = authres.ResultNone
	case spf.Pass:
		spfRes.Value = authres.ResultPass
	case spf.Fail:
		spfRes.Value = authres.ResultFail
	case spf.SoftFail:
		spfRes.Value = authres.ResultSoftFail
	case spf.TempError:
		spfRes.Value = authres.ResultTempError
	case spf.PermError:
		spfRes.Value = authres.ResultPermError
	default:
		spfRes.Value = authres.ResultNeutral
	}
	return spfRes
}

func verifyDmarc(dmarcDomain, authDomain string, result authres.ResultValue, mode dmarc.AlignmentMode) (bool, error) {
	switch mode {
	case dmarc.AlignmentStrict:
		if dmarcDomain == authDomain {
			if result == authres.ResultPass {
				return true, nil
			} else if result == authres.ResultTempError {
				return false, errors.New("auth temp failed")
			}
		}
	case dmarc.AlignmentRelaxed:
		if isSubDomain(dmarcDomain, authDomain) || isSubDomain(authDomain, dmarcDomain) {
			if result == authres.ResultPass {
				return true, nil
			} else if result == authres.ResultTempError {
				return false, errors.New("auth temp failed")
			}
		}
	}
	return false, nil
}

func (s *authServiceImpl) dmarc(ctx context.Context, authRes *data.AuthResult, from mail.Address) authres.DMARCResult {
	dmarcDomain := getDomain(from)
	record, err := dmarc.Lookup(dmarcDomain)
	if err != nil {
		if dmarc.IsTempFail(err) {
			return authres.DMARCResult{
				Value: authres.ResultTempError,
				From:  from.Address,
			}
		}
		return authres.DMARCResult{
			Value: authres.ResultNone,
			From:  from.Address,
		}
	}

	// check based on spf
	spfDomain := getDomain(mail.Address{Address: authRes.Spf.From})
	ok, spfErr := verifyDmarc(dmarcDomain, spfDomain, authRes.Spf.Value, record.SPFAlignment)
	if ok {
		return authres.DMARCResult{
			Value: authres.ResultPass,
			From:  from.Address,
		}
	}

	// dkim check
	hasDkimErr := false
	for _, dkimRes := range authRes.Dkim {
		ok, dkimErr := verifyDmarc(dmarcDomain, dkimRes.Domain, dkimRes.Value, record.DKIMAlignment)
		if ok {
			return authres.DMARCResult{
				Value: authres.ResultPass,
				From:  from.Address,
			}
		}
		if dkimErr != nil {
			hasDkimErr = true
		}
	}

	// authentication error
	if (spfErr != nil) || hasDkimErr {
		return authres.DMARCResult{
			Value: authres.ResultTempError,
			From:  from.Address,
		}
	}

	return authres.DMARCResult{
		Value: authres.ResultFail,
		From:  from.Address,
	}
}

func (s *authServiceImpl) dkim(ctx context.Context, mime data.MimeData) []authres.DKIMResult {
	reader := bytes.NewReader(mime.RawData)
	// TODO: 設定から指定できるようにする
	dkims, err := dkim.VerifyWithOptions(reader, &dkim.VerifyOptions{MaxVerifications: 3})
	if err != nil {
		s.log.WithError(err).Errorf("dkim error")
	}

	// no dkim-signature header
	if len(dkims) == 0 {
		return []authres.DKIMResult{
			{
				Value: authres.ResultNone,
			},
		}
	}

	res := make([]authres.DKIMResult, len(dkims))

	for idx, d := range dkims {
		dkimRes := authres.DKIMResult{
			Domain:     d.Domain,
			Identifier: d.Identifier,
		}

		if d.Err == nil {
			dkimRes.Value = authres.ResultPass
		} else if dkim.IsTempFail(d.Err) {
			dkimRes.Value = authres.ResultTempError
		} else if dkim.IsPermFail(d.Err) {
			dkimRes.Value = authres.ResultPermError
		} else {
			dkimRes.Value = authres.ResultFail
		}
		res[idx] = dkimRes
	}

	return res
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

func isSubDomain(child, parent string) bool {
	if strings.Compare(parent, child) == 0 {
		return true
	}

	return strings.HasSuffix(child, "."+parent)
}
