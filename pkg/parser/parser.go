package parser

import (
	"github.com/while-loop/todo/pkg/log"
	"github.com/while-loop/todo/pkg/tracker"
	"net/http"
	"sync"
	"time"
)

const (
	maxGoroutines = 5
)

var (
	client = &http.Client{Timeout: 30 * time.Second}
)

type TodoParser interface {
	GetTodos(urls ...string) []tracker.Issue
}

type todoParser struct {
}

func New() TodoParser {
	return &todoParser{}
}

func (p *todoParser) GetTodos(urls ...string) []tracker.Issue {

	jobs := make(chan string, 100)
	results := make(chan tracker.Issue, 100)
	finished := make(chan struct{})

	wg := &sync.WaitGroup{}

	log.Debugf("spinning %d goroutines for todoParser", maxGoroutines)
	for i := 0; i < maxGoroutines; i++ {
		wg.Add(1) // track how many goroutines we spin up
		go worker(wg, jobs, results)
	}

	go func() {
		log.Debug("sending urls to worker routines")
		for _, u := range urls {
			jobs <- u
		}
		close(jobs) // close the jobs channel so goroutines gracefully stop when no jobs are left
	}()

	issues := make([]tracker.Issue, 0)
	go func() {
		log.Debug("collecting results from workers")
		for issue := range results {
			issues = append(issues, issue)
		}
		finished <- struct{}{} // let the main thread know we're done collecting issues
	}()

	log.Debug("waiting for goroutines to finish")
	// wait for all the goroutines to gracefully finish
	wg.Wait()

	// tell the result collector that we're done waiting on worker goroutines
	close(results)

	log.Debug("waiting for collector routine to finish")
	// wait for the result collector to finish appending issues
	<-finished

	log.Debug("done waiting for collector")
	return issues
}

// worker downloads and parses a file given a url
func worker(wg *sync.WaitGroup, urlChan <-chan string, results chan<- tracker.Issue) {
	for u := range urlChan {
		log.Info("worker recvd ", u)
		rc, err := DownloadFile(client, u)
		if err != nil {
			log.Error("worker failed to download file", err)
			continue
		}

		log.Debug("working parsing file")
		iss, err := ParseFile(u, rc)
		rc.Close()
		if err != nil {
			log.Error("worker failed to parse file", err)
			// don't return because we could have recvd partial issues w/ an error
		}

		if len(iss) > 0 {
			for _, is := range iss {
				results <- is
			}
		}
	}
	wg.Done()
}
