# How to add a test

## 1. Create the check file

Add a `testname.go` in the `checks` folder. The filename doesn't matter, only the `.go` is important.

```go
package checks

func CheckName(
    cr *cmt.CheckResult,
    globals map[string]interface{},
    settings map[string]interface{},
) {
    // code for the check goes here
}
```

## 2. Add check to map of all checks

The runner (in `runner.go`) needs to know about this new check. Just add it to the map `allchecks` (`checkname => check function`), like so:

```go
// check name: check function
var allchecks = map[string]checkerfunction{
    "cpu": checks.CPU,
    "mem": checks.Mem,
    "checkname": checkname, // notice the comma
}
```

## 3. Enable the check

Now the check just needs to be enabled in the configuration.

FIXME: the current system to enable/disable checks doesn't allow a child conf
to disable checks that are enabled in a parent conf.
