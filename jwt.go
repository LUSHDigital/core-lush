package corelush

import (
	"crypto/rsa"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	// TimeFunc is a variable with a function to determine the current time.
	// Can be overridden in a test environment to set the current time to whatever you want it to be.
	TimeFunc = time.Now
)

const (
	// ValidationErrorExpired happens when EXP validation failed
	ValidationErrorExpired uint32 = 1 << iota
	// ValidationErrorIssuedAt happens when IAT validation failed
	ValidationErrorIssuedAt
	// ValidationErrorNotValidYet happens when NBF validation failed
	ValidationErrorNotValidYet
	// ValidationErrorIssuer happens when ISS validation failed
	ValidationErrorIssuer
	// ValidationErrorID happens when JTI validation failed
	ValidationErrorID
)

// RSAKeyFunc is used with the jwt-go library to validate that a token is using the correct signing algorithm.
func RSAKeyFunc(pk *rsa.PublicKey) func(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return pk, fmt.Sprintf("unexpected signing method: %v", token.Header["alg"])
	}
	return pk, nil
}

// NewValidationError is for constructing a ValidationError with a string error message
func NewValidationError(errorText string, errorFlags uint32) *ValidationError {
	return &ValidationError{
		text:   errorText,
		Errors: errorFlags,
	}
}

// ValidationError happens when we're unable to validate a JWT.
type ValidationError struct {
	Inner  error  // stores the error returned by external dependencies, i.e.: KeyFunc
	Errors uint32 // bitfield.  see ValidationError... constants
	text   string // errors that do not have a valid error just have text
}

// Validation error is an error type
func (e ValidationError) Error() string {
	if e.Inner != nil {
		return e.Inner.Error()
	} else if e.text != "" {
		return e.text
	}
	return "token is invalid"
}

// Valid can be used to see if the validation error has any errors.
func (e *ValidationError) Valid() bool {
	return e.Errors == 0
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
	now := TimeFunc().Unix()
	err := c.validate(now)
	if c.VerifyExpiresAt(now) == false {
		delta := time.Unix(now, 0).Sub(time.Unix(c.ExpiresAt, 0))
		err.Inner = fmt.Errorf("token is expired by %v", delta)
		err.Errors |= ValidationErrorExpired
	}
	if !err.Valid() {
		return err
	}
	return nil
}

// WasValid validates time based claims (IAT, NBF) as well as the identifiers (ISS, JTI).
func (c *Claims) WasValid() error {
	now := TimeFunc().Unix()
	err := c.validate(now)
	if !err.Valid() {
		return *err
	}
	return nil
}

func (c *Claims) validate(now int64) *ValidationError {
	err := new(ValidationError)
	if c.Issuer == "" {
		err.Inner = fmt.Errorf("token does not have an issuer")
		err.Errors |= ValidationErrorIssuer
	}
	if c.ID == "" {
		err.Inner = fmt.Errorf("token does not have an id")
		err.Errors |= ValidationErrorID
	}
	if c.VerifyIssuedAt(now) == false {
		err.Inner = fmt.Errorf("token used before issued")
		err.Errors |= ValidationErrorIssuedAt
	}
	if c.VerifyNotBefore(now) == false {
		err.Inner = fmt.Errorf("token is not valid yet")
		err.Errors |= ValidationErrorNotValidYet
	}
	return err
}

// VerifyExpiresAt compares the exp claim against a timestamp.
// Will change behaviour depending on the value of corelush.TimeFunc
func (c *Claims) VerifyExpiresAt(ts int64) bool {
	if c.ExpiresAt == 0 {
		return false
	}
	return ts <= c.ExpiresAt
}

// VerifyIssuedAt compares the iat claim against a timestamp.
// Will change behaviour depending on the value of corelush.TimeFunc
func (c *Claims) VerifyIssuedAt(ts int64) bool {
	if c.IssuedAt == 0 {
		return false
	}
	return ts >= c.IssuedAt
}

// VerifyNotBefore compares the nbf claim against a timestamp.
// Will change behaviour depending on the value of corelush.TimeFunc
func (c *Claims) VerifyNotBefore(ts int64) bool {
	if c.NotBefore == 0 {
		return false
	}
	return ts >= c.NotBefore
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

// HasAnyRole checks if a consumer possess any of a given set of roles
func (c *Consumer) HasAnyRole(roles ...string) bool {
	return hasAny(c.Roles, roles...)
}

// HasAnyNeed checks if a consumer has any of the given needs
func (c *Consumer) HasAnyNeed(needs ...string) bool {
	return hasAny(c.Needs, needs...)
}

// IsUser checks if a consumer has the same ID as a user
func (c *Consumer) IsUser(userID int64) bool {
	return c.ID == userID
}

// HasUUID checks if a consumer has the same uuid as a user
func (c *Consumer) HasUUID(id string) bool {
	return c.UUID == id
}

// hasAny checks if a set contains one or more members.
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
