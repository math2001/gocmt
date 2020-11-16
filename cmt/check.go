package cmt

import (
	"bytes"
	"log"
)

// Check is basically a getter setter class. I don't think it's the
// typical go way, but it makes writing checks easier to read (AddError and
// AddItem), and that's what matters here.
type Check struct {
	name   string
	argset map[string]interface{}

	errors     []error
	checkitems []*CheckItem

	debugbuf bytes.Buffer

	isAlert      bool
	alertMessage string

	panicData *panicData
}

func (c *Check) AddError(err error) {
	c.errors = append(c.errors, err)
}

func (c *Check) AddItem(ci *CheckItem) {
	c.checkitems = append(c.checkitems, ci)
}

func (c *Check) Errors() []error {
	return c.errors
}

func (c *Check) CheckItems() []*CheckItem {
	return c.checkitems
}

func (c *Check) Name() string {
	return c.name
}

func (c *Check) ArgumentSet() map[string]interface{} {
	return c.argset
}

func (c *Check) DebugBuffer() *bytes.Buffer {
	return &c.debugbuf
}

func (c *Check) SetAlert(msg string) {
	if c.isAlert {
		log.Printf("[checkresult] warning: alert already set to %q (overwrote by %q)", c.alertMessage, msg)
	}
	c.isAlert = true
	c.alertMessage = msg
}

func (c *Check) GetAlert() (is_alert bool, message string) {
	return c.isAlert, c.alertMessage
}

func (c *Check) SetPanic(message interface{}, stack []byte) {
	if c.panicData != nil {
		log.Printf("[checkresult] warning: panic already set %q (overwrote by %q)", c.panicData.msg, message)
	}
	c.panicData = &panicData{
		msg:   message,
		stack: stack,
	}
}

func (c *Check) GetPanic() (r interface{}, stack []byte) {
	if c.panicData == nil {
		return nil, nil
	}
	return c.panicData.msg, c.panicData.stack
}

func NewCheckResult(name string, argset map[string]interface{}) *Check {
	return &Check{
		name:   name,
		argset: argset,
	}
}

type CheckItem struct {
	Name        string
	Value       interface{}
	Description string
	Unit        string
}

type panicData struct {
	msg   interface{} // whatever is passed to panic
	stack []byte
}
