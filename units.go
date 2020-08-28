package unit

import (
	"sort"
)

// Units represent a ratio of unit lists that are attached to qualified values.
// Units and its contents may be re-used internally and must not be directly
// modified (including its numerator and denominator slices).  Units implements
// Qualified, representing itself as a qualified value of 1.
type Units struct {
	N []Unit
	D []Unit
}

/*
func (us Units) Singular() Unit {
	if len(us.N) == 1 && len(us.D) == 0 {
		return us.N[0]
	}
	return nil
}
*/

// Empty returns true if this Units has no units in the numerator or denominator.
func (us Units) Empty() bool {
	return len(us.N) == 0 && len(us.D) == 0
}

// Recip returns the reciprocal of Units, swapping numerator and denominator.
func (a Units) Recip() Units {
	var r Units
	r.N = append(r.N, a.D...)
	r.D = append(r.D, a.N...)
	return r
}

// Value returns 1.
func (a Units) Value() float64 { return 1 }

// Units returns itself.
func (a Units) Units() Units { return a }

// unitList represents a list of Units.  It implements sort.Interface to sort
// by symbol name (sorting nil values last).
type unitList []Unit

func (ul unitList) Len() int      { return len(ul) }
func (ul unitList) Swap(a, b int) { ul[a], ul[b] = ul[b], ul[a] }
func (ul unitList) Less(a, b int) bool {
	if ul[a] == nil {
		return false
	}
	if ul[b] == nil {
		return true
	}
	return ul[a].Symbol() < ul[b].Symbol()
}

// Equal returns true if both unit lists contain equal units.  This method
// assumes the unitList is sorted.
func (a unitList) Equal(b unitList) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if !a[i].Equal(b[i]) {
			return false
		}
	}
	return true
}

/*
func (ul unitList) Format(sep string, expFn func(s string, n int) string) string {
	var r []string
	var exp int
	var last Unit
	if expFn == nil {
		expFn = defaultExp
	}
	emit := func() {
		if exp > 0 {
			r = append(r, expFn(last.Symbol(), exp+1))
		} else {
			r = append(r, last.Symbol())
		}
	}
	for _, u := range ul {
		if u == nil {
			continue
		}
		if Equal(u, last) {
			exp++
		} else {
			if last != nil {
				emit()
			}
			last = u
			exp = 0
		}
	}
	if last != nil {
		emit()
	}
	return strings.Join(r, sep)
}

func defaultExp(s string, n int) string {
	return fmt.Sprintf("%s^%d", s, n)
}

func (ul unitList) String() string {
	return ul.Format(" ", defaultExp)
}
*/

// Mul returns the multiplication of the two units, effectively creating
// Units{N: a.N+b.N, D: a.D+b.D}.
func (a Units) Mul(b Units) Units {
	r := a.mul(b)
	r.cancel()
	return r
}

func (a Units) mul(b Units) Units {
	var r Units
	r.N = append(r.N, a.N...)
	r.N = append(r.N, b.N...)
	r.D = append(r.D, a.D...)
	r.D = append(r.D, b.D...)
	return r
}

// Div returns the division of the two units, equivalent to a.Mul(b.Recip()).
func (a Units) Div(b Units) Units {
	return a.Mul(b.Recip())
}

// Pow returns a raised to the power of p, equivalent to multiplying a by
// its original value p-1 times.  If p is 0, an empty Units will be returned.
func (a Units) Pow(p int) Units {
	// The type of p is int, not a float. This means it's not possible to
	// reverse the effect of this function to move, say, from m^2/s^2 to m/s by
	// passing in 0.5.
	if p < 0 {
		return a.Recip().Pow(-p).Recip()
	}
	if p == 0 {
		return Units{}
	}
	r := a
	for p-1 > 0 {
		r = r.mul(a)
		p--
	}
	r.cancel()
	return r
}

// Equal returns true if a and b contain the same units.
func (a Units) Equal(b Units) bool {
	return unitList(a.N).Equal(b.N) && unitList(a.D).Equal(b.D)
}

// Equiv returns true if a and b reduce to the same primitive units.
func (a Units) Equiv(b Units) bool {
	if a.Equal(b) {
		return true
	}
	ar := a.Reduce()
	br := b.Reduce()
	return ar.U.Equal(br.U)
}

