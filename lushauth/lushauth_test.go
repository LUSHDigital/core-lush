package lushauth_test

import (
	"corelush/lushauth"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"os"
	"testing"
	"time"
)

var (
	now           = time.Unix(1257894000, 0)
	grantTagRoles = []string{
		"test.foo",
		"test.bar",
		"test.baz",
	}
	consumer = lushauth.Consumer{
		ID:        1,
		UUID:      "1234",
		FirstName: "John",
		LastName:  "Doe",
		Language:  "en",
		Grants:    grantTagRoles,
		Roles:     grantTagRoles,
		Needs:     grantTagRoles,
	}
	validClaims = lushauth.Claims{
		ID:        "1234",
		Issuer:    "Test Suite",
		Audience:  "Testers",
		Subject:   "Test",
		ExpiresAt: now.Add(23 * time.Hour).Unix(),
		IssuedAt:  now.Add(-1 * time.Hour).Unix(),
		NotBefore: now.Add(-1 * time.Hour).Unix(),
		Consumer:  consumer,
	}
	expiredClaims = lushauth.Claims{
		ID:        "1234",
		Issuer:    "Test Suite",
		Audience:  "Testers",
		Subject:   "Test",
		ExpiresAt: now.Add(-1 * time.Hour).Unix(),
		IssuedAt:  now.Add(-24 * time.Hour).Unix(),
		NotBefore: now.Add(-24 * time.Hour).Unix(),
		Consumer:  consumer,
	}
	invalidClaims = lushauth.Claims{
		ExpiresAt: now.Add(-23 * time.Hour).Unix(),
		IssuedAt:  now.Add(1 * time.Hour).Unix(),
		NotBefore: now.Add(1 * time.Hour).Unix(),
		Consumer:  consumer,
	}
	emptyClaims = lushauth.Claims{
		Consumer: consumer,
	}

	ecPriv   *ecdsa.PrivateKey
	rsaPriv  *rsa.PrivateKey
	rsaPriv2 *rsa.PrivateKey
	public   *rsa.PublicKey
)

func TestMain(m *testing.M) {
	var err error
	lushauth.TimeFunc = func() time.Time {
		return now
	}
	rsaPriv, err = rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	public = &rsaPriv.PublicKey
	rsaPriv2, err = rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	ecPriv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
