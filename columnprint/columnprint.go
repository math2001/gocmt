package columnprint

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// U simplifies column printing
type U struct {
	// cols stores the format string for each column
	cols        []string
	verbsPerCol []int
	widths      []int

	record    []interface{}
	recording bool
}

// SetColumns sets the column format strings (one argument = one column).
// There can be more than one verb per column.
func (u *U) SetColumns(cols ...string) {
	u.cols = cols
	u.verbsPerCol = make([]int, len(u.cols))
	for i, f := range cols {
		u.verbsPerCol[i] = countVerbs(f)
	}
	u.widths = make([]int, len(u.cols))
}

// WouldPrint figures out the minimum width required for each column
func (u *U) WouldPrint(vs ...interface{}) error {
	var head int
	for i := 0; i < len(u.cols); i++ {
		width, err := fmt.Fprintf(ioutil.Discard, u.cols[i], vs[head:head+u.verbsPerCol[i]]...)
		if err != nil {
			return err
		}
		if width > u.widths[i] {
			u.widths[i] = width
		}
		head += u.verbsPerCol[i]
	}
	if u.recording {
		u.record = append(u.record, vs)
	}
	return nil
}

// WouldPrintLiteral makes sure the columns are wide enough to fit the literal strings
func (u *U) WouldPrintLiteral(s ...string) {
	if len(s) != len(u.cols) {
		panic("invalid usage")
	}
	for i, v := range s {
		if len(v) > u.widths[i] {
			u.widths[i] = len(v)
		}
	}
	if u.recording {
		// // silly Go, I need to put it in a []interface{} buf
		// buf := make([]interface{}, len(s))
		// for i := 0; i < len(buf); i++ {
		// 	buf[i] = s[i]
		// }
		u.record = append(u.record, s)
	}
}

// Fprint prints out the formated strings in columns, padded with spaces to the
// right of the text
func (u *U) Fprint(w io.Writer, vs ...interface{}) error {
	var head int
	for i := 0; i < len(u.cols); i++ {
		width, err := fmt.Fprintf(w, u.cols[i], vs[head:head+u.verbsPerCol[i]]...)
		if err != nil {
			return err
		}
		for j := 0; j < u.widths[i]-width; j++ {
			_, err := fmt.Fprintf(w, " ")
			if err != nil {
				return err
			}
		}
		head += u.verbsPerCol[i]

		var char string
		if i == len(u.cols)-1 {
			char = "\n"
		} else {
			char = " "
		}

		_, err = fmt.Fprintf(w, char)
		if err != nil {
			return err
		}
	}
	return nil
}

// FprintLiteral prints out the string literals in columns, padded with spaces
// to the right of the text
func (u *U) FprintLiteral(w io.Writer, s ...string) error {
	if len(s) != len(u.cols) {
		panic("invalid usage")
	}
	for i, v := range s {
		space := " "
		if i == len(s)-1 {
			space = ""
		}
		_, err := fmt.Fprintf(w, "%*s%s", -u.widths[i], v, space)
		if err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(w, "\n")
	return err
}

// Print calls Fprint with os.Stdout
func (u *U) Print(vs ...interface{}) error {
	return u.Fprint(os.Stdout, vs...)
}

// PrintLiteral calls FprintLiteral with os.Stdout
func (u *U) PrintLiteral(s ...string) error {
	return u.FprintLiteral(os.Stdout, s...)
}

// Record records the argument given to would print so that FprintFromRecord can
// act as the second Print loop
func (u *U) Record(numRowsGuess int) {
	u.record = make([]interface{}, 0, numRowsGuess)
	u.recording = true
}

// FprintFromRecord acts as second Print loop, printing from the record
func (u *U) FprintFromRecord(w io.Writer) error {
	if !u.recording {
		panic("invalid usage")
	}
	u.recording = false
	for _, args := range u.record {
		var err error
		if s, ok := args.([]string); ok {
			err = u.FprintLiteral(w, s...)
		} else if vs, ok := args.([]interface{}); ok {
			err = u.Fprint(w, vs...)
		} else {
			panic("internal error")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// PrintFromRecord calls FprintFromRecord with os.Stdout
func (u *U) PrintFromRecord() error {
	return u.FprintFromRecord(os.Stdout)
}

func countVerbs(format string) (n int) {
	for i := 0; i < len(format); i++ {
		if format[i] == '%' {
			if i < len(format)-1 && format[i+1] == '%' {
				i++
			} else {
				n++
			}
		}
	}
	return
}
