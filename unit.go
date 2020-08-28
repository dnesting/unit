package unit

import "fmt"

// Unit is a type that defines a single named unit that can then be attached with
// other units to a scalar to form a qualified value.
type Unit interface {
	Qualified
	// Symbol is the name attached to a unit.
	Symbol() string
	// Deriv describes the qualified value that this unit is derived from.
	// For primitive units, this will return a unitless value of 1.
	Deriv() Value
	// Make creates a new qualified value with this unit in its numerator.
	Make(v float64) Value
	// Equal compares units.  Units are equal if they have the same symbol and
	// their derived values are equal.
	Equal(u Unit) bool
}

type unitType struct {
	symbol string
	deriv  Value
	units  Units
}

func (u unitType) String() string       { return fmt.Sprintf("Unit(%q = %v)", u.Symbol(), u.Deriv()) }
func (u unitType) Symbol() string       { return u.symbol }
func (u unitType) Deriv() Value         { return u.deriv }
func (u unitType) Units() Units         { return u.units }
func (u unitType) Value() float64       { return 1 }
func (u unitType) Make(v float64) Value { return u.units.Make(v) }
func (u unitType) Equal(o Unit) bool    { return unitEqual(u, o) }

func unitEqual(a, b Unit) bool {
	if a == nil {
		return b == nil
	}
	if b == nil {
		return false
	}
	if a.Symbol() != b.Symbol() {
		return false
	}
	return a.Deriv().Equal(b.Deriv())
}

// Scalar returns a Maker that produces unitless values multiplied by v.
func Scalar(v float64) Maker { return Value{S: v}.MulN }

var (
	Unity = Scalar(1)
)

// Maker is a function that produces qualified values.  It implements
// Qualified (producing a value of 1).  Makers can be composed together
// using Div, Mul, and Pow to make qualified values representing a
// composite of different units.
type Maker func(v float64) Value

// Value returns the scalar multiplier used by this Maker.
func (m Maker) Value() float64 { return m(1).Value() }

// Units returns the Units that this Maker attaches to values.
func (m Maker) Units() Units { return m(1).Units() }

// Div returns a Maker that divides by b.
func (m Maker) Div(b Maker) Maker {
	return func(f float64) Value { return m(f).Div(b) }
}

// Mul returns a Maker that multiplies by b.
func (m Maker) Mul(b Maker) Maker {
	return func(f float64) Value { return m(f).Mul(b) }
}

// Pow returns a Maker that raises itself by the power of n.
func (m Maker) Pow(n int) Maker {
	return func(f float64) Value { return m(f).Mul(m(1).Pow(n - 1)) }
}

func (m Maker) String() string {
	v := m(1)
	f := v.Value()
	if f != 1 || v.U.Empty() {
		return fmt.Sprintf("*%g", f)
	}
	return v.U.String()
}

// Unit returns the singular unit that this Maker attaches to
// values.  If the Maker attaches a composite of Units, this will
// return nil.
func (m Maker) Unit() Unit {
	us := m.Units()
	if len(us.N) == 1 && len(us.D) == 0 {
		return us.N[0]
	}
	return nil
}

// FromQualified creates a Value from an implementation of Qualified.
func FromQualified(qual Qualified) Value {
	switch v := qual.(type) {
	case Value:
		return v
	case *Value:
		return *v
	default:
		return Value{
			S: qual.Value(),
			U: qual.Units(),
		}
	}
}

func newUnit(symbol string, value Qualified) *unitType {
	u := &unitType{
		symbol: symbol,
		deriv:  FromQualified(value),
	}
	// We populate u.units once so we don't have to re-allocate it
	// for every value derived from the unit.
	u.units = Units{N: []Unit{u}}
	return u
}

// Primitive creates a Maker associated with a new unit named symbol.
// Primitive units are considered irreducible and in unit systems should
// be either the base units of the system, or unitless fundamental
// constants.
func Primitive(symbol string) Maker {
	return Derive(symbol, Unity)
}

// IsPrimitive returns true if u is non-nil and is derived from a unitless
// scalar value of 1.
func IsPrimitive(u Unit) bool {
	if u == nil {
		return false
	}
	v := u.Deriv()
	return v.S == 1 && v.U.Empty()
}

// Derive creates a Maker associated with a new unit named symbol and
// derived from value.
func Derive(symbol string, value Qualified) Maker {
	return newUnit(symbol, FromQualified(value)).Make
}

func Must(v Value, e error) Value {
	if e != nil {
		panic(e)
	}
	return v
}

type prefixType struct {
	prefix string
	mult   float64
	inner  Unit
	units  Units
}

func (u prefixType) Symbol() string    { return u.prefix + u.inner.Symbol() }
func (u prefixType) Deriv() Value      { return u.inner.Make(u.mult) }
func (u prefixType) Equal(o Unit) bool { return unitEqual(u, o) }
func (u prefixType) Units() Units      { return u.units }
func (u prefixType) Value() float64    { return 1 }
func (u prefixType) String() string {
	return fmt.Sprintf("Prefix(%q = %g*%v)", u.prefix, u.mult, u.inner)
}
func (u prefixType) Make(v float64) Value {
	val := u.inner.Make(v)
	val.U = u.units
	return val
}

// Prefix creates a prefix with the given multiplier.  Prefixes can
// wrap other Maker instances to produce a composite unit.
func Prefix(symbol string, mult float64) func(inner Maker) Maker {
	return func(inner Maker) Maker {
		iu := inner.Unit()
		if iu == nil {
			panic(fmt.Sprintf("prefix %q must wrap a singular unit, got %q", symbol, inner.Units()))
		}
		p := prefixType{
			prefix: symbol,
			mult:   mult,
			inner:  iu,
		}
		p.units = Units{N: []Unit{p}}
		return p.Make
	}
}
