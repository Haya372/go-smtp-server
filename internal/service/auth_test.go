package service

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net"
	"net/mail"
	"strings"
	"testing"

	"github.com/Haya372/smtp-server/internal/data"
	"github.com/Haya372/smtp-server/internal/mock"
	"github.com/emersion/go-msgauth/authres"
	"github.com/emersion/go-msgauth/dkim"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mockResolver struct {
	txt    []string
	mx     []*net.MX
	ipAddr []net.IPAddr
	addr   []string
	err    error
}

func (r *mockResolver) LookupTXT(ctx context.Context, name string) ([]string, error) {
	return r.txt, r.err
}

func (r *mockResolver) LookupMX(ctx context.Context, name string) ([]*net.MX, error) {
	return r.mx, r.err
}

func (r *mockResolver) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	return r.ipAddr, r.err
}

func (r *mockResolver) LookupAddr(ctx context.Context, addr string) (names []string, err error) {
	return r.addr, r.err
}

func TestSpf(t *testing.T) {
	tests := []struct {
		name      string
		txtRecord []string
		err       error
		expect    authres.ResultValue
	}{
		{
			name:      "none",
			txtRecord: []string{},
			expect:    authres.ResultNone,
		},
		{
			name: "pass",
			txtRecord: []string{
				"v=spf1 +ip4:1.2.3.4 -all",
			},
			expect: authres.ResultPass,
		},
		{
			name: "fail",
			txtRecord: []string{
				"v=spf1 +ip4:1.2.3.5 -all",
			},
			expect: authres.ResultFail,
		},
		{
			name: "softfail",
			txtRecord: []string{
				"v=spf1 +ip4:1.2.3.5 ~all",
			},
			expect: authres.ResultSoftFail,
		},
		{
			name:   "temperror",
			err:    &net.DNSError{IsTemporary: true},
			expect: authres.ResultTempError,
		},
		{
			name: "permerror",
			txtRecord: []string{
				"+ip4:1.2.3. ~all",
			},
			err:    errors.New("permerror"),
			expect: authres.ResultPermError,
		},
		{
			name: "neutral",
			txtRecord: []string{
				"v=spf1 ?ip4:1.2.3.4 ~all",
			},
			expect: authres.ResultNeutral,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock.NewInitializedMockLogger(ctrl)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resolver := &mockResolver{
				txt: test.txtRecord,
				err: test.err,
			}

			ip := net.IPv4(1, 2, 3, 4)
			helo := "example.com"
			sender := "example.net"

			target := authServiceImpl{
				log:      log,
				resolver: resolver,
			}

			res := target.spf(context.TODO(), ip, helo, sender)

			assert.Equal(t, test.expect, res.Value)
			assert.Equal(t, helo, res.Helo)
			assert.Equal(t, sender, res.From)
		})
	}
}

func getPemPublicKey(key crypto.PublicKey) (string, error) {
	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("der fail")
	}

	der := x509.MarshalPKCS1PublicKey(rsaKey)

	var b bytes.Buffer
	if err := pem.Encode(&b, &pem.Block{Bytes: der}); err != nil {
		return "", err
	}

	lines := strings.Split(b.String(), "\n")

	return strings.Join(lines[1:len(lines)-2], "\n"), nil
}

func TestDkim(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		assert.Fail(t, "fail to generate private key")
	}

	rawData := `
	From: test-from
	To: test-to
	Subject: test

	test mail
	`

	var b bytes.Buffer
	opt := &dkim.SignOptions{
		Domain:   "example.com",
		Selector: "test",
		HeaderKeys: []string{
			"from", "to", "subject",
		},
		Signer: key,
	}

	dkim.Sign(&b, strings.NewReader(rawData), opt)

	pem, err := getPemPublicKey(key.Public())
	if err != nil {
		assert.Fail(t, "public key fail")
	}
	dkimRecord := "v=DKIM1; p=" + pem

	tests := []struct {
		name   string
		data   []byte
		txt    []string
		err    error
		result []authres.DKIMResult
	}{
		{
			name: "none",
			data: []byte(rawData),
			result: []authres.DKIMResult{
				{
					Value: authres.ResultNone,
				},
			},
		},
		{
			name: "pass",
			data: b.Bytes(),
			txt:  []string{dkimRecord},
			result: []authres.DKIMResult{
				{
					Value:      authres.ResultPass,
					Domain:     "example.com",
					Identifier: "test",
				},
			},
		},
		{
			name: "temperror",
			data: b.Bytes(),
			err:  &net.DNSError{IsTemporary: true, IsTimeout: true},
			result: []authres.DKIMResult{
				{
					Value:      authres.ResultTempError,
					Domain:     "example.com",
					Identifier: "test",
				},
			},
		},
		{
			name: "permerror",
			data: b.Bytes(),
			err:  errors.New("permerror"),
			result: []authres.DKIMResult{
				{
					Value:      authres.ResultPermError,
					Domain:     "example.com",
					Identifier: "test",
				},
			},
		},
		{
			name: "fail",
			data: b.Bytes()[0 : b.Len()-2],
			txt:  []string{dkimRecord},
			result: []authres.DKIMResult{
				{
					Value:      authres.ResultFail,
					Domain:     "example.com",
					Identifier: "test",
				},
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock.NewInitializedMockLogger(ctrl)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mime := data.MimeData{
				RawData: test.data,
			}

			resolver := &mockResolver{
				txt: test.txt,
				err: test.err,
			}

			target := authServiceImpl{
				log:      log,
				resolver: resolver,
			}

			result := target.dkim(context.TODO(), mime)

			assert.Equal(t, len(test.result), len(result))
			for idx := range result {
				assert.Equal(t, test.result[idx].Value, result[idx].Value)
			}
		})
	}
}

