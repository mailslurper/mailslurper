package contexts

import "github.com/labstack/echo"

type AdminUserContext struct {
	echo.Context
	User string
}

func GetAdminContext(ctx echo.Context) *AdminUserContext {
	var ok bool

	if _, ok = ctx.(*AdminUserContext); ok {
		return ctx.(*AdminUserContext)
	}

	return &AdminUserContext{
		Context: ctx,
	}
}
