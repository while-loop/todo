package parser

import (
	"bufio"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"context"

	"github.com/pkg/errors"
	"github.com/while-loop/todo/pkg/issue"
	"github.com/while-loop/todo/pkg/log"
)

const (
	cmtToken = `{comment_token}`
	cmtRegex = `(?i)^\s*` + cmtToken + `.*todo:?\s*(\((?P<assignee>.*)\):?)?[\s|$]\s*(?P<title>.*)$`
	slash    = `//`
	hash     = `#`
)

var (
	slashRegex    = regexp.MustCompile(strings.Replace(cmtRegex, cmtToken, slash, 1))
	hashRegex     = regexp.MustCompile(strings.Replace(cmtRegex, cmtToken, hash, 1))
	mentionsRegex = regexp.MustCompile(`(@[^\s]+)`)
	labelsRegex   = regexp.MustCompile(`\+([^\s]+)`)
	slashLangs    = []string{"go", "java", "c", "cpp", "h", "hpp", "test"}
	hashLangs     = []string{"py", "sh", "bash", "yml", "yaml"}
)

func init() {
	sort.Strings(slashLangs)
	sort.Strings(hashLangs)
}

func ParseFile(fileName string, file io.ReadCloser) ([]*issue.Issue, error) {
	defer file.Close()
	issues := make([]*issue.Issue, 0)
	ext := strings.TrimLeft(filepath.Ext(fileName), ".")
	if ext == "" {
		log.Errorf("failed to get file ext for %s", fileName)
		return nil, fmt.Errorf("unknown file ext: %s", ext)
	}

	rexp := commentRegexes(ext)
	if rexp == nil {
		log.Warnf("parser regex. unknown ext type: %s", ext)
		return nil, nil
	}

	scan := bufio.NewScanner(file)
	scan.Split(bufio.ScanLines)

	lineNum := 0 // set as 0 so 1st iteration is 1. increment is at beginning to accommodate for continues
	for scan.Scan() {
		lineNum++
		line := scan.Text()

		// todo(while-loop): add ignore keyword to yml config (ParseFile will be a todoParser func)
		if strings.Contains(line, "!todo") {
			continue
		}

		// todo(while-loop): treat subsequent comment lines as description
		is, found := parseLine(rexp, line)
		if found {
			is.File = fileName
			is.Line = lineNum
			is.Commit = "master"
			issues = append(issues, is)
		}
	}

	if scan.Err() != nil {
		return issues, errors.Wrapf(scan.Err(), "error while scanning file: %s", fileName)
	}

	return issues, nil
}

func commentRegexes(ext string) *regexp.Regexp {
	idx := sort.SearchStrings(slashLangs, ext)
	if idx < len(slashLangs) && slashLangs[idx] == ext {
		return slashRegex
	}

	idx = sort.SearchStrings(hashLangs, ext)
	if idx < len(hashLangs) && hashLangs[idx] == ext {
		return hashRegex
	}

	return nil
}

func parseLine(rexp *regexp.Regexp, line string) (*issue.Issue, bool) {

	finds := rexp.FindStringSubmatch(line)
	if len(finds) <= 0 {
		return nil, false
	}

	ms := mentionsRegex.FindAllString(line, -1)
	if ms == nil {
		ms = []string{}
	}

	i := &issue.Issue{
		Mentions:    ms,
		Labels:      parseLabels(line),
		File:        "",
		Line:        0,
		ID:          "",
		Title:       "",
		Assignee:    "",
		Author:      "",
		Description: "",
		Ctx:         context.Background(),
	}

	for idx, name := range rexp.SubexpNames() {
		if name == "assignee" {
			i.Assignee = finds[idx]
		}
		if name == "title" {
			i.Title = finds[idx]
		}
	}

	i.Title = filterTitle(i.Title, i.Mentions, i.Labels)
	return i, true
}

func filterTitle(line string, mentions, labels []string) string {
	for _, m := range mentions {
		line = regexp.MustCompile(`\s*`+m+`\s*`).ReplaceAllString(line, "")
	}

	for _, l := range labels {
		line = regexp.MustCompile(`\s*\+`+l+`\s*`).ReplaceAllString(line, "")
	}
	return line
}

func parseLabels(line string) []string {
	labels := []string{"todo"}

	for _, groups := range labelsRegex.FindAllStringSubmatch(line, -1) {
		if len(groups) < 2 {
			log.Warn("only found 1 group in label regex", groups)
			continue
		}

		labels = append(labels, groups[1])
	}

	return labels
}
