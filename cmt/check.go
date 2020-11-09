package cmt

type CheckResult struct {
	errors     []error
	checkitems []*CheckItem
}

func (cr *CheckResult) AddError(err error) {
	cr.errors = append(cr.errors, err)
}

func (cr *CheckResult) AddItem(ci *CheckItem) {
	cr.checkitems = append(cr.checkitems, ci)
}

type CheckItem struct {
	Name        string
	Value       interface{}
	Description string
	Unit        string
}
