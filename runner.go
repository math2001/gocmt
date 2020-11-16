// Manages the second componenet: run the tests, and collect the results
package main

import (
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/math2001/gocmt/checks"
	"github.com/math2001/gocmt/cmt"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// check name: check function
var allchecks = map[string]checkerfunction{
	"cpu":      checks.CPU,
	"boottime": checks.Boottime,
	"load":     checks.Load,
	"disks":    checks.Disks,
	"folders":  checks.Folders,
}

// This function returns before all the tests have finished running. It returns
// a channel on which the check results are send. The channel is closed as soon
// as all the tests have finished running.
func runChecks(conf Config) <-chan *cmt.Check {

	var wg sync.WaitGroup

	checkresults := make(chan *cmt.Check)

	// producer (produces check results)

	for name, fn := range allchecks {
		if name == "_globals" {
			panic("'_globals' is a reserved name (found check named _globals)")
		}

		if !isCheckEnabled(conf.FrameworkSettings, name) {
			continue
		}

		value, ok := conf.ChecksArguments[name]
		if ok {
			sets, ok := value.([]interface{})
			if !ok {
				panic(fmt.Sprintf("%s: invalid argument sets. It should be a list, got %#v", name, value))
			}

			for _, set := range sets {
				wg.Add(1)
				var argSet map[string]interface{}
				if err := mapstructure.Decode(set, &argSet); err != nil {
					panic(errors.Wrapf(err, "decoding %#v into map[string]interface{}. You probably mess up check_arguments for %s", set, name))
				}
				go runCheck(&wg, name, fn, checkresults, argSet)
			}

		} else {
			wg.Add(1)
			go runCheck(&wg, name, fn, checkresults, nil)
		}

	}

	go func() {
		wg.Wait() // waits for the producer to finish
		close(checkresults)
	}()

	return checkresults
}

func isCheckEnabled(fs *FrameworkSettings, name string) bool {
	// TODO#perf: sort checks names and binary search
	for _, n := range fs.Checks {
		if n == name {
			return true
		}
	}
	return false
}

func runCheck(
	wg *sync.WaitGroup,
	name string,
	fn checkerfunction,
	checkresults chan<- *cmt.Check,
	argset map[string]interface{},
) {
	defer wg.Done()
	checkresult := cmt.NewCheckResult(name, argset)

	defer func() {
		if r := recover(); r != nil {
			checkresult.SetPanic(r, debug.Stack())
		}
		// send to the channel *after* we are done with object
		checkresults <- checkresult
	}()

	fn(checkresult, argset)
}
