package main

import (
	"net/http"

	"github.com/while-loop/todo/pkg"
	"github.com/while-loop/todo/pkg/log"
)

func main() {
	app := todo.New(nil, nil)

	router := app.Handler()
	log.Info("Running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
