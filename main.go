package main

import (
	"log"
	"net/http"

	"github.com/betNevS/easyhttp/core"
)

func main() {
	core := core.New()
	RegisterRouter(core)
	s := &http.Server{
		Handler: core,
		Addr:    ":8080",
	}

	log.Fatal(s.ListenAndServe())
}
