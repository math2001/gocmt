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
	DB         map[string]interface{}

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

func NewCheck(name string, argset map[string]interface{}, db map[string]interface{}) *Check {
	return &Check{
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
