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

	debugbuf bytes.Buffer

	isAlert      bool
	alertMessage string

	panicData *panicData
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

func (cr *CheckResult) ArgumentSet() map[string]interface{} {
	return cr.argset
}

func (cr *CheckResult) DebugBuffer() *bytes.Buffer {
	return &cr.debugbuf
}

func (cr *CheckResult) SetAlert(msg string) {
	if cr.isAlert {
		log.Printf("[checkresult] warning: alert already set to %q (overwrote by %q)", cr.alertMessage, msg)
	}
	cr.isAlert = true
	cr.alertMessage = msg
}

func (cr *CheckResult) GetAlert() (is_alert bool, message string) {
	return cr.isAlert, cr.alertMessage
}

func (cr *CheckResult) SetPanic(message interface{}, stack []byte) {
	if cr.panicData != nil {
		log.Printf("[checkresult] warning: panic already set %q (overwrote by %q)", cr.panicData.msg, message)
	}
	cr.panicData = &panicData{
		msg:   message,
		stack: stack,
	}
}

func (cr *CheckResult) GetPanic() (r interface{}, stack []byte) {
	if cr.panicData == nil {
		return nil, nil
	}
	return cr.panicData.msg, cr.panicData.stack
}

func NewCheckResult(name string, argset map[string]interface{}) *CheckResult {
	return &CheckResult{
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
