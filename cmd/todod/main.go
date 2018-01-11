package main

import (
	"net/http"

	"flag"

	"github.com/while-loop/todo/pkg"
	"github.com/while-loop/todo/pkg/config"
	"github.com/while-loop/todo/pkg/log"
)

var (
	configFile = flag.String("i", "", "path to config file")
)

func main() {
	flag.Parse()

	if *configFile == "" {
		log.Fatal("config file not given")
	}

	conf, err := config.ParseFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	app := todo.New(conf)

	log.Info(app.RepoMan.Services())
	log.Info(app.PublisherMan.Publishers())
	log.Info(app.TrackerMan.Trackers())

	log.Info("Running on :8080")
	log.Fatal(http.ListenAndServe(":8080", app.Handler()))
}
