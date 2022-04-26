package main

import "github.com/betNevS/easyhttp/core"

func RegisterRouter(core *core.Core) {
	core.Set("foo", FooControllerHandler)
}
