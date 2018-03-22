package github

import (
	"context"

	"strconv"

	"fmt"
	"net/http"

	"strings"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"github.com/while-loop/todo/pkg/issue"
	"github.com/while-loop/todo/pkg/tracker/config"
)

const (
	name     = "github"
	ownerTag = "{owner}"
	repoTag  = "{repo}"
	shaTag   = "{sha}"
	pathTag  = "{path}"
)

var (
	blobUrl = fmt.Sprintf("https://github.com/%s/%s/blob/%s/%s", ownerTag, repoTag, shaTag, pathTag)
)

// code snippet https://github.com/while-loop/todo/blob/cc6b554cccfd3598f6b6342d69c78abcbc5d0128/pkg/app.go#L17-L25
// footer  ###### This issue was generated by [todo](https://github.com/while-loop/todo) on behalf of %s.
type tracker struct {
	conf *config.GithubConfig
}

func NewTracker(cfg *config.GithubConfig) *tracker {

	return &tracker{
		conf: cfg,
	}
}

func (t *tracker) GetIssues(ctx context.Context, ref *issue.Issue) ([]*issue.Issue, error) {
	c, err := t.client(ref.GetInt("installation"))
	if err != nil {
		return nil, err
	}

	gIss, _, err := c.Issues.ListByRepo(ctx, ref.Owner, ref.Repo, &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{PerPage: 50},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get issues for %s/%s", ref.Owner, ref.Repo)
	}

	// todo support pagination
	iss := make([]*issue.Issue, 0)

	for _, is := range gIss {
		iss = append(iss, ghIssue2todoIssue(ref.Owner, ref.Repo, is))
	}

	return iss, nil
}

func (t *tracker) CreateIssue(ctx context.Context, issue *issue.Issue) (*issue.Issue, error) {
	c, err := t.client(issue.GetInt("installation"))
	if err != nil {
		return nil, err
	}

	issue.Description += fmt.Sprintf("\n\n%s\n\n", createBlobUrl(issue))
	issue.Description += fmt.Sprintf("\n\n###### This issue was generated by [todo](https://github.com/while-loop/todo) on behalf of [%s](https://github.com/%s)", issue.Author, issue.Author)
	is, _, err := c.Issues.Create(ctx, issue.Owner, issue.Repo, &github.IssueRequest{
		Title:    &issue.Title,
		Body:     &issue.Description,
		Labels:   &issue.Labels,
		Assignee: &issue.Assignee,
		State:    pString("open"),
	})

	if err != nil {
		return nil, errors.Wrapf(err, "failed to create issue %s/%s", issue.Owner, issue.Repo)
	}

	return ghIssue2todoIssue(issue.Owner, issue.Repo, is), nil
}

func (t *tracker) DeleteIssue(ctx context.Context, issue *issue.Issue) error {
	c, err := t.client(ctx.Value("installation").(int))
	if err != nil {
		return err
	}

	iID, _ := strconv.Atoi(issue.ID)
	_, resp, err := c.Issues.Edit(ctx, issue.Owner, issue.Repo, iID, &github.IssueRequest{
		State: pString("closed"),
	})

	if err != nil {
		return errors.Wrapf(err, "failed to close issue %s/%s/%d", issue.Owner, issue.Repo, iID)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to close issue %s/%s/%d.. http status: %d", issue.Owner, issue.Repo, iID, resp.StatusCode)
	}

	return nil
}

func (t *tracker) Name() string {
	return name
}

func (t *tracker) client(installationID int) (*github.Client, error) {
	itr, err := ghinstallation.New(http.DefaultTransport, t.conf.IssueNumber, installationID, []byte(t.conf.PrivateKey))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get client for tracker")
	}

	return github.NewClient(&http.Client{
		Transport: itr,
	}), nil
}

func createBlobUrl(is *issue.Issue) string {
	fpath := is.File[strings.Index(is.File, is.Commit+"/")+len(is.Commit)+1:]
	lines := is.GetInt("lines")
	bUrl := strings.NewReplacer(pathTag, fpath, repoTag, is.Repo, ownerTag, is.Owner, shaTag, is.Commit).Replace(blobUrl)

	start := is.Line
	if lines != 0 {
		start -= 3
	}

	if start < 1 {
		start = 1
	}
	end := start + 6

	if end > lines {
		end = lines
	}

	if lines == 0 && start != 1 {
		bUrl += fmt.Sprintf("#L%d", start)
	} else if start < end && lines != 0 {
		bUrl += fmt.Sprintf("#L%d-L%d", start, end)
	}

	return bUrl
}

func parseLabels(gLs []github.Label) []string {
	ls := []string{}
	for _, gL := range gLs {
		ls = append(ls, gL.GetName())
	}
	return ls
}

func pString(s string) *string {
	return &s
}

func ghIssue2todoIssue(owner, repo string, gIs *github.Issue) *issue.Issue {
	return &issue.Issue{
		ID:          strconv.Itoa(int(gIs.GetID())),
		Title:       gIs.GetTitle(),
		Description: gIs.GetBody(),
		Assignee:    gIs.GetAssignee().GetName(),
		Author:      gIs.GetUser().GetName(),
		Mentions:    []string{},
		Labels:      parseLabels(gIs.Labels),
		File:        "",
		Line:        0,
		Owner:       owner,
		Repo:        repo,
	}
}
