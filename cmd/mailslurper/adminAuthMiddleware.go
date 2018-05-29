// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

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
		var err error
		var s *sessions.Session
		var temp interface{}
		var ok bool

		adminUserContext := &contexts.AdminUserContext{
			Context: ctx,
			User:    "",
		}

		if config.AuthenticationScheme == authscheme.NONE {
			return next(adminUserContext)
		}

		if s, err = session.Get("session", ctx); err != nil {
			logger.WithError(err).Errorf("There was a problem retrieving the admin session")
			return err
		}

		s.Options = &sessions.Options{
			Path:   "/",
			MaxAge: 0,
		}

		if temp, ok = s.Values["user"]; !ok {
			return ctx.Redirect(http.StatusTemporaryRedirect, "/login")
		}

		adminUserContext.User = temp.(string)

		logger.WithField("user", adminUserContext.User).Debugf("Admin middleware")
		return next(adminUserContext)
	}
}