// Cancel identifies units in both the numerator and denominator and
// removes them.  Units must be exact matches; no reduction is performed.
func (a *Units) cancel() {
	var nr, dr int // read index to r.N and a.D
	var nw, dw int // write index, always <= the read index

	// First sort by symbol, with nil values at the end
	sort.Sort(unitList(a.N))
	sort.Sort(unitList(a.D))

	for nr < len(a.N) && dr < len(a.D) {
		if a.N[nr] == nil {
			// skip nil values in the numerator
			nr++
			continue
		}
		if a.D[dr] == nil {
			// skip nil values in the denominator
			dr++
			continue
		}
		if a.N[nr].Symbol() < a.D[dr].Symbol() {
			// if the denominator symbol is after the numerator symbol, that
			// means the numerator has nothing to cancel it, so we advance the
			// numerator.
			if nr != nw {
				a.N[nw] = a.N[nr]
			}
			nr++
			nw++
		} else if a.D[dr].Symbol() < a.N[nr].Symbol() {
			// if the numerator symbol is after the denominator symbol, that
			// means the denominator has nothing to cancel it, so we advance the
			// denominator.
			if dr != dw {
				a.D[dw] = a.D[dr]
			}
			dr++
			dw++
		} else {
			// TODO: We need to scan forward for all equal symbols to find one
			// for which Equal returns true, since we could have multiple units that
			// have the same symbol but different definitions.
			if !a.N[nr].Equal(a.D[dr]) {
				if nr != nw {
					a.N[nw] = a.N[nr]
				}
				nw++
				if dr != dw {
					a.D[dw] = a.D[dr]
				}
				dw++
			}
			nr++
			dr++
		}
	}
	// anything left at the end of a.N or a.D should be kept
	for nr < len(a.N) {
		if a.N[nr] != nil {
			a.N[nw] = a.N[nr]
			nw++
		}
		nr++
	}
	for dr < len(a.D) {
		if a.D[dr] != nil {
			a.D[dw] = a.D[dr]
			dw++
		}
		dr++
	}
	a.N = a.N[:nw]
	a.D = a.D[:dw]
}

func reduceLine(left []Unit, start int) (mult float64, res, right []Unit) {
	defer tracein("reduceLine(%v, %d)", left, start)()
	mult = 1
	for i := start; i < len(left); i++ {
		tracemsg("item %d=%v primitive? %v", i, left[i], IsPrimitive(left[i]))
		if !IsPrimitive(left[i]) {
			tracemsg("- %v.Deriv() = %v", left[i], left[i].Deriv())
			val := left[i].Deriv().Reduce()
			mult *= val.Value()
			us := val.Units()
			left = append(left, us.N...)
			right = append(right, us.D...)
			tracemsg("- replacing %v with %v", left[i], us)
			left[i] = nil
		}
	}
	tracemsg("= mult=%v left=%v right=%v", mult, left, right)
	return mult, left, right
}

// Reduce reduces us to primitive units.  The return type is a Value since
// the act of reducing may introduce a multiplier.
func (us Units) Reduce() (r Value) {
	defer tracein("%q.Reduce()", us)()
	var u Units
	u.N = append([]Unit{}, us.N...)
	u.D = append([]Unit{}, us.D...)

	var mult float64 = 1
	var lastn, lastd int
	var recip []Unit
	var n float64
	for {
		tracemsg("starting pass: %q", u)
		n, u.N, recip = reduceLine(u.N, lastn)
		mult *= n
		lastn = len(u.N)
		u.D = append(u.D, recip...)

		n, u.D, recip = reduceLine(u.D, lastd)
		mult /= n
		lastd = len(u.D)
		if recip == nil {
			break
		}
		u.N = append(u.N, recip...)
	}
	tracemsg("done: %q", u)
	u.cancel()
	tracemsg("= %q", Value{S: mult, U: u})
	return Value{S: mult, U: u}
}

// Make creates a new qualified value.
func (us Units) Make(v float64) Value {
	return Value{
		S: v,
		U: us,
	}
}

/*
func (us Units) format(unitSep, fracSep string, expFn func(s string, exp int) string) string {
	if us.Empty() {
		return ""
	}
	var sb strings.Builder
	n := unitList(us.N).Format(unitSep, expFn)
	if len(n) > 0 {
		sb.WriteString(n)
	} else {
		sb.WriteString("1")
	}
	if len(us.D) > 0 {
		sb.WriteString(fracSep)
		sb.WriteString(unitList(us.D).Format(unitSep, expFn))
	}
	return sb.String()
}

func (us Units) String() string {
	return us.format(" ", "/", nil)
}
*/

func (us Units) String() string {
	return DefaultFormatter.FormatUnits(us)
}
