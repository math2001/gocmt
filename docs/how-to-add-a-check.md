# How to add a test

## 1. Create the check file

Add a `testname.go` in the `checks` folder. The filename doesn't matter, only the `.go` is important.

```go
package checks

func CheckName(
    cr *cmt.CheckResult,
    args map[string]interface{},
) {
    // code for the check goes here
}
```

The function's name **has to start with a capital letter**. This way the symbol is available to other packages (namely the main package, where the framework lives).

See [how to write a check](./how-to-write-a-check.md) for CMT's API.

## 2. Add check to map of all checks

The runner (in `runner.go`) needs to know about this new check. Just add it to the map `allchecks` (`checkname => check function`), like so:

```go
// check name: check function
var allchecks = map[string]checkerfunction{
    "cpu": checks.CPU,
    "mem": checks.Mem,
    "checkname": checks.CheckName, // notice the comma
}
```

## 3. Build

```bash
$ go build
```

If you (and I) didn't make any mistake, it should remain quiet and build the
binary.

## 4. Enable the check

Now the check just needs to be enabled in the configuration.

FIXME: the current system to enable/disable checks doesn't allow a child conf
to disable checks that are enabled in a parent conf.

## Name correspondance

The function name doesn't have to be related to the *actual* check name: the key in the `allchecks` map. The actual check name has to be used in the conf.

```yaml
...
checks:
    - actual_name
    ...
```

```go
// check name: check function
var allchecks = map[string]checkerfunction{
    "cpu": checks.CPU,
    "mem": checks.Mem,
    "acutalname": checks.WhateverName, // notice the comma
}
```

```go
func WhateverName(...) { ... }
```