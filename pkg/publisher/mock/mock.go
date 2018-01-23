package mock

import (
	"os"

	"fmt"

	"github.com/while-loop/todo/pkg/issue"
	"github.com/while-loop/todo/pkg/publisher/config"
)

const (
	name = "mock"
)

type MockPub struct {
	out *os.File
}

func NewPublisher(config *config.MockConfig) *MockPub {
	o := os.Stdout
	if config.Output == "stderr" {
		o = os.Stderr
	}

	return &MockPub{out: o}
}

func (s *MockPub) Name() string {
	return name
}

func (s *MockPub) Publish(issue *issue.Issue) error {
	fmt.Fprintln(s.out, issue)
	return nil
}
