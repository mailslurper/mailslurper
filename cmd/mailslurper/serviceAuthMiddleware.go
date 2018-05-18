package main

import (
	"net/http"

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

		if token, err = jwtService.Parse(ctx.Request(), config.AuthSecret); err != nil {
			logger.WithError(err).Errorf("Error parsing JWT token in service authorization middleware")
			return ctx.String(http.StatusForbidden, "Error parsing token")
		}

		if err = jwtService.IsTokenValid(token); err != nil {
			logger.WithError(err).Errorf("Invalid token")
			return ctx.String(http.StatusForbidden, "Invalid token")
		}

		adminUserContext.User = jwtService.GetUserFromToken(token)

		logger.WithField("user", adminUserContext.User).Debugf("Service middleware")
		return next(adminUserContext)
	}
}
