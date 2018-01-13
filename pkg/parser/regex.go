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
	cmtRegex = `(?i)^\s*` + cmtToken + `.*todo(\((.*)\):")?\s*(.*)$`
)

var (
	slashRegex = regexp.MustCompile(strings.Replace(cmtRegex, cmtToken, `//`, 1))
	hashRegex  = regexp.MustCompile(strings.Replace(cmtRegex, cmtToken, `#`, 1))
	slashLangs = []string{"go", "java", "c", "cpp", "h", "hpp"}
	hashLangs  = []string{"py", "sh", "bash", "yml", "yaml"}
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

	for scan.Scan() {
		line := scan.Text()
		iss, found := parseLine(rexp, line)
		if found {
			log.Debugf("found issue: %s\n%s", line, iss.String())
			issues = append(issues, iss)
		}
	}

	if scan.Err() != nil {
		return issues, errors.Wrapf(scan.Err(), "error while scanning file: %s", fileName)
	}

	return issues, nil
}

func commentRegexes(ext string) *regexp.Regexp {
	idx := sort.SearchStrings(slashLangs, ext)
	fmt.Println(idx, len(slashLangs), slashLangs)
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
	var i tracker.Issue

	finds := rexp.FindStringSubmatch(line)
	if len(finds) <= 0 {
		return i, false
	}

	log.Info(finds)
	i.Title = "temp"

	return i, true
}
