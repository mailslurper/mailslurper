package contexts

import "github.com/labstack/echo"

type AdminUserContext struct {
	echo.Context
	User string
}

func GetAdminContext(ctx echo.Context) *AdminUserContext {
	return ctx.(*AdminUserContext)
}
