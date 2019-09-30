package lushauthmw_test

import (
	"context"
	"log"
	"testing"

	"github.com/LUSHDigital/core-lush/lushauth"
	"github.com/LUSHDigital/core-lush/middleware/lushauthmw"
	"github.com/LUSHDigital/core/middleware"
	"github.com/LUSHDigital/core/test"
	"github.com/LUSHDigital/core/workers/keybroker/keybrokermock"
	"github.com/LUSHDigital/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestGRPCMiddleware(t *testing.T) {
	claims, err := lushauth.NewClaimsForConsumer("Test", lushauth.Consumer{
		ID:        1,
		UUID:      uuid.Must(uuid.NewV4()).String(),
		FirstName: "John",
		LastName:  "Doe",
		Language:  "en",
		Grants:    []string{"read"},
		Roles:     []string{"guest"},
		Needs:     []string{"password_reset"},
	})
	if err != nil {
		log.Fatalln(err)
	}
	token, err := issuer.Issue(&claims)
	if err != nil {
		log.Fatalln(err)
	}
	brk := keybrokermock.MockRSAPublicKey(public)
	middleware.ChainStreamServer(
		lushauthmw.StreamServerInterceptor(brk),
	)
	middleware.ChainUnaryServer(
		lushauthmw.UnaryServerInterceptor(brk),
	)
	middleware.ChainStreamClient(
		lushauthmw.StreamClientInterceptor(token),
	)
	middleware.ChainUnaryClient(
		lushauthmw.UnaryClientInterceptor(token),
	)
}

func TestInterceptServerJWT(t *testing.T) {
	cases := []struct {
		name    string
		token   string
		errors  bool
		code    codes.Code
		message string
	}{
		{
			name:   "valid claims",
			token:  validToken,
			errors: false,
			code:   codes.OK,
		},
		{
			name:    "missing token",
			errors:  true,
			code:    codes.InvalidArgument,
			message: "metadata missing: auth-token",
		},
		{
			name:    "malformed token",
			token:   "123",
			errors:  true,
			code:    codes.InvalidArgument,
			message: "token contains an invalid number of segments",
		},
		{
			name:    "incorrect signing method",
			token:   "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJjb25zdW1lciI6eyJpZCI6OTk5LCJmaXJzdF9uYW1lIjoiVGVzdHkiLCJsYXN0X25hbWUiOiJNY1Rlc3QiLCJsYW5ndWFnZSI6IiIsImdyYW50cyI6WyJ0ZXN0aW5nLnJlYWQiLCJ0ZXN0aW5nLmNyZWF0ZSJdfSwiZXhwIjoxNTE4NjAzNzIwLCJqdGkiOiIyNTAwYjk3MS0wNTcxLTQ4Y2UtYmUzOS1jNWJhNGQwZmU0MGIiLCJpc3MiOiJ0ZXN0aW5nIn0.",
			errors:  true,
			code:    codes.InvalidArgument,
			message: "unexpected signing method (needs to be RSA): none",
		},
		{
			name:    "invalid claims",
			token:   invalidToken,
			errors:  true,
			code:    codes.InvalidArgument,
			message: "crypto/rsa: verification error",
		},
		{
			name:    "expired token",
			token:   expiredToken,
			errors:  true,
			code:    codes.Unauthenticated,
			message: "could not verify token: has expired",
		},
		{
			name:    "token not valid yet",
			token:   futureToken,
			errors:  true,
			code:    codes.Unauthenticated,
			message: "could not verify token: is not valid yet",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			broker := keybrokermock.MockRSAPublicKey(public)
			md := metadata.MD{}
			if c.token != "" {
				md.Set("auth-token", c.token)
			}
			ctx := metadata.NewIncomingContext(context.Background(), md)
			_, err := lushauthmw.InterceptServerJWT(ctx, broker)
			if c.errors {
				s, ok := status.FromError(err)
				if !ok {
					t.Errorf("unknown status from err: %v", err)
				}
				test.Equals(t, c.message, s.Message())
				test.Equals(t, c.code, s.Code())
			} else {
				test.Equals(t, nil, err)
			}
		})
	}
}

func TestContextWithAuthTokenMetadata(t *testing.T) {
	cases := []struct {
		name  string
		token string
	}{
		{
			name:  "invalid claims",
			token: validToken,
		},
		{
			name:  "invalid claims",
			token: invalidToken,
		},
		{
			name:  "expired token",
			token: expiredToken,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			out := lushauthmw.ContextWithAuthTokenMetadata(ctx, c.token)
			md, ok := metadata.FromOutgoingContext(out)
			test.Equals(t, true, ok)
			test.Equals(t, c.token, md.Get("auth-token")[0])
		})
	}
}
