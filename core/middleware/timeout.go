package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/betNevS/easyweb/core"
)

func Timeout(d time.Duration) core.ControllerHandler {
	return func(ctx *core.Context) error {
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)

		ctxTimeout, cancel := context.WithTimeout(ctx.BaseContext(), d)
		defer cancel()

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()

			ctx.Next()

			finish <- struct{}{}
		}()

		select {
		case <-finish:
			fmt.Println("perfect finish")
		case p := <-panicChan:
			fmt.Println(p)
			ctx.ExecPanic()
		case <-ctxTimeout.Done():
			ctx.ExecTimeout()
		}
		return nil
	}
}
