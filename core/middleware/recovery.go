package middleware

import (
	"net/http"

	"github.com/betNevS/easyhttp/core"
)

func Recovery() core.ControllerHandler {
	return func(ctx *core.Context) error {
		defer func() {
			if p := recover(); p != nil {
				ctx.JSON(http.StatusInternalServerError, "INTERNAL ERROR")
			}
		}()
		ctx.Next()
		return nil
	}
}