func TestDmarc(t *testing.T) {
	tests := []struct {
		name    string
		authRes data.AuthResult
		txt     []string
		err     error
		expect  authres.DMARCResult
	}{
		{
			name: "pass(spf: pass, dkim: none, strict)",
			authRes: data.AuthResult{
				Spf: authres.SPFResult{
					From:  "test@example.com",
					Value: authres.ResultPass,
				},
				Dkim: []authres.DKIMResult{},
			},
			txt: []string{"v=DMARC1; p=quarantine; adkim=r; aspf=s; rua=mailto:rua@example.com"},
			expect: authres.DMARCResult{
				Value: authres.ResultPass,
			},
		},
		{
			name: "temperror(spf: temperror, dkim: none, strict)",
			authRes: data.AuthResult{
				Spf: authres.SPFResult{
					From:  "test@example.com",
					Value: authres.ResultTempError,
				},
				Dkim: []authres.DKIMResult{},
			},
			txt: []string{"v=DMARC1; p=quarantine; adkim=r; aspf=s; rua=mailto:rua@example.com"},
			expect: authres.DMARCResult{
				Value: authres.ResultTempError,
			},
		},
		{
			name: "fail(spf: pass but other domain, dkim: none, strict)",
			authRes: data.AuthResult{
				Spf: authres.SPFResult{
					From:  "test@child.example.com",
					Value: authres.ResultPass,
				},
				Dkim: []authres.DKIMResult{},
			},
			txt: []string{"v=DMARC1; p=quarantine; adkim=r; aspf=s; rua=mailto:rua@example.com"},
			expect: authres.DMARCResult{
				Value: authres.ResultFail,
			},
		},
		{
			name: "pass(spf: pass, dkim: none, relax)",
			authRes: data.AuthResult{
				Spf: authres.SPFResult{
					From:  "test@child.example.com",
					Value: authres.ResultPass,
				},
				Dkim: []authres.DKIMResult{},
			},
			txt: []string{"v=DMARC1; p=quarantine; adkim=r; aspf=r; rua=mailto:rua@example.com"},
			expect: authres.DMARCResult{
				Value: authres.ResultPass,
			},
		},
		{
			name: "temperror(spf: temperror, dkim: none, relax)",
			authRes: data.AuthResult{
				Spf: authres.SPFResult{
					From:  "test@example.com",
					Value: authres.ResultTempError,
				},
				Dkim: []authres.DKIMResult{},
			},
			txt: []string{"v=DMARC1; p=quarantine; adkim=r; aspf=r; rua=mailto:rua@example.com"},
			expect: authres.DMARCResult{
				Value: authres.ResultTempError,
			},
		},
		{
			name: "fail(spf: fail, dkim: none, relax)",
			authRes: data.AuthResult{
				Spf: authres.SPFResult{
					From:  "test@example.com",
					Value: authres.ResultFail,
				},
				Dkim: []authres.DKIMResult{},
			},
			txt: []string{"v=DMARC1; p=quarantine; adkim=r; aspf=r; rua=mailto:rua@example.com"},
			expect: authres.DMARCResult{
				Value: authres.ResultFail,
			},
		},
		{
			name: "pass(spf: none, dkim: pass, strict)",
			authRes: data.AuthResult{
				Spf: authres.SPFResult{
					From:  "test@example.com",
					Value: authres.ResultNone,
				},
				Dkim: []authres.DKIMResult{
					{
						Domain: "example.com",
						Value:  authres.ResultPass,
					},
				},
			},
			txt: []string{"v=DMARC1; p=quarantine; adkim=s; aspf=s; rua=mailto:rua@example.com"},
			expect: authres.DMARCResult{
				Value: authres.ResultPass,
			},
		},
		{
			name: "temperror(spf: none, dkim: temperror, strict)",
			authRes: data.AuthResult{
				Spf: authres.SPFResult{
					From:  "test@example.com",
					Value: authres.ResultNone,
				},
				Dkim: []authres.DKIMResult{
					{
						Domain: "example.com",
						Value:  authres.ResultTempError,
					},
				},
			},
			txt: []string{"v=DMARC1; p=quarantine; adkim=s; aspf=s; rua=mailto:rua@example.com"},
			expect: authres.DMARCResult{
				Value: authres.ResultTempError,
			},
		},
		{
			name: "fail(spf: none, dkim: pass but other domain, strict)",
			authRes: data.AuthResult{
				Spf: authres.SPFResult{
					From:  "test@child.example.com",
					Value: authres.ResultNone,
				},
				Dkim: []authres.DKIMResult{
					{
						Domain: "child.example.com",
						Value:  authres.ResultPass,
					},
				},
			},
			txt: []string{"v=DMARC1; p=quarantine; adkim=s; aspf=s; rua=mailto:rua@example.com"},
			expect: authres.DMARCResult{
				Value: authres.ResultFail,
			},
		},
		{
			name: "pass(spf: none, dkim: pass, relax)",
			authRes: data.AuthResult{
				Spf: authres.SPFResult{
					From:  "test@example.com",
					Value: authres.ResultNone,
				},
				Dkim: []authres.DKIMResult{
					{
						Domain: "child.example.com",
						Value:  authres.ResultPass,
					},
				},
			},
			txt: []string{"v=DMARC1; p=quarantine; adkim=r; aspf=r; rua=mailto:rua@example.com"},
			expect: authres.DMARCResult{
				Value: authres.ResultPass,
			},
		},
		{
			name: "temperror(spf: none, dkim: temperror, relax)",
			authRes: data.AuthResult{
				Spf: authres.SPFResult{
					From:  "test@example.com",
					Value: authres.ResultNone,
				},
				Dkim: []authres.DKIMResult{
					{
						Domain: "child.example.com",
						Value:  authres.ResultTempError,
					},
				},
			},
			txt: []string{"v=DMARC1; p=quarantine; adkim=r; aspf=r; rua=mailto:rua@example.com"},
			expect: authres.DMARCResult{
				Value: authres.ResultTempError,
			},
		},
		{
			name: "fail(spf: none, dkim: fail, relax)",
			authRes: data.AuthResult{
				Spf: authres.SPFResult{
					From:  "test@example.com",
					Value: authres.ResultFail,
				},
				Dkim: []authres.DKIMResult{
					{
						Domain: "child.example.com",
						Value:  authres.ResultFail,
					},
				},
			},
			txt: []string{"v=DMARC1; p=quarantine; adkim=r; aspf=r; rua=mailto:rua@example.com"},
			expect: authres.DMARCResult{
				Value: authres.ResultFail,
			},
		},
		{
			name: "tempfail(dns temp error)",
			authRes: data.AuthResult{
				Spf: authres.SPFResult{
					From:  "test@example.com",
					Value: authres.ResultPass,
				},
				Dkim: []authres.DKIMResult{
					{
						Domain: "child.example.com",
						Value:  authres.ResultPass,
					},
				},
			},
			err: &net.DNSError{IsTemporary: true},
			expect: authres.DMARCResult{
				Value: authres.ResultTempError,
			},
		},
		{
			name: "none(dns perm error)",
			authRes: data.AuthResult{
				Spf: authres.SPFResult{
					From:  "test@example.com",
					Value: authres.ResultPass,
				},
				Dkim: []authres.DKIMResult{
					{
						Domain: "child.example.com",
						Value:  authres.ResultPass,
					},
				},
			},
			err: errors.New("permerror"),
			expect: authres.DMARCResult{
				Value: authres.ResultNone,
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock.NewInitializedMockLogger(ctrl)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resolver := &mockResolver{
				txt: test.txt,
				err: test.err,
			}

			target := authServiceImpl{
				log:      log,
				resolver: resolver,
			}

			result := target.dmarc(context.TODO(), &test.authRes, mail.Address{Address: "test@example.com"})

			assert.Equal(t, test.expect.Value, result.Value)
			assert.Equal(t, "test@example.com", result.From)
		})
	}
}

func TestAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := mock.NewInitializedMockLogger(ctrl)
	resolver := mockResolver{}

	target := authServiceImpl{
		log:      log,
		resolver: &resolver,
	}

	res := target.Auth(context.TODO(), data.MimeData{
		Ip:           net.IPv4(1, 2, 3, 4),
		EnvelopeFrom: &mail.Address{Address: "test@example.com"},
		SenderDomain: "example.com",
		RawData:      []byte(`Subject: test`),
	})

	assert.Equal(t, authres.ResultNone, res.Spf.Value)
	assert.Equal(t, authres.ResultNone, res.Dmarc.Value)
	assert.Equal(t, 1, len(res.Dkim))
	assert.Equal(t, authres.ResultNone, res.Dkim[0].Value)
}
