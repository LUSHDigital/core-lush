package lushauthmw

import (
	"net/http"
	"strings"

	"github.com/LUSHDigital/core-lush/lushauth"
	"github.com/LUSHDigital/core/auth"
	"github.com/LUSHDigital/core/rest"
)

const (
	authHeader               = "Authorization"
	authHeaderPrefix         = "Bearer "
	msgMissingRequiredGrants = "missing required grants"
)

// MiddlewareFunc is a function which receives an http.Handler and returns another http.Handler.
// Typically, the returned handler is a closure which does something with the http.ResponseWriter and http.Request passed
// to it, and then calls the handler passed as parameter to the MiddlewareFunc.
type MiddlewareFunc func(http.Handler) http.Handler

// Middleware allows MiddlewareFunc to implement the middleware interface.
func (mw MiddlewareFunc) Middleware(handler http.Handler) http.Handler {
	return mw(handler)
}

// JWTMiddleware returns the middleware function for a jwt.
func JWTMiddleware(cr CopierRenewer) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return JWTHandler(cr, next.ServeHTTP)
	}
}

// JWTHandler takes a JWT from the request headers, attempts validation and returns a http handler.
func JWTHandler(cr CopierRenewer, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := strings.TrimPrefix(r.Header.Get(authHeader), authHeaderPrefix)
		pk := cr.Copy()
		parser := auth.NewParser(&pk, lushauth.RSAKeyFunc)
		var claims lushauth.Claims
		err := parser.Parse(raw, &claims)
		if err != nil {
			switch err.(type) {
			case lushauth.JWTSigningMethodError:
				cr.Renew() // Renew the public key if there's an error validating the token signature
			}
			res := &rest.Response{Code: http.StatusUnauthorized, Message: err.Error()}
			res.WriteTo(w)
			return
		}
		ctx := lushauth.ContextWithConsumer(r.Context(), claims.Consumer)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// HandlerGrants is an HTTP handler to check that the consumer in the request context has the required grants.
func HandlerGrants(grants []string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		consumer := lushauth.ConsumerFromContext(r.Context())
		if !consumer.HasAnyGrant(grants...) {
			res := &rest.Response{Code: http.StatusUnauthorized, Message: msgMissingRequiredGrants}
			res.WriteTo(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// HandlerRoles is an HTTP handler to check that the consumer in the request context has the required roles.
func HandlerRoles(roles []string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		consumer := lushauth.ConsumerFromContext(r.Context())
		if !consumer.HasAnyRole(roles...) {
			res := &rest.Response{Code: http.StatusUnauthorized, Message: msgMissingRequiredGrants}
			res.WriteTo(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}
