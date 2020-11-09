// Manages the second componenet: run the tests, and collect the results
package main

import (
	"sync"

	"github.com/math2001/gocmt/checks"
	"github.com/math2001/gocmt/cmt"
)

// This function returns before all the tests have finished running. It returns
// a channel on which the check results are send. The channel is closed as soon
// as all the tests have finished running.
func runChecks(conf cmt.Conf) <-chan *cmt.CheckResult {

	c := make(map[string]checkerfunction)

	c["cpu"] = checks.CheckCPU

	var wg sync.WaitGroup

	checkresults := make(chan *cmt.CheckResult)

	// producer (produces check results)
	var globals map[string]interface{}
	if conf.CheckSettings["_globals"] != nil {
		globals = conf.CheckSettings["_globals"].(map[string]interface{})
	}

	for name, fn := range c {
		if name == "_globals" {
			panic("'_globals' is a reserved name (found check named _globals)")
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
