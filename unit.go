package unit

import "fmt"

type Unit interface {
	Symbol() string
	Value() Value
	Equal(b Unit) bool
}

type unitType struct {
	symbol string
	value  Value
	units  Units
}

func (u unitType) String() string       { return fmt.Sprintf("Unit(%q = %v)", u.Symbol(), u.Value()) }
func (u unitType) Symbol() string       { return u.symbol }
func (u unitType) Value() Value         { return u.value }
func (u unitType) Equal(o Unit) bool    { return unitEqual(u, o) }
func (u unitType) Make(v float64) Value { return u.units.Make(v) }

func unitEqual(a, b Unit) bool {
	if a == nil {
		return b == nil
	}
	if b == nil {
		return false
	}
	return a.Symbol() == b.Symbol()
}

func Scalar(v float64) Maker { return Value{S: v}.MulN }

var (
	Unity = Scalar(1)
)

type Maker func(v float64) Value

func (m Maker) Value() float64 { return m(1).Value() }
func (m Maker) Units() Units   { return m(1).Units() }

func (m Maker) Div(b Maker) Maker {
	return func(f float64) Value { return m(f).Div(b) }
}
func (m Maker) Mul(b Maker) Maker {
	return func(f float64) Value { return m(f).Mul(b) }
}
func (m Maker) Pow(n int) Maker {
	return func(f float64) Value { return m(f).Mul(m(1).Pow(n - 1)) }
}

/*
func (m Maker) Add(n float64) Maker {
	return func(f float64) Value { return m(f).AddN(n) }
}
func (m Maker) Sub(n float64) Maker {
	return func(f float64) Value { return m(f).SubN(n) }
}
*/

func (m Maker) String() string {
	v := m(1)
	f := v.Value()
	if f != 1 || v.U.IsEmpty() {
		return fmt.Sprintf("*%g", f)
	}
	return v.U.String()
}
func (m Maker) Unit() Unit {
	us := m.Units()
	if len(us.N) == 1 && len(us.D) == 0 {
		return us.N[0]
	}
	return nil
}

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
		value:  FromQualified(value),
	}
	u.units = Units{N: []Unit{u}}
	return u
}

func Primitive(symbol string) Maker {
	return Derive(symbol, Unity)
}

func IsPrimitive(u Unit) bool {
	if u == nil {
		return false
	}
	v := u.Value()
	return v.S == 1 && v.U.IsEmpty()
}

func Derive(symbol string, value Qualified) Maker {
	return newUnit(symbol, FromQualified(value)).Make
}

func Must(v Value, e error) Value {
	if e != nil {
		panic(e)
	}
	return v
}
