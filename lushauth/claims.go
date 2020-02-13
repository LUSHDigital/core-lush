package lushauth

import (
	"crypto"
	"fmt"
	"strings"
	"time"

	"github.com/LUSHDigital/uuid"

	jwt "github.com/dgrijalva/jwt-go"
)

// NewClaimsForConsumer spawns new claims for
func NewClaimsForConsumer(issuer string, consumer Consumer) (Claims, error) {
	var c Claims
	now := TimeFunc()
	id, err := uuid.NewV4()
	if err != nil {
		return c, err
	}
	return Claims{
		ID:        id.String(),
		Issuer:    issuer,
		ExpiresAt: now.Add(DefaultValidPeriod).Unix(),
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
		Consumer:  consumer,
	}, nil
}

const (
	// JWTValidationErrorExpired happens when EXP validation failed
	JWTValidationErrorExpired uint32 = 1 << iota
	// JWTValidationErrorUsedBeforeIssued happens when IAT validation failed
	JWTValidationErrorUsedBeforeIssued
	// JWTValidationErrorNotValidYet happens when NBF validation failed
	JWTValidationErrorNotValidYet
	// JWTValidationErrorIssuer happens when ISS validation failed
	JWTValidationErrorIssuer
	// JWTValidationErrorID happens when JTI validation failed
	JWTValidationErrorID
)

// JWTSigningMethodError happens when the RSA
type JWTSigningMethodError struct {
	Algorithm interface{}
}

func (e JWTSigningMethodError) Error() string {
	return fmt.Sprintf("unexpected signing method (needs to be RSA): %v", e.Algorithm)
}

// JWTVerificationError happens when one or more token fields could not be verified.
type JWTVerificationError struct {
	Errors uint32
}

func (e JWTVerificationError) Error() string {
	var messages []string
	if (e.Errors & JWTValidationErrorID) > 0 {
		messages = append(messages, "does not have an id")
	}
	if (e.Errors & JWTValidationErrorIssuer) > 0 {
		messages = append(messages, "does not have an issuer")
	}
	if (e.Errors & JWTValidationErrorExpired) > 0 {
		messages = append(messages, "has expired")
	}
	if (e.Errors & JWTValidationErrorNotValidYet) > 0 {
		messages = append(messages, "is not valid yet")
	}
	if (e.Errors & JWTValidationErrorUsedBeforeIssued) > 0 {
		messages = append(messages, "used before issued")
	}
	return fmt.Sprintf("could not verify token: %s", strings.Join(messages, ", "))
}

// RSAKeyFunc is used with the jwt-go library to validate that a token is using the correct signing algorithm.
func RSAKeyFunc(pk crypto.PublicKey) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return pk, JWTSigningMethodError{token.Header["alg"]}
		}
		return pk, nil
	}
}

// Claims hold information of the power exerted by a JWT.
// A structured version of the Claims section, as referenced at
// https://tools.ietf.org/html/rfc7519#section-4.1
type Claims struct {
	ID       string `json:"jti,omitempty"`
	Issuer   string `json:"iss,omitempty"`
	Audience string `json:"aud,omitempty"`
	Subject  string `json:"sub,omitempty"`

	ExpiresAt int64 `json:"exp,omitempty"`
	IssuedAt  int64 `json:"iat,omitempty"`
	NotBefore int64 `json:"nbf,omitempty"`

	Consumer Consumer `json:"consumer"`
}

// Valid validates time based claims (EXP, IAT, NBF) as well as the identifiers (ISS, JTI).
func (c *Claims) Valid() error {
	now := TimeFunc()
	errors := c.verify(now, uint32(0))
	errors = c.verifyExpiresAt(now, errors)
	if errors > 0 {
		return JWTVerificationError{errors}
	}
	return nil
}

func (c *Claims) verifyExpiresAt(now time.Time, errors uint32) uint32 {
	if c.VerifyExpiresAt(now) == false {
		errors |= JWTValidationErrorExpired
	}
	return errors
}

func (c *Claims) verify(now time.Time, errors uint32) uint32 {
	if c.Issuer == "" {
		errors |= JWTValidationErrorIssuer
	}
	if c.ID == "" {
		errors |= JWTValidationErrorID
	}
	if c.VerifyIssuedAt(now) == false {
		errors |= JWTValidationErrorUsedBeforeIssued
	}
	if c.VerifyNotBefore(now) == false {
		errors |= JWTValidationErrorNotValidYet
	}
	return errors
}

// VerifyExpiresAt compares the exp claim against a timestamp.
// Will change behaviour depending on the value of corelush.TimeFunc
func (c *Claims) VerifyExpiresAt(now time.Time) bool {
	if c.ExpiresAt == 0 {
		return false
	}
	return now.Unix() <= c.ExpiresAt
}

// VerifyIssuedAt compares the iat claim against a timestamp.
// Will change behaviour depending on the value of corelush.TimeFunc
func (c *Claims) VerifyIssuedAt(now time.Time) bool {
	if c.IssuedAt == 0 {
		return false
	}
	return now.Unix() >= c.IssuedAt
}

// VerifyNotBefore compares the nbf claim against a timestamp.
// Will change behaviour depending on the value of corelush.TimeFunc
func (c *Claims) VerifyNotBefore(now time.Time) bool {
	if c.NotBefore == 0 {
		return false
	}
	return now.Unix() >= c.NotBefore
}

// RefreshableClaims hold information of the power exerted by a JWT.
// A structured version of the Claims section, as referenced at
// https://tools.ietf.org/html/rfc7519#section-4.1
//
// The difference between RefreshableClaims and Claims is that this struct
// will not attempt to validate whether the token is expired.
type RefreshableClaims struct {
	Claims
}

// Valid verifies time based claims (IAT, NBF) as well as the identifiers (ISS, JTI).
func (c *RefreshableClaims) Valid() error {
	now := TimeFunc()
	errors := c.verify(now, uint32(0))
	if errors > 0 {
		return JWTVerificationError{errors}
	}
	return nil
}
