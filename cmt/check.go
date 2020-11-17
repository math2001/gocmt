package cmt

import (
	"bytes"
	"log"
)

// CheckResult is basically a getter setter class. I don't think it's the
// typical go way, but it makes writing checks easier to read (AddError and
// AddItem), and that's what matters here.
type CheckResult struct {
	name   string
	argset map[string]interface{}

	errors     []error
	checkitems []*CheckItem
	DB         map[string]interface{}

	debugbuf bytes.Buffer

	isAlert      bool
	alertMessage string

	panicData *panicData
}

func (c *CheckResult) AddError(err error) {
	c.errors = append(c.errors, err)
}

func (c *CheckResult) AddItem(ci *CheckItem) {
	c.checkitems = append(c.checkitems, ci)
}

func (c *CheckResult) Errors() []error {
	return c.errors
}

func (c *CheckResult) CheckItems() []*CheckItem {
	return c.checkitems
}

func (c *CheckResult) Name() string {
	return c.name
}

func (c *CheckResult) ArgumentSet() map[string]interface{} {
	return c.argset
}

func (c *CheckResult) DebugBuffer() *bytes.Buffer {
	return &c.debugbuf
}

func (c *CheckResult) SetPanic(message interface{}, stack []byte) {
	if c.panicData != nil {
		log.Printf("[checkresult] warning: panic already set %q (overwrote by %q)", c.panicData.msg, message)
	}
	c.panicData = &panicData{
		msg:   message,
		stack: stack,
	}
}

func (c *CheckResult) GetPanic() (r interface{}, stack []byte) {
	if c.panicData == nil {
		return nil, nil
	}
	return c.panicData.msg, c.panicData.stack
}

func NewCheckResult(name string, argset map[string]interface{}, db map[string]interface{}) *CheckResult {
	return &CheckResult{
		name:   name,
		argset: argset,
		DB:     db,
	}
}

type CheckItem struct {
	Name        string
	Value       interface{}
	Description string
	Unit        string

	IsAlert      bool
	AlertMessage string
}

type panicData struct {
	msg   interface{} // whatever is passed to panic
	stack []byte
}
