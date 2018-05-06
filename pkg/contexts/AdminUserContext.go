package contexts

import "github.com/labstack/echo"

type AdminUserContext struct {
	echo.Context
	User string
}
