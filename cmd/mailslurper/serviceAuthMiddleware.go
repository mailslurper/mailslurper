// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

import (
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/mailslurper/mailslurper/pkg/auth/authscheme"
	jwtservice "github.com/mailslurper/mailslurper/pkg/auth/jwt"
	"github.com/mailslurper/mailslurper/pkg/contexts"
)

func serviceAuthorization(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var err error
		var token *jwt.Token

		adminUserContext := &contexts.AdminUserContext{
			Context: ctx,
			User:    "",
		}

		if config.AuthenticationScheme == authscheme.NONE {
			return next(adminUserContext)
		}

		jwtService := &jwtservice.JWTService{
			Config: config,
		}

		logger.Debugf("Starting parse of JWT token")
		tokenFromHeader := ctx.Request().Header.Get("Authorization")
		tokenFromHeader = strings.TrimPrefix(tokenFromHeader, "Bearer ")

		if tokenFromHeader == "" {
			logger.Errorf("No bearer and token in authorization header")
			return ctx.String(http.StatusForbidden, "Unauthorized")
		}

		if token, err = jwtService.Parse(tokenFromHeader, config.AuthSecret); err != nil {
			logger.WithError(err).Errorf("Error parsing JWT token in service authorization middleware")
			return ctx.String(http.StatusForbidden, "Error parsing token")
		}

		if err = jwtService.IsTokenValid(token); err != nil {
			logger.WithError(err).Errorf("Invalid token")
			return ctx.String(http.StatusForbidden, "Invalid token")
		}

		adminUserContext.User = jwtService.GetUserFromToken(token)
		logger.WithField("user", adminUserContext.User).Debugf("Service middleware")

		if _, ok := cacheService.Get(adminUserContext.User); !ok {
			logger.WithField("user", adminUserContext.User).Errorf("User not found in JWT token cache")
			return ctx.String(http.StatusForbidden, "JWT token not found in cache")
		}

		return next(adminUserContext)
	}
}
