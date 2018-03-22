package main

import (
	"net/http"

	"os"

	"github.com/akrylysov/algnhsa"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/while-loop/todo/pkg"
	"github.com/while-loop/todo/pkg/config"
	"github.com/while-loop/todo/pkg/log"
)

const (
	configKey = "TODO_CONFIG_PATH"
)

func main() {
	log.Infof("%s %s %s %s", todo.Name, todo.Version, todo.BuildTime, todo.Commit)

	if len(os.Args) >= 2 {
		return
	}

	configFile := os.Getenv(configKey)
	if configFile == "" {
		log.Fatal("config file not given")
	}

	log.Info("Using config: ", configFile)
	conf, err := config.ParseFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	app := todo.New(conf, mux.NewRouter())

	log.Info("Repo services ", app.RepoMan.Services())
	log.Info("Tracker services ", app.TrackerMan.Trackers())
	algnhsa.ListenAndServe(wrapAppHandler(app), nil)
}

func wrapAppHandler(handler http.Handler) http.Handler {
	h := handlers.LoggingHandler(os.Stdout, handler)
	h = handlers.CORS()(h)
	h = handlers.RecoveryHandler()(h)
	return h
}
