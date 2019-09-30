package lushauthmw_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LUSHDigital/core-lush/lushauth"
	"github.com/LUSHDigital/core-lush/middleware/lushauthmw"
	"github.com/LUSHDigital/core/rest"
	"github.com/LUSHDigital/core/test"
	"github.com/LUSHDigital/core/workers/keybroker/keybrokermock"
)

func TestJWTHandler(t *testing.T) {
	cases := []struct {
		name               string
		token              string
		expectedStatusCode int
	}{
		{
			name:               "token is good",
			token:              validToken,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "token is missing",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "token is malformed",
			expectedStatusCode: http.StatusUnauthorized,
			token:              "i am invalid",
		},
		{
			name:               "token has expired",
			token:              expiredToken,
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "token is not ready yet",
			token:              futureToken,
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "IAT is in the future",
			token:              unissuedToken,
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "token not signed with matching key",
			token:              invalidToken,
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			broker := keybrokermock.MockRSAPublicKey(public)
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			if c.token != "" {
				req.Header.Add("Authorization", "Bearer "+c.token)
			}
			recorder := httptest.NewRecorder()
			handler := lushauthmw.JWTHandler(broker, func(w http.ResponseWriter, r *http.Request) {
				consumer := lushauth.ConsumerFromContext(r.Context())
				rest.Response{Code: http.StatusOK, Message: "", Data: &rest.Data{Type: "consumer", Content: consumer}}.WriteTo(w)
			})
			handler.ServeHTTP(recorder, req)
			test.Equals(t, c.expectedStatusCode, recorder.Code)
			if c.expectedStatusCode == http.StatusOK {
				var consumer lushauth.Consumer
				rest.UnmarshalJSONResponse(recorder.Body.Bytes(), &consumer)
				test.NotEquals(t, "", consumer.ID)
			}
		})
	}
}
