package main

import (
	"net/http"

	"flag"

	"os"

	"github.com/gorilla/handlers"
	"github.com/while-loop/todo/pkg"
	"github.com/while-loop/todo/pkg/config"
	"github.com/while-loop/todo/pkg/log"
)

var (
	configFile = flag.String("i", "", "path to config file")
	v          = flag.Bool("v", false, todo.Name+" version")
	laddr      = flag.String("laddr", ":8675", "local address to bind to")
)

func main() {
	flag.Parse()
	if *v {
		log.Infof("%s %s %s %s", todo.Name, todo.Version, todo.BuildTime, todo.Commit)
		return
	}

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

	log.Info("Running on " + *laddr)
	log.Fatal(http.ListenAndServe(*laddr, wrapAppHandler(app)))
}

func wrapAppHandler(handler http.Handler) http.Handler {
	h := handlers.LoggingHandler(os.Stdout, handler)
	h = handlers.CORS()(h)
	h = handlers.RecoveryHandler()(h)
	return h
}
