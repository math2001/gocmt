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
}

// SetColumns sets the column format verbs (one argument = one column).
// There can be more than one verb per column
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
	return nil
}

// Fprint prints out the columns, padded with spaces (spaces are to the right of the text)
func (u *U) Fprint(w io.Writer, vs ...interface{}) error {
	var head int
	for i := 0; i < len(u.cols); i++ {
		width, err := fmt.Fprintf(w, u.cols[i], vs[head:head+u.verbsPerCol[i]]...)
		if err != nil {
			return err
		}
		for j := 0; j < u.widths[i]-width; j++ {
			fmt.Fprintf(w, " ")
		}
		head += u.verbsPerCol[i]
		if i == len(u.cols)-1 {
			fmt.Fprintf(w, "\n")
		} else {
			fmt.Fprintf(w, " ")
		}
	}
	return nil
}

// Print calls Fprint with os.Stdout
func (u *U) Print(vs ...interface{}) {
	u.Fprint(os.Stdout, vs...)
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
