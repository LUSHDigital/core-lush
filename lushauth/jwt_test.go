package lushauth_test

import (
	"corelush/lushauth"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/LUSHDigital/core/test"
	"github.com/LUSHDigital/uuid"
	jwt "github.com/dgrijalva/jwt-go"
)

var (
	now = time.Unix(1257894000, 0)

	consumer = lushauth.Consumer{
		ID:        1,
		UUID:      "1234",
		FirstName: "John",
		LastName:  "Doe",
		Language:  "en",
		Grants: []string{
			"test.foo",
			"test.bar",
			"test.baz",
		},
		Roles: []string{
			"test.foo",
			"test.bar",
			"test.baz",
		},
		Needs: []string{
			"test.foo",
			"test.bar",
			"test.baz",
		},
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

func must(s string, err error) string {
	if err != nil {
		panic(err)
	}
	return s
}

func TestClaims_Valid(t *testing.T) {
	type Test struct {
		name        string
		token       string
		key         *rsa.PublicKey
		expectedErr error
	}
	cases := []Test{
		Test{
			name:  "valid RS256 token",
			token: must(jwt.NewWithClaims(jwt.SigningMethodRS256, &validClaims).SignedString(rsaPriv)),
		},
		Test{
			name:  "expired RS256 token",
			token: must(jwt.NewWithClaims(jwt.SigningMethodRS256, &expiredClaims).SignedString(rsaPriv)),
			expectedErr: jwt.ValidationError{
				Inner: lushauth.JWTVerificationError{
					Errors: lushauth.JWTValidationErrorExpired,
				},
			},
		},
		Test{
			name:  "valid RS256 token (invalid claims)",
			token: must(jwt.NewWithClaims(jwt.SigningMethodRS256, &invalidClaims).SignedString(rsaPriv)),
			expectedErr: jwt.ValidationError{
				Inner: lushauth.JWTVerificationError{
					Errors: lushauth.JWTValidationErrorID | lushauth.JWTValidationErrorIssuer | lushauth.JWTValidationErrorNotValidYet | lushauth.JWTValidationErrorExpired | lushauth.JWTValidationErrorUsedBeforeIssued,
				},
			},
		},
		Test{
			name:  "valid RS256 token (empty claims)",
			token: must(jwt.NewWithClaims(jwt.SigningMethodRS256, &emptyClaims).SignedString(rsaPriv)),
			expectedErr: jwt.ValidationError{
				Inner: lushauth.JWTVerificationError{
					Errors: lushauth.JWTValidationErrorID | lushauth.JWTValidationErrorIssuer | lushauth.JWTValidationErrorNotValidYet | lushauth.JWTValidationErrorExpired | lushauth.JWTValidationErrorUsedBeforeIssued,
				},
			},
		},
		Test{
			name:  "invalid RS256 token",
			token: must(jwt.NewWithClaims(jwt.SigningMethodRS256, &validClaims).SignedString(rsaPriv2)),
			expectedErr: jwt.ValidationError{
				Inner: rsa.ErrVerification,
			},
		},
		Test{
			name:  "valid RS384 token",
			token: must(jwt.NewWithClaims(jwt.SigningMethodRS384, &validClaims).SignedString(rsaPriv)),
		},
		Test{
			name:  "invalid RS384 token",
			token: must(jwt.NewWithClaims(jwt.SigningMethodRS384, &validClaims).SignedString(rsaPriv2)),
			expectedErr: jwt.ValidationError{
				Inner: rsa.ErrVerification,
			},
		},
		Test{
			name:  "valid RS512 token",
			token: must(jwt.NewWithClaims(jwt.SigningMethodRS512, &validClaims).SignedString(rsaPriv)),
		},
		Test{
			name:  "invalid RS512 token",
			token: must(jwt.NewWithClaims(jwt.SigningMethodRS512, &validClaims).SignedString(rsaPriv2)),
			expectedErr: jwt.ValidationError{
				Inner: rsa.ErrVerification,
			},
		},
		Test{
			name:  "valid ECDSA token",
			token: must(jwt.NewWithClaims(jwt.SigningMethodES512, &validClaims).SignedString(ecPriv)),
			expectedErr: jwt.ValidationError{
				Inner: lushauth.JWTSigningMethodError{"ES512"},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var claims lushauth.Claims
			token, err := jwt.ParseWithClaims(c.token, &claims, lushauth.RSAKeyFunc(public))
			test.Equals(t, c.expectedErr, err)
			if err == nil {
				c, ok := token.Claims.(*lushauth.Claims)
				test.Equals(t, true, ok)
				test.Equals(t, validClaims.ID, c.ID)
			}

		})
	}
}

func TestRefreshableClaims_Valid(t *testing.T) {
	type Test struct {
		name        string
		token       string
		key         *rsa.PublicKey
		expectedErr error
	}
	cases := []Test{
		Test{
			name:  "valid RS256 token",
			token: must(jwt.NewWithClaims(jwt.SigningMethodRS256, &validClaims).SignedString(rsaPriv)),
		},
		Test{
			name:  "expired RS256 token",
			token: must(jwt.NewWithClaims(jwt.SigningMethodRS256, &expiredClaims).SignedString(rsaPriv)),
		},
		Test{
			name:  "valid RS256 token (invalid claims)",
			token: must(jwt.NewWithClaims(jwt.SigningMethodRS256, &invalidClaims).SignedString(rsaPriv)),
			expectedErr: jwt.ValidationError{
				Inner: lushauth.JWTVerificationError{
					Errors: lushauth.JWTValidationErrorID | lushauth.JWTValidationErrorIssuer | lushauth.JWTValidationErrorNotValidYet | lushauth.JWTValidationErrorUsedBeforeIssued,
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var claims lushauth.RefreshableClaims
			token, err := jwt.ParseWithClaims(c.token, &claims, lushauth.RSAKeyFunc(public))
			test.Equals(t, c.expectedErr, err)
			if err == nil {
				c, ok := token.Claims.(*lushauth.RefreshableClaims)
				test.Equals(t, true, ok)
				test.Equals(t, validClaims.ID, c.ID)
			}

		})
	}
}

type hasAnyTest struct {
	name     string
	values   []string
	expected bool
}

func hasAnyCases(kind string) []hasAnyTest {
	return []hasAnyTest{
		hasAnyTest{
			name:     fmt.Sprintf("when using one %s that exists", kind),
			values:   []string{"test.foo"},
			expected: true,
		},
		hasAnyTest{
			name:     fmt.Sprintf("when using two %ss where one does not exist", kind),
			values:   []string{"test.foo", "doesnot.exist"},
			expected: true,
		},
		hasAnyTest{
			name:     fmt.Sprintf("when using one %s that does not exist", kind),
			values:   []string{"doesnot.exist"},
			expected: false,
		},
		hasAnyTest{
			name:     fmt.Sprintf("when using two %ss that does not exist", kind),
			values:   []string{"doesnot.exist", "has.no.access"},
			expected: false,
		},
	}
}

func TestConsumer_HasAnyGrant(t *testing.T) {
	for _, c := range hasAnyCases("grant") {
		test.Equals(t, c.expected, consumer.HasAnyGrant(c.values...))
	}
}

func TestConsumer_HasNoMatchingGrant(t *testing.T) {
	for _, c := range hasAnyCases("grant") {
		test.Equals(t, !c.expected, consumer.HasNoMatchingGrant(c.values...))
	}
}

func TestConsumer_HasAnyNeed(t *testing.T) {
	for _, c := range hasAnyCases("need") {
		test.Equals(t, c.expected, consumer.HasAnyNeed(c.values...))
	}
}

func TestConsumer_HasNoMatchingNeed(t *testing.T) {
	for _, c := range hasAnyCases("need") {
		test.Equals(t, !c.expected, consumer.HasNoMatchingNeed(c.values...))
	}
}

func TestConsumer_HasAnyRole(t *testing.T) {
	for _, c := range hasAnyCases("role") {
		test.Equals(t, c.expected, consumer.HasAnyRole(c.values...))
	}
}

func TestConsumer_HasNoMatchingRole(t *testing.T) {
	for _, c := range hasAnyCases("role") {
		test.Equals(t, !c.expected, consumer.HasNoMatchingRole(c.values...))
	}
}

func TestConsumer_IsUser(t *testing.T) {
	consumer := &lushauth.Consumer{
		ID: 1,
	}
	t.Run("when its the same user", func(t *testing.T) {
		test.Equals(t, true, consumer.IsUser(1))
	})
	t.Run("when its not the same user", func(t *testing.T) {
		test.Equals(t, false, consumer.IsUser(2))
	})
}

func TestConsumer_HasUUID(t *testing.T) {
	id1 := uuid.Must(uuid.NewV4()).String()
	id2 := uuid.Must(uuid.NewV4()).String()
	consumer := &lushauth.Consumer{
		UUID: id1,
	}
	t.Run("when its the same user", func(t *testing.T) {
		test.Equals(t, true, consumer.HasUUID(id1))
	})
	t.Run("when its not the same user", func(t *testing.T) {
		test.Equals(t, false, consumer.HasUUID(id2))
	})
}
