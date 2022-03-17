package middleware

import (
	"net/http"

	"github.com/strick-j/scimfe/internal/model/auth"
	"github.com/strick-j/scimfe/internal/service"
	"github.com/strick-j/scimfe/internal/web"
)

const authHeader = "X-Auth-Token"

// NewAuthMiddleware returns a new middleware which checks if user is authenticated.
//
// If user is authenticated, user session will be populated into request context.
func NewAuthMiddleware(authSvc *service.AuthService) web.MiddlewareFunc {
	return func(rw http.ResponseWriter, req *http.Request) (*http.Request, error) {
		token := req.Header.Get(authHeader)
		ssid, err := auth.ParseToken(token)
		if err != nil {
			return req, service.ErrAuthRequired
		}

		sess, err := authSvc.GetSession(req.Context(), ssid)
		if err != nil {
			return req, err
		}

		ctx := auth.ContextWithSession(req.Context(), sess)
		return req.WithContext(ctx), nil
	}
}
