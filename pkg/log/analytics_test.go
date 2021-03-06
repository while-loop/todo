package log

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stretchr/testify/assert"
	"github.com/while-loop/todo/pkg/issue"
	"testing"
)

const (
	defUrl    = "https://search-todo-zwjfwycj3n3zhfhw6wydtt6sb4.us-east-2.es.amazonaws.com"
	testIndex = "test-todo"
)

func TestAWS_Login(t *testing.T) {
	needsAWS(t)
	_, err := newESLogger(testIndex, defUrl)
	assert.NoError(t, err)
}

func TestPutCommit(t *testing.T) {
	needsAWS(t)
	logger, err := newESLogger(testIndex, defUrl)
	assert.NoError(t, err)
	assert.NoError(t, logger.LogCommit("help", "me", "pls", "no"))
}

func TestPutIssue(t *testing.T) {
	needsAWS(t)
	logger, err := newESLogger(testIndex, defUrl)
	assert.NoError(t, err)
	assert.NoError(t, logger.LogIssue(&issue.Issue{
		Title:  "halp",
		Author: "while-loop",
	}))
}

func TestPutInstall(t *testing.T) {
	needsAWS(t)
	logger, err := newESLogger(testIndex, defUrl)
	assert.NoError(t, err)
	assert.NoError(t, logger.LogInstallation("created", "while-loop"))
}

func TestPutRepoInstall(t *testing.T) {
	needsAWS(t)
	logger, err := newESLogger(testIndex, defUrl)
	assert.NoError(t, err)
	assert.NoError(t, logger.LogRepoInstallation("created", "while-loop", []string{"myrepo1", "myrepo2"}))
}

func needsAWS(t *testing.T) {
	sesh, err := session.NewSession()
	if err != nil {
		t.Skip("aws_creds:", err.Error())
	}
	if sesh.Config.Credentials == nil {
		t.Skip("aws_creds are nil")
	}

	_, err = sesh.Config.Credentials.Get()
	if err != nil {
		t.Skip("aws_creds:", err.Error())
	}
}
