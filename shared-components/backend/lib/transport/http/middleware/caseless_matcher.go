package middleware

import (
	"strings"

	"github.com/kataras/iris/v12"
)

func CaselessMatcherMiddleware(ctx iris.Context) {
	ctx.Request().URL.Path = strings.ToLower(ctx.Path())
	ctx.Next()
}
