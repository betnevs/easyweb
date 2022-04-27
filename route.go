package main

import (
	"net/http"
	"time"

	"github.com/betNevS/easyhttp/core/middleware"

	"github.com/betNevS/easyhttp/core"
)

func RegisterRouter(c *core.Core) {
	c.Use(middleware.Recovery())
	c.Get("/user/login", middleware.Cost(), func(ctx *core.Context) error {
		time.Sleep(time.Second)
		ctx.JSON(http.StatusOK, "user login")
		return nil
	})

	subApi := c.Group("/aa")
	{
		subApi.Use(middleware.Test2())
		subApi.Get("/bb", func(c *core.Context) error {
			c.JSON(http.StatusOK, "OK")
			return nil
		})

		subApi.Get("/cc", func(c *core.Context) error {
			c.JSON(http.StatusOK, "OK")
			return nil
		})
	}

}
