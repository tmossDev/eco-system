package middleware

import (
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"tmossDev.github.com/eco-system/shared-components/backend/package/constants"
)

func RequestIDMiddleware(ctx iris.Context) {
	requestId := ctx.Request().Header.Get(constants.CTXRequestIdKey)
	if len(requestId) == 0 {
		requestId = uuid.NewString()
	}
	ctx.Values().Set(constants.CTXRequestIdKey, requestId)
	ctx.Header(constants.CTXRequestIdKey, requestId)
	ctx.Next()
}
