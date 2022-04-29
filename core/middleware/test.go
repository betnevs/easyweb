package middleware

import (
	"fmt"

	"github.com/betNevS/easyweb/core"
)

func Test1() core.ControllerHandler {
	return func(ctx *core.Context) error {
		fmt.Println("middleware pre test1")
		ctx.Next()
		fmt.Println("middleware post test1")
		return nil
	}
}

func Test2() core.ControllerHandler {
	return func(ctx *core.Context) error {
		fmt.Println("middleware pre test2")
		ctx.Next()
		fmt.Println("middleware post test2")
		return nil
	}
}

func Test3() core.ControllerHandler {
	return func(ctx *core.Context) error {
		fmt.Println("middleware pre test3")
		ctx.Next()
		fmt.Println("middleware post test3")
		return nil
	}
}
