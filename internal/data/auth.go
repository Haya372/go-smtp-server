package data

import (
	"github.com/emersion/go-msgauth/authres"
)

type AuthResult struct {
	Spf   authres.SPFResult
	Dkim  []authres.DKIMResult
	Dmarc authres.DMARCResult
}
