package unit

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

var Incomparable = errors.New("unit mismatch")

// Qualified represents a value with associated units.  It is implemented by
// Maker, Value, and Units to allow them to be used anywhere flexibility in
// type is useful.
type Qualified interface {
	Value() float64
	Units() Units
}

// Value is a concrete value with associated units.
type Value struct {
	S float64
	U Units
}

// Compare matches units of a and b, and calls cmpFn to compare the
// scalar values.  Returns the result of cmpFn, or an Incomparable error
// if a and b do not have conforming units.
func (a Value) Compare(b Qualified, cmpFn func(a, b float64) bool) (bool, error) {
	if a.Units().Equal(b.Units()) {
		return cmpFn(a.Value(), b.Value()), nil
	}
	ab, remain := a.Convert(b.Units().Make)
	if !remain.IsEmpty() {
		return false, fmt.Errorf("%w: %q != %q (diff=%q)", Incomparable, a.Units(), b.Units(), remain)
	}
	return cmpFn(ab.Value(), b.Value()), nil
}

// Equal returns true if the units for a and b are equivalent, and the
// scalar values are equal.  Returns an Incomparable error if the units
// are not conformable.
func (a Value) Equal(b Qualified) (bool, error) {
	return a.Compare(b, func(a, b float64) bool { return a == b })
}

// Less returns true if the units for a and b are equivalent, and the
// scalar value of a is less than b.  Returns an Incomparable error if
// the units are not conformable.
func (a Value) Less(b Qualified) (bool, error) {
	return a.Compare(b, func(a, b float64) bool { return a < b })
}

// Value returns the scalar (float64) component of v.
func (v Value) Value() float64 { return v.S }

// Units returns the Units for v.
func (v Value) Units() Units { return v.U }

// MulN multiplies b by the scalar value of a, returning the result,
// and keeping the units from a.
func (a Value) MulN(b float64) (r Value) {
	r.S = a.S * b
	r.U = a.U
	return
}

// Mul multiplies a and b, returning the result.
func (a Value) Mul(b Qualified) (r Value) {
	r.S = a.S * b.Value()
	r.U = a.U.Mul(b.Units())
	return
}

// DivN divides the scalar value of a by b, returning the result, and
// keeping the units from a.
func (a Value) DivN(b float64) (r Value) {
	r.S = a.S / b
	r.U = a.U
	return
}

// Div divides a by b, returning the result.
func (a Value) Div(b Qualified) (r Value) {
	r.S = a.S / b.Value()
	r.U = a.U.Div(b.Units())
	return
}

// AddN adds b to the scalar value of a, returning the result, and keeping the units from a.
func (a Value) AddN(b float64) (r Value) {
	r.S = a.S + b
	r.U = a.U
	return
}

func (a Value) conform(b Value) (r Value, ok bool) {
	if !a.U.Equal(b.Units()) {
		var remain Units
		b, remain = b.Convert(a.U.Make)
		if !remain.IsEmpty() {
			return
		}
	}
	return b, true
}

// Add adds a and b, returning the result.  If the units are not
// equivalent, r will be empty and ok will be false.
func (a Value) Add(b Value) (r Value, ok bool) {
	if b, ok = a.conform(b); !ok {
		return
	}
	r.S = a.S + b.Value()
	r.U = a.U
	ok = true
	return
}

// SubN subtracts b from the scalar value of a, returning the result, and keeping the units from a.
func (a Value) SubN(b float64) (r Value) {
	r.S = a.S - b
	r.U = a.U
	return
}

// Sub subtracts b from a, returning the result.  If the units are not
// equivalent, r will be empty and ok will be false.
func (a Value) Sub(b Value) (r Value, ok bool) {
	if b, ok = a.conform(b); !ok {
		return
	}
	r.S = a.S - b.Value()
	r.U = a.U
	ok = true
	return
}

// Pow raises a to the power of n, which may be negative, returning the result.
func (a Value) Pow(n int) (r Value) {
	r.S = math.Pow(a.S, float64(n))
	r.U = a.U.Pow(n)
	return
}

// Convert forces a into the units of wanted.  Returns the resulting value
// and any "extra" units that indicate non-conformability between a and wanted.
// For instance, "5 m/s".Convert("m") will return ("5 m", "1/s").
func (a Value) Convert(wanted Maker) (result Value, remainder Units) {
	if a.Units().Equal(wanted.Units()) {
		result = a
		return
	}
	defer tracein("%q.Convert(%q)", a, wanted)()
	rv := a.Div(wanted).Units().Reduce() // Divide out the wanted units
	tracemsg("result=%q remain=%q", wanted(a.S*rv.S), rv.U)
	return wanted(a.S * rv.S), rv.U
}

func MustConvert(a Qualified, wanted Maker) Value {
	v := FromQualified(a)
	result, remain := v.Convert(wanted)
	if !remain.IsEmpty() {
		panic(fmt.Sprintf("Cannot convert %q to units %q (got: %q, remainder: %q)", a, wanted, result, remain))
	}
	return result
}

// Reduce reduces the units for a to primitive units and returns the resulting
// value.
func (a Value) Reduce() (v Value) {
	defer tracein("%q.Reduce()", a)()
	v = a.U.Reduce()
	v.S *= a.S
	return
}

func formatStr(f fmt.State, c rune) string {
	var sb strings.Builder
	for _, c := range "+-# 0" {
		if f.Flag(int(c)) {
			fmt.Fprint(&sb, string(c))
		}
	}
	if wid, ok := f.Width(); ok {
		fmt.Fprint(&sb, wid)
	}
	if prec, ok := f.Precision(); ok {
		fmt.Fprint(&sb, ".", prec)
	}
	fmt.Fprint(&sb, string(c))
	return sb.String()
}

// Format implements fmt.Printf-style verbs to format the value.
// The following verbs are supported:
//
//    %f %g     render only the scalar portion of the value
//    %s %q %v  render both the scalar and units portion of the value as a
//              string, using the %g format to render the scalar portion,
//              and adding a space in between scalar and units, if present.
//
// For more precise control over the output, you can access the S and
// U fields of the type directly.
func (a Value) Format(f fmt.State, c rune) {
	switch c {
	case 'v', 's', 'q':
		/*
			var sb strings.Builder
			fmt.Fprintf(&sb, "%g", a.S)
			if !a.U.IsEmpty() {
				fmt.Fprintf(&sb, " %s", a.U.String())
			}
			s := sb.String()
		*/
		s := DefaultFormatter.Format(a)
		fmt.Fprintf(f, "%"+formatStr(f, c), s)
	default:
		fmt.Fprintf(f, "%"+formatStr(f, c), a.S)
	}
}

func (a Value) String() string {
	return fmt.Sprintf("%s", a)
}

func (a Value) Recip() Value {
	var r Value
	r.S = 1 / a.S
	r.U = a.U.Recip()
	return r
}
