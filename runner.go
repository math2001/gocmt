// Manages the second componenet: run the tests, and collect the results
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"sort"
	"sync"

	"github.com/math2001/gocmt/checks"
	"github.com/math2001/gocmt/cmt"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// check name: check function
var allchecks = map[string]checkerFunction{
	"cpu":              checks.CPU,
	"boottime":         checks.Boottime,
	"load":             checks.Load,
	"disks":            checks.Disks,
	"folders":          checks.Folders,
	"network_counters": checks.NetworkCounters,
	"memory":           checks.Memory,
	"process":          checks.Process,
	"swap":             checks.Swap,
	"mounts":           checks.Mounts,
	"pings":            checks.Pings,
	"urls":             checks.URLs,
}

// This function returns before all the tests have finished running. It returns
// a channel on which the check results are send. The channel is closed as soon
// as all the tests have finished running.
func runChecks(checkresults chan<- *cmt.CheckResult, conf Config, sequential bool) {

	var wg sync.WaitGroup

	db := loadDatabaseFromFile(conf.FrameworkSettings.DatabaseFile)

	// always start the tests in the same order (map aren't sorted in Go)
	// consistency can only be good. And it's pretty cheap here
	names := make([]string, len(allchecks))
	i := 0
	for name := range allchecks {
		names[i] = name
		i++
	}
	sort.Slice(names, func(i, j int) bool {
		return names[i] < names[j]
	})

	// producer (produces check results)
	for _, name := range names {
		fn := allchecks[name]
		if name == "_globals" {
			panic("'_globals' is a reserved name (found check named _globals)")
		}

		if !isCheckEnabled(conf.FrameworkSettings, name) {
			continue
		}

		if _, ok := db[name]; !ok {
			db[name] = make(map[string]interface{})
		}

		value, ok := conf.ArgumentSets[name]
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
				if sequential {
					runCheck(&wg, name, fn, checkresults, argSet, db[name])
				} else {
					go runCheck(&wg, name, fn, checkresults, argSet, db[name])
				}
			}

		} else {
			wg.Add(1)
			if sequential {
				runCheck(&wg, name, fn, checkresults, nil, db[name])
			} else {
				go runCheck(&wg, name, fn, checkresults, nil, db[name])
			}
		}

	}

	go func(db map[string]map[string]interface{}, dbfile string) {
		wg.Wait() // waits for the producer to finish
		saveDatabaseToFile(db, dbfile)
		close(checkresults)

	}(db, conf.FrameworkSettings.DatabaseFile)
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
	fn checkerFunction,
	checkresults chan<- *cmt.CheckResult,
	argset map[string]interface{},
	db map[string]interface{},
) {
	defer wg.Done()
	checkresult := cmt.NewCheckResult(name, argset, db)

	defer func() {
		if r := recover(); r != nil {
			checkresult.SetPanic(r, debug.Stack())
		}
		// send to the channel *after* we are done with object
		checkresults <- checkresult
	}()

	fn(checkresult, argset)
}

func loadDatabaseFromFile(filename string) map[string]map[string]interface{} {
	if filename == "" {
		panic("database_file (string) is a required configuration option")
	}
	f, err := os.Open(filename)
	if err != nil {
		return make(map[string]map[string]interface{})
	}
	defer f.Close()

	var db map[string]map[string]interface{}
	if err := json.NewDecoder(f).Decode(&db); err != nil {
		log.Printf("[load db from file]: %s", err)
		return make(map[string]map[string]interface{})
	}
	return db
}

func saveDatabaseToFile(database map[string]map[string]interface{}, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Printf("[save db to file]: %s", err)
		return
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(database); err != nil {
		log.Printf("[save db to file]: %s", err)
		return
	}
}
