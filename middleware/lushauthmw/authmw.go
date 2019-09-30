package lushauthmw

import (
	"crypto/rsa"
)

// CopierRenewer represents the combination of a Copier and Renewer interface
type CopierRenewer interface {
	Copy() rsa.PublicKey
	Renew()
}
