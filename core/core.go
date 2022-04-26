package core

import (
	"log"
	"net/http"
)

type Core struct {
	router map[string]ControllerHandler
}

func New() *Core {
	return &Core{
		router: make(map[string]ControllerHandler),
	}
}

func (c *Core) Set(url string, handler ControllerHandler) {
	c.router[url] = handler
}

func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	log.Println("core.ServeHTTP")
	ctx := NewContext(request, response)

	r := c.router["foo"]
	if r == nil {
		return
	}

	r(ctx)
}
