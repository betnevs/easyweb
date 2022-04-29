package main

import (
	"github.com/betNevS/easyweb/core"
	"log"
	"net/http"
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
