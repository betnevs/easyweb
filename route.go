package main

import (
	"fmt"

	"github.com/betNevS/easyweb/core/middleware"

	"github.com/betNevS/easyweb/core"
)

func RegisterRouter(c *core.Core) {
	c.Use(middleware.Recovery())
	c.Get("/user/login", middleware.Cost(), func(ctx *core.Context) error {
		ctx.JSON("user login")
		return nil
	})

	c.Get("/", func(ctx *core.Context) error {
		fmt.Println(ctx.GetRequest().URL.Path)
		ctx.JSON("xxxx")
		return nil
	})

}
