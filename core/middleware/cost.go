package middleware

import (
	"log"
	"time"

	"github.com/betNevS/easyhttp/core"
)

func Cost() core.ControllerHandler {
	return func(ctx *core.Context) error {
		start := time.Now()

		ctx.Next()

		end := time.Now()
		cost := end.Sub(start)

		log.Printf("api uri: %v, cost: %v", ctx.GetRequest().RequestURI, cost.Seconds())
		return nil
	}
}
