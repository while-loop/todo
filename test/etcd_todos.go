package main

import (
	"github.com/while-loop/todo/pkg/log"
	"github.com/while-loop/todo/pkg/parser"
	"github.com/while-loop/todo/pkg/tracker"
	"os"
	"reflect"
)

var fileName = `test/etcdTodos.test`

func main() {
	// this is a list of TODOs found in the root directory of https://github.com/coreos/etcd
	// using the command `grep -nri -E "^.*//\s*todo.*" > etcdTodos.txt`
	// empty todos are not treated as issues

	fail := false

	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	issues, err := parser.ParseFile(fileName, f)
	if err != nil {
		panic(err)
	}

	if len(issues) != 441 {
		log.Errorf("wrong issue len. want %v, got %v", 441, len(issues))
		fail = true
	}

	for _, ei := range expIssues {
		if !reflect.DeepEqual(ei, issues[ei.Line-1]) {
			log.Errorf("issue not equal\nwant:\n%v\ngot:\n%v", ei, issues[ei.Line-1])
			fail = true
		}
	}

	if fail {
		os.Exit(2)
	} else {
		log.Info("core os test passed")
	}
}

// expIssues has a select of todos given in the file
// more will be added as time goes on
var expIssues = []tracker.Issue{
	{Title: "save_fsync latency?", File: fileName, Line: 1},
	{Title: "ensure the entries are continuous and", Assignee: "xiangli", File: fileName, Line: 25},
	{Title: "The original rationale of passing in a pre-allocated", Assignee: "beorn7", File: fileName, Line: 133},
	{Title: "consider a more generally-known optimization for reflect.Value ==> Interface", File: fileName, Line: 200},
	{Title: "support bincUnicodeOther (for now, just use string or bytearray)", File: fileName, Line: 201},
}
