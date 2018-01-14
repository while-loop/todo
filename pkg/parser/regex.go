package parser

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"github.com/while-loop/todo/pkg/log"
	"github.com/while-loop/todo/pkg/tracker"
	"io"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
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
	slashLangs    = []string{"go", "java", "c", "cpp", "h", "hpp"}
	hashLangs     = []string{"py", "sh", "bash", "yml", "yaml"}
)

func init() {
	sort.Strings(slashLangs)
	sort.Strings(hashLangs)
}

func ParseFile(fileName string, file io.ReadCloser) ([]tracker.Issue, error) {
	defer file.Close()
	issues := make([]tracker.Issue, 0)
	ext := strings.TrimLeft(filepath.Ext(fileName), ".")
	if ext == "" {
		log.Errorf("failed to get file ext for %s", fileName)
		return nil, fmt.Errorf("unknown file ext: %s", ext)
	}

	rexp := commentRegexes(ext)
	if rexp == nil {
		log.Errorf("parser regex. unknown ext type: %s", ext)
		return nil, fmt.Errorf("unknown file ext: %s", ext)
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

		is, found := parseLine(rexp, line)
		if found {
			is.File = fileName
			is.Line = lineNum
			//log.Debugf("found issue: %s\n%s", line, is.String())
			issues = append(issues, is)
		} else {
			log.Debug(line)
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

func parseLine(rexp *regexp.Regexp, line string) (tracker.Issue, bool) {

	finds := rexp.FindStringSubmatch(line)
	if len(finds) <= 0 {
		return tracker.Issue{}, false
	}

	i := tracker.Issue{
		Mentions:    mentionsRegex.FindAllString(line, -1),
		Labels:      parseLabels(line),
		File:        "",
		Line:        0,
		ID:          "",
		Title:       "",
		Assignee:    "",
		Author:      "",
		Description: "",
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
	var labels []string

	for _, groups := range labelsRegex.FindAllStringSubmatch(line, -1) {
		if len(groups) < 2 {
			log.Warn("only found 1 group in label regex", groups)
			continue
		}

		labels = append(labels, groups[1])
	}

	return labels
}
