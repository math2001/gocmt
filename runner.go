// Manages the second componenet: run the tests, and collect the results
package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"sync"

	"github.com/math2001/gocmt/checks"
	"github.com/math2001/gocmt/cmt"
)

// check name: check function
var allchecks = map[string]checkerfunction{
	"cpu":      checks.CPU,
	"boottime": checks.Boottime,
	"load":     checks.Load,
}

// This function returns before all the tests have finished running. It returns
// a channel on which the check results are send. The channel is closed as soon
// as all the tests have finished running.
func runChecks(conf cmt.Conf) <-chan *cmt.CheckResult {

	var wg sync.WaitGroup

	checkresults := make(chan *cmt.CheckResult)

	// producer (produces check results)
	var globals map[string]interface{}
	if conf.CheckSettings["_globals"] != nil {
		globals = conf.CheckSettings["_globals"].(map[string]interface{})
	}

	for name, fn := range allchecks {
		if name == "_globals" {
			panic("'_globals' is a reserved name (found check named _globals)")
		}

		if !isCheckEnabled(conf.FrameworkSettings, name) {
			continue
		}

		// doesn't panic if name isn't a key (not like Python)
		var subconf map[string]interface{}
		if conf.CheckSettings[name] != nil {
			subconf = conf.CheckSettings[name].(map[string]interface{})
		}

		// important, because fn and name are used in the goroutine below
		fn := fn
		name := name

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					fmt.Fprintf(os.Stderr, "panic: %s\n", r)
					debug.PrintStack()
				}
				// TODO: report panic properly
				// checkresult.SetPanic(r, debug.Stack())
			}()

			checkresult := cmt.NewCheckResult(name)
			fn(checkresult, globals, subconf)
			checkresults <- checkresult
		}()
	}

	go func() {
		wg.Wait() // waits for the producer to finish
		close(checkresults)
	}()

	return checkresults
}

func isCheckEnabled(fs *cmt.FrameworkSettings, name string) bool {
	// TODO#perf: sort checks names and binary search
	for _, n := range fs.Checks {
		if n == name {
			return true
		}
	}
	return false
}
