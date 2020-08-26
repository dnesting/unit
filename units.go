package unit

import (
	"fmt"
	"sort"
	"strings"
)

type Units struct {
	N []Unit
	D []Unit
}

func (us Units) Singular() Unit {
	if len(us.N) == 1 && len(us.D) == 0 {
		return us.N[0]
	}
	return nil
}

func (us Units) IsEmpty() bool {
	return len(us.N) == 0 && len(us.D) == 0
}

// Recip returns the reciprocal of units, swapping numerator and denominator.
func (a Units) Recip() Units {
	var r Units
	r.N = append(r.N, a.D...)
	r.D = append(r.D, a.N...)
	return r
}

func (a Units) Value() float64 { return 1 }
func (a Units) Units() Units   { return a }

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
		if u.Equal(last) {
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

// Mul returns the multiplication of the two units, effectively creating
// Units{N: a.N+b.N, D: a.D+b.D}.
func (a Units) Mul(b Units) Units {
	var r Units
	r.N = append(r.N, a.N...)
	r.N = append(r.N, b.N...)
	r.D = append(r.D, a.D...)
	r.D = append(r.D, b.D...)
	r.cancel()
	return r
}

// Div returns the division of the two units, equivalent to a.Mul(b.Recip()).
func (a Units) Div(b Units) Units {
	return a.Mul(b.Recip())
}

// Pow returns a raised to the power of p, equivalent to multiplying a by
// its original value p-1 times.  If p is 0, an empty Units will be returned.
func (a Units) Pow(p int) Units {
	// XXX: The type of p is int, not a float. This means it's not possible to
	// reverse the effect of this function to move, say, from m^2/s^2 to m/s by
	// passing in 0.5.  This opens the door for non-integer exponents for units,
	// necessitating some redesign.
	if p < 0 {
		return a.Recip().Pow(-p).Recip()
	}
	if p == 0 {
		return Units{}
	}
	r := a
	for p-1 > 0 {
		r = r.Mul(a)
		p--
	}
	return r
}

// Equal returns true if a and b contain the same units.
func (a Units) Equal(b Units) bool {
	return unitList(a.N).Equal(b.N) && unitList(a.D).Equal(b.D)
}

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
			// neither symbol compares less than the other, but consult the Unit
			// implementation of Equal just to be sure.
			if !a.N[nr].Equal(a.D[dr]) {
				// equal symbols but Equal() returned false, so just keep both
				// XXX this is potentally buggy if we sort two "foo" units according
				// to symbol.  If we're going to have a Unit.Equal, that probably
				// means we need a Unit.Less.  Ugh.
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
			tracemsg("- %v.Value() = %v", left[i], left[i].Value())
			val := left[i].Value().Reduce()
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

// Reduce reduce us to primitive units.  The return type is a Value since
// there may be a scalar multiplier that needs to be applied to the final
// result.
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

func (us Units) Make(v float64) Value {
	return Value{
		S: v,
		U: us,
	}
}

func (us Units) format(unitSep, fracSep string, expFn func(s string, exp int) string) string {
	if us.IsEmpty() {
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
