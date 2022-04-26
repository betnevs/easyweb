package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/betNevS/easyhttp/core"
)

func FooControllerHandler(c *core.Context) error {
	finish := make(chan struct{}, 1)
	panicChan := make(chan interface{}, 1)

	ctx, cancel := context.WithTimeout(c.BaseContext(), 2*time.Second)
	defer cancel()

	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()

		time.Sleep(3 * time.Second)
		c.JSON(http.StatusOK, "OK")

		finish <- struct{}{}

	}()

	select {
	case <-ctx.Done():
		c.ExecTimeout()
	case <-finish:
		fmt.Println("Perfect finish")
	case p := <-panicChan:
		fmt.Println(p)
		c.ExecPanic()
	}
	return nil
}
