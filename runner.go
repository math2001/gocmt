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
func runChecks(conf map[string]interface{}) <-chan *cmt.CheckResult {
	c := make(map[string]func(map[string]interface{}, map[string]interface{}) *cmt.CheckResult)
	c["cpu"] = checks.CheckCPU

	var wg sync.WaitGroup

	checkresults := make(chan *cmt.CheckResult)

	// producer (produces check results)
	globals := conf["checks_settings"]["_globals"]
	for name, fn := range c {
		// doesn't panic if name isn't a key (not like Python)
		subconf := conf["check_settings"][name]
		wg.Add(1)
		go func() {
			defer wg.Done()
			checkresults <- fn(globals, subconf)
		}()
	}

	go func() {
		wg.Wait() // waits for the producer to finish
		close(checkresults)
	}()

	return checkresults
}
