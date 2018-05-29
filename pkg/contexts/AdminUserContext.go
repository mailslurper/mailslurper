// Copyright 2013-2018 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

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
