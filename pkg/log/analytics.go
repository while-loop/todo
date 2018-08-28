package log

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
	"github.com/sha1sum/aws_signing_client"
	"github.com/while-loop/todo/pkg/issue"
	"net/http"
	"time"
)

const (
	defIndex = "todo"
)

type AnalysisLogger interface {
	LogIssue(issue *issue.Issue) error
	LogCommit(author, owner, repo, commit string) error
	LogInstallation(action, owner string) error
	LogRepoInstallation(action, owner string, repos []string) error
}

type esLogger struct {
	client *elastic.Client
	index  string
}

func NewESLogger(endpoint string) (AnalysisLogger, error) {
	return newESLogger(defIndex, endpoint)
}

func newESLogger(index, endpoint string) (AnalysisLogger, error) {
	sesh, err := session.NewSession()
	signer := v4.NewSigner(sesh.Config.Credentials)
	awsClient, err := aws_signing_client.New(signer, &http.Client{
		Timeout: 5 * time.Second,
	}, "es", "us-east-2")
	if err != nil {
		return nil, errors.Wrap(err, "analysis")
	}

	client, err := elastic.NewClient(
		elastic.SetURL(endpoint),
		elastic.SetScheme("https"),
		elastic.SetHttpClient(awsClient),
		elastic.SetSniff(false), // See note below
	)
	if err != nil {
		return nil, errors.Wrap(err, "analysis")
	}

	if exists, err := client.IndexExists(index).Do(context.Background()); err != nil {
		return nil, errors.Wrap(err, "analysis")
	} else if !exists {
		_, err := client.CreateIndex(index).Do(context.Background())
		if err != nil {
			return nil, errors.Wrap(err, "analysis")
		}
	}

	return &esLogger{client: client, index: index}, nil
}

func (es *esLogger) LogIssue(issue *issue.Issue) error {
	var issueJson map[string]interface{}
	inrec, _ := json.Marshal(issue)
	json.Unmarshal(inrec, &issueJson)

	issueJson["timestamp"] = time.Now()
	_, err := es.client.Index().
		Index(es.index + "-issue").
		Type("issue").
		BodyJson(issueJson).
		Do(context.Background())

	if err != nil {
		return errors.Wrap(err, "logissue")
	}

	return nil
}

func (es *esLogger) LogCommit(author, owner, repo, commit string) error {
	_, err := es.client.Index().
		Index(es.index + "-commit").
		Type("commit").
		BodyJson(map[string]interface{}{
			"author":    author,
			"owner":     owner,
			"commit":    commit,
			"timestamp": time.Now(),
		}).
		Do(context.Background())

	if err != nil {
		return errors.Wrap(err, "logcommit")
	}
	return nil
}

func (es *esLogger) LogInstallation(action, owner string) error {
	_, err := es.client.Index().
		Index(es.index + "-installation").
		Type("installation").
		BodyJson(map[string]interface{}{
			"action":    action,
			"owner":     owner,
			"timestamp": time.Now(),
		}).
		Do(context.Background())

	if err != nil {
		return errors.Wrap(err, "loginstall")
	}
	return nil
}

func (es *esLogger) LogRepoInstallation(action, owner string, repos []string) error {
	_, err := es.client.Index().
		Index(es.index + "-installation_repositories").
		Type("installation_repositories").
		BodyJson(map[string]interface{}{
			"action":    action,
			"owner":     owner,
			"repos":     repos,
			"timestamp": time.Now(),
		}).
		Do(context.Background())

	if err != nil {
		return errors.Wrap(err, "logrepoinstall")
	}
	return nil
}
