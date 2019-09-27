package lushauth

import (
	"crypto/rsa"
	"fmt"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

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
func RSAKeyFunc(pk *rsa.PublicKey) func(token *jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return pk, JWTSigningMethodError{token.Header["alg"]}
		}
		return pk, nil
	}
}

// Claims hold information of the power exherted by a JWT.
// Structured version of Claims Section, as referenced at
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

// RefreshableClaims hold information of the power exherted by a JWT.
// Structured version of Claims Section, as referenced at
// https://tools.ietf.org/html/rfc7519#section-4.1
//
// The difference between RefreshableClaims and Claims is that this struct
// will not attempt to validate wether or not the token is expired.
type RefreshableClaims struct {
	Claims
}

// Valid verifys time based claims (IAT, NBF) as well as the identifiers (ISS, JTI).
func (c *RefreshableClaims) Valid() error {
	now := TimeFunc()
	errors := c.verify(now, uint32(0))
	if errors > 0 {
		return JWTVerificationError{errors}
	}
	return nil
}

// Consumer represents an API user for the LUSH infrastructure.
type Consumer struct {
	// ID is a unique identifier for a user but should not be used in favour of UUID.
	ID int64 `json:"id"`
	// UUID is the unique identifier for a user.
	UUID string `json:"uuid"`
	// FirstName is the given name of a user.
	FirstName string `json:"first_name"`
	// LastName is the surname of a user.
	LastName string `json:"last_name"`
	// Language is the preferred language of a user.
	Language string `json:"language"`
	// Grants are any specific given permissions for a user.
	// e.g. products.create, pages.read or  tills.close
	Grants []string `json:"grants"`
	// Roles are what purpose a user server within the context of LUSH
	// e.g. guest, staff, creator or admin
	Roles []string `json:"roles"`
	// Needs are things that the user needs to do and that a front-end can react to.
	// e.g. password_reset, confirm_email or accept_terms
	Needs []string `json:"needs"`
}

// HasAnyGrant checks if a consumer possess any of a given set of grants
func (c *Consumer) HasAnyGrant(grants ...string) bool {
	return hasAny(c.Grants, grants...)
}

// HasNoMatchingGrant checks if a consumer is missing any of a given set of grants
func (c Consumer) HasNoMatchingGrant(grants ...string) bool {
	return hasNoMatching(c.Grants, grants...)
}

// HasAnyRole checks if a consumer possess any of a given set of roles
func (c *Consumer) HasAnyRole(roles ...string) bool {
	return hasAny(c.Roles, roles...)
}

// HasNoMatchingRole checks if a consumer is missing any of a given set of roles
func (c *Consumer) HasNoMatchingRole(roles ...string) bool {
	return hasNoMatching(c.Roles, roles...)
}

// HasAnyNeed checks if a consumer has any of the given needs
func (c *Consumer) HasAnyNeed(needs ...string) bool {
	return hasAny(c.Needs, needs...)
}

// HasNoMatchingNeed checks if a consumer has any of the given needs
func (c *Consumer) HasNoMatchingNeed(needs ...string) bool {
	return hasNoMatching(c.Needs, needs...)
}

// IsUser checks if a consumer has the same ID as a user
func (c *Consumer) IsUser(userID int64) bool {
	return c.ID == userID
}

// HasUUID checks if a consumer has the same uuid as a user
func (c *Consumer) HasUUID(id string) bool {
	return c.UUID == id
}

// hasAny checks if a set contains any of the given members.
func hasAny(set []string, members ...string) bool {
	for _, member := range members {
		for _, m := range set {
			if member == m {
				return true
			}
		}
	}
	return false
}

// hasNoMatching checks if a set does not contain any and all of the given members.
func hasNoMatching(set []string, members ...string) bool {
	for _, member := range members {
		for _, m := range set {
			if member == m {
				return false
			}
		}
	}
	return true
}
