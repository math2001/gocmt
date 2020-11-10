# How to write a check

## Check signature

```go
func Check(
    cr      *cmt.CheckResult,
    args    map[string]interface{},
)
```

Arguments:

1. `cr *cmt.CheckResult`: this is where the results are stored. You'll be calling `.AddItem`, `.AddError` and `.SaveToDatabase`. CMT will then manage everything based on this report.
3. `args map[string]interface{}`: the arguments for this check. Different argument set can be provided in the conf, 

Return value: Nothing.

If you want to early return (after an error for example) just write `return`.

## Runtime casting

A `map[string]interface{}` is like a Python dictionary, where keys are strings, and the values can be anything (because everything satisfies the empty interface).

If you know that `globals["verbosity"]` will be an int, just doing:

```go
if globals["verbosity"] > 2 {
    // doesn't work
}
```

will not work because Go needs you to explicitely say "it will be an integer", like so:

```go
if globals["verbosity"].(int) > 2 {
    // works
}
```

If the casting fails (verbosity isn't a integer), then Go will panic. This means that the goroutine in which the check is running will stop, but *it will not have any impact on the other checks*.

Maybe you'll find the following code snippet useful:

```go
integer_value, ok := globals["verbosity"].(int)
// ok is set to true if casting succeeded, and false otherwise (never panics).
```

## Example

TODO. For now, have a look at the existing check, they should be quite readable.

## API

Go can automatically generate the documentation for the API from the code.

```bash
$ godoc -http=:6060
```

Open `http://localhost:6060/pkg/github.com/math2001/gocmt/cmt/` in the browser.

(the path might change in the future, when we move the github repo. In truth the actual github repo doesn't matter, only the string in `go.mod`)