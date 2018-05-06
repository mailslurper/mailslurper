package main

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/mailslurper/mailslurper/pkg/auth/authscheme"
	"github.com/mailslurper/mailslurper/pkg/contexts"
)

func adminAuthorization(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if config.AuthenticationScheme == authscheme.NONE {
			return nil
		}

		var err error
		var s *sessions.Session
		var temp interface{}
		var ok bool
		var adminUserContext *contexts.AdminUserContext

		if s, err = session.Get("session", ctx); err != nil {
			logger.WithError(err).Errorf("There was a problem retrieving the admin session")
			return err
		}

		s.Options = &sessions.Options{
			Path:   "/",
			MaxAge: 0,
		}

		if temp, ok = s.Values["user"]; !ok {
			// if ctx.Path() == "/login" {
			// 	return next(ctx)
			// }

			return ctx.Redirect(http.StatusTemporaryRedirect, "/login")
		}

		adminUserContext = &contexts.AdminUserContext{
			Context: ctx,
			User:    temp.(string),
		}

		return next(adminUserContext)
	}
}
