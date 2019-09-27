package lushauth_test

import (
	"corelush/lushauth"
	"crypto/rsa"
	"testing"

	"github.com/LUSHDigital/core/test"
	jwt "github.com/dgrijalva/jwt-go"
)

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
