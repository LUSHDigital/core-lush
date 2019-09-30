package lushauthmw_test

import (
	"crypto/rsa"
	"log"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/LUSHDigital/core/auth/authmock"
	"github.com/LUSHDigital/uuid"
	jwt "github.com/dgrijalva/jwt-go"

	"github.com/LUSHDigital/core-lush/lushauth"
	"github.com/LUSHDigital/core-lush/middleware/lushauthmw"
	"github.com/LUSHDigital/core/auth"
	"google.golang.org/grpc"
)

var (
	private, incorrectPrivate                                                            *rsa.PrivateKey
	public, incorrectPublic                                                              *rsa.PublicKey
	issuer, incorrectIssuer                                                              *auth.Issuer
	parser, incorrectParser                                                              *auth.Parser
	broker                                                                               lushauthmw.CopierRenewer
	now, then, at                                                                        time.Time
	validClaims, otherClaims, invalidClaims, expiredClaims, futureClaims, unissuedClaims lushauth.Claims
	validToken, otherToken, invalidToken, expiredToken, futureToken, unissuedToken       string
)

func mustIssue(s string, err error) string {
	if err != nil {
		panic(err)
	}
	return s
}

func mustClaims(c lushauth.Claims, err error) lushauth.Claims {
	if err != nil {
		panic(err)
	}
	return c
}

func TestMain(m *testing.M) {
	now = time.Now()
	then = now.Add(-76 * time.Hour)
	at = now.Add(76 * time.Hour)
	lushauth.TimeFunc = func() time.Time { return now }
	jwt.TimeFunc = func() time.Time { return now }
	private, public = authmock.MustNewRSAKeyPair()
	issuer, parser = authmock.NewRSAIssuerAndParserFromKeyPair(private, public)
	incorrectPrivate, incorrectPublic = authmock.MustNewRSAKeyPair()
	incorrectIssuer, incorrectParser = authmock.NewRSAIssuerAndParserFromKeyPair(incorrectPrivate, incorrectPublic)
	defaultConsumer := lushauth.Consumer{
		ID:        999,
		UUID:      uuid.Must(uuid.NewV4()).String(),
		FirstName: "John",
		LastName:  "Doe",
		Language:  "en",
		Grants:    []string{},
		Roles:     []string{},
		Needs:     []string{},
	}
	validClaims = mustClaims(lushauth.NewClaimsForConsumer("Test", defaultConsumer))
	validToken = mustIssue(issuer.Issue(&validClaims))
	otherClaims = mustClaims(lushauth.NewClaimsForConsumer("", defaultConsumer))
	otherToken = mustIssue(incorrectIssuer.Issue(&otherClaims))
	invalidClaims = lushauth.Claims{
		IssuedAt:  now.Add(-2 * time.Hour).Unix(),
		NotBefore: now.Add(-1 * time.Hour).Unix(),
		ExpiresAt: now.Add(-1 * time.Minute).Unix(),
		Consumer:  defaultConsumer,
	}
	invalidToken = mustIssue(incorrectIssuer.Issue(&invalidClaims))
	expiredClaims = lushauth.Claims{
		ID:        uuid.Must(uuid.NewV4()).String(),
		Issuer:    "Test",
		IssuedAt:  now.Add(-2 * time.Hour).Unix(),
		NotBefore: now.Add(-1 * time.Hour).Unix(),
		ExpiresAt: now.Add(-1 * time.Minute).Unix(),
		Consumer:  defaultConsumer,
	}
	expiredToken = mustIssue(issuer.Issue(&expiredClaims))
	futureClaims = lushauth.Claims{
		ID:        uuid.Must(uuid.NewV4()).String(),
		Issuer:    "Test",
		IssuedAt:  now.Add(-2 * time.Hour).Unix(),
		NotBefore: now.Add(1 * time.Minute).Unix(),
		ExpiresAt: now.Add(1 * time.Hour).Unix(),
		Consumer:  defaultConsumer,
	}
	futureToken = mustIssue(issuer.Issue(&futureClaims))
	unissuedClaims = lushauth.Claims{
		ID:        uuid.Must(uuid.NewV4()).String(),
		Issuer:    "Test",
		IssuedAt:  now.Add(1 * time.Hour).Unix(),
		NotBefore: now.Add(1 * time.Minute).Unix(),
		ExpiresAt: now.Add(1 * time.Hour).Unix(),
		Consumer:  defaultConsumer,
	}
	unissuedToken = mustIssue(issuer.Issue(&unissuedClaims))
	os.Exit(m.Run())
}

func ExampleNewStreamServerInterceptor() {
	srv := grpc.NewServer(
		lushauthmw.NewStreamServerInterceptor(broker),
	)
	l, err := net.Listen("tpc", ":50051")
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(srv.Serve(l))
}

func ExampleNewUnaryServerInterceptor() {
	srv := grpc.NewServer(
		lushauthmw.NewUnaryServerInterceptor(broker),
	)

	l, err := net.Listen("tpc", ":50051")
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(srv.Serve(l))
}

func ExampleJWTHandler() {
	http.Handle("/users", lushauthmw.JWTHandler(broker, func(w http.ResponseWriter, r *http.Request) {
		consumer := lushauth.ConsumerFromContext(r.Context())
		if !consumer.HasAnyGrant("users.read") {
			http.Error(w, "access denied", http.StatusUnauthorized)
		}
	}))
}
