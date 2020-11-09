package cmt

import "bytes"

// CheckResult is basically a getter setter class. I don't think it's the
// typical go way, but it makes writing checks easier to read (AddError and
// AddItem), and that's what matters here.
type CheckResult struct {
	name string

	errors     []error
	checkitems []*CheckItem

	debugbuf bytes.Buffer
}

func (cr *CheckResult) AddError(err error) {
	cr.errors = append(cr.errors, err)
}

func (cr *CheckResult) AddItem(ci *CheckItem) {
	cr.checkitems = append(cr.checkitems, ci)
}

func (cr *CheckResult) Errors() []error {
	return cr.errors
}

func (cr *CheckResult) CheckItems() []*CheckItem {
	return cr.checkitems
}

func (cr *CheckResult) Name() string {
	return cr.name
}

func (cr *CheckResult) DebugBuffer() *bytes.Buffer {
	return &cr.debugbuf
}

func NewCheckResult(name string) *CheckResult {
	return &CheckResult{
		name: name,
	}
}

type CheckItem struct {
	Name        string
	Value       interface{}
	Description string
	Unit        string
}
