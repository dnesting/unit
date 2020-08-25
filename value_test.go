package unit_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/dnesting/unit"
)

func TestValueEqual(t *testing.T) {
	meter := unit.Primitive("m")
	second := unit.Primitive("s")
	a := meter(5)
	if eq, err := a.Equal(a); !eq || err != nil {
		t.Errorf("%v and itself should be equal, err=%v", a, err)
	}
	b := meter(5)
	if eq, err := a.Equal(b); !eq || err != nil {
		t.Errorf("%v and %v should be equal, err=%v", a, b, err)
	}
	b = meter(6)
	if eq, err := a.Equal(b); eq || err != nil {
		t.Errorf("%v and %v should NOT be equal but have compatible units, err=%v", a, b, err)
	}
	b = second(5)
	if eq, err := a.Equal(b); eq || err == nil {
		t.Errorf("%v and %v should NOT be equal with incompatible units", a, b)
	}

	mps := unit.Derive("mps", meter.Div(second))
	a = mps(2)
	b = meter.Div(second)(2)
	if eq, err := a.Equal(b); !eq || err != nil {
		t.Errorf("%v and %v should be equal because their units are equivalent, err=%v", a, b, err)
	}
}

func TestValueLess(t *testing.T) {
	meter := unit.Primitive("m")
	second := unit.Primitive("s")
	mps := unit.Derive("mps", meter.Div(second))

	a := mps(5)
	b := mps(6)

	if lt, err := a.Less(b); !lt || err != nil {
		t.Errorf("%v should be less than %v, got err=%v", a, b, err)
	}
	if lt, err := b.Less(a); lt || err != nil {
		t.Errorf("%v should NOT be less than %v, got err=%v", b, a, err)
	}
	b = meter.Div(second)(7)
	if lt, err := a.Less(b); !lt || err != nil {
		t.Errorf("%v should be less than %v with equivalent units, got err=%v", a, b, err)
	}
	b = second(8)
	if lt, err := a.Less(b); lt || err == nil {
		t.Errorf("%v should be NOT less than %v with incompatible units", a, b)
	}
}

func TestMath(t *testing.T) {
	meter := unit.Primitive("m")
	second := unit.Primitive("s")
	mps := meter.Div(second)

	a := mps(2)
	e := mps(6)
	r := a.MulN(3)
	if ok, _ := e.Equal(r); !ok {
		t.Errorf("%v.MulN(3) should give us %v, got %v", a, e, r)
	}

	b := mps(3)
	e = meter.Pow(2).Div(second.Pow(2))(6)
	r = a.Mul(b)
	if eq, err := e.Equal(r); !eq || err != nil {
		t.Errorf("%v.Mul(%v) should give us %v, got %v (err=%v)", a, b, e, r, err)
	}

	a = meter.Pow(2).Div(second.Pow(2))(6)
	b = mps(3)
	e = mps(2)
	r = a.Div(b)
	if eq, err := e.Equal(r); !eq || err != nil {
		t.Errorf("%v.Div(%v) should give us %v, got %v (err=%v)", a, b, e, r, err)
	}

	a = mps(6)
	e = mps(2)
	r = a.DivN(3)
	if eq, err := e.Equal(r); !eq || err != nil {
		t.Errorf("%v.DivN(3) should give us %v, got %v", a, e, r)
	}

	a = mps(2)
	e = meter.Pow(3).Div(second.Pow(3))(8)
	r = a.Pow(3) // 8 m^3/s^3
	if eq, err := e.Equal(r); !eq || err != nil {
		t.Errorf("%v.Pow(3) should give us %v, got %v (err=%v)", a, e, r, err)
	}

	a = meter.Pow(2).Div(second.Pow(2))(4)
	e = meter.Pow(4).Div(second.Pow(4))(float64(1) / 16)
	r = a.Pow(-2)
	if eq, err := e.Equal(r); !eq || err != nil {
		t.Errorf("%v.Pow(-2) should give us %v, got %v (err=%v)", a, e, r, err)
	}

	a = mps(4)
	e = unit.Scalar(1)(1)
	r = a.Pow(0)
	if eq, err := e.Equal(r); !eq || err != nil {
		t.Errorf("%v.Pow(0) should give us %v, got %v (err=%v)", a, e, r, err)
	}

	e = a
	r = a.Pow(1)
	if eq, err := e.Equal(r); !eq || err != nil {
		t.Errorf("%v.Pow(1) should give us %v, got %v (err=%v)", a, e, r, err)
	}

	a = mps(3)
	e = mps(5)
	r = a.AddN(2)
	if eq, err := e.Equal(r); !eq || err != nil {
		t.Errorf("%v.AddN(2) should give us %v, got %v (err=%v)", a, e, r, err)
	}

	a = mps(2)
	b = mps(3)
	e = mps(5)
	r, ok := a.Add(b)
	if eq, err := e.Equal(r); !ok || !eq || err != nil {
		t.Errorf("%v.Add(%v) should give us %v, got %v (ok=%v eq=%v err=%v)", a, b, e, r, ok, eq, err)
	}
	b = meter(1) // different unit than a
	r, ok = a.Add(b)
	if ok {
		t.Errorf("%v.Add(%v) should fail, got %q", a, b, r)
	}

	a = mps(5)
	b = mps(2)
	e = mps(3)
	r, ok = a.Sub(b)
	if eq, err := e.Equal(r); !ok || !eq || err != nil {
		t.Errorf("%v.Sub(%v) should give us %v, got %v (ok=%v eq=%v err=%v)", a, b, e, r, ok, eq, err)
	}
	b = meter(1)
	r, ok = a.Sub(b)
	if ok {
		t.Errorf("%v.Sub(%v) should fail, got %q", a, b, r)
	}
}

func TestConvert(t *testing.T) {
	m := unit.Primitive("m")
	cm := unit.Derive("cm", m.Div(unit.Scalar(100)))
	in := unit.Derive("in", cm(2.54))

	a := in(2)
	e := m(2.54 * 2 / 100)
	r, extra := a.Convert(m)
	eq, err := r.Compare(e, func(a, b float64) bool { return math.Round(a*10000)/10000 == math.Round(b*10000)/10000 })
	if !eq || err != nil || !extra.IsEmpty() {
		t.Errorf("%q.Convert(%v) should give %v, got %v (extra=%q eq=%v err=%v)", a.String(), m, e, r, extra, eq, err)
	}

	a = m(2.54 / 100 * 2)
	e = in(2)
	r, extra = a.Convert(m)
	eq, err = r.Compare(e, func(a, b float64) bool { return math.Round(a*10000)/10000 == math.Round(b*10000)/10000 })
	if !eq || err != nil || !extra.IsEmpty() {
		t.Errorf("%q.Convert(%v) should give %v, got %v (extra=%q eq=%v err=%v)", a.String(), m, e, r, extra, eq, err)
	}
}

func TestString(t *testing.T) {
	m := unit.Primitive("m")
	s := unit.Primitive("s")
	mps2 := m.Div(s.Pow(2))
	hz := unit.Scalar(1).Div(s)

	for _, c := range []struct {
		expected string
		actual   unit.Qualified
		desc     string
	}{
		{"5 m/s^2", mps2(5), "basic"},
		{"5.1234 m/s^2", mps2(5.1234), "real"},
		{"0.1234 m/s^2", mps2(0.1234), "real starting with zero"},

		{"m", m, "unit m only"},
		{"1/s", hz, "unit 1/s only"},
		{"m/s^2", mps2, "unit mps2 only"},
	} {
		t.Run(c.desc, func(t *testing.T) {
			actual := fmt.Sprintf("%s", c.actual)
			if actual != c.expected {
				t.Errorf("expected %q, got %q", c.expected, actual)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	m := unit.Primitive("m")
	s := unit.Primitive("s")
	mps2 := m.Div(s.Pow(2))

	for _, c := range []struct {
		format   string
		expected string
		actual   unit.Qualified
		desc     string
	}{
		{"%g", "5", mps2(5), "g-integer"},
		{"%g", "5.1234", mps2(5.1234), "g-real"},
		{"%g", "0.1234", mps2(0.1234), "g-real with zero"},
		{"%10.2g", "         5", mps2(5), "g-integer-10.2"},
		{"%10.2g", "       5.1", mps2(5.1234), "g-real-10.2"},
		{"%10.2g", "      0.12", mps2(0.1234), "g-real with zero-10.2"},
		{"%f", "5.000000", mps2(5), "f-integer"},
		{"%f", "5.123400", mps2(5.1234), "f-real"},
		{"%f", "0.123400", mps2(0.1234), "f-real with zero"},
		{"%10.2f", "      5.00", mps2(5), "f-integer-10.2"},
		{"%10.2f", "      5.12", mps2(5.1234), "f-real-10.2"},
		{"%10.2f", "      0.12", mps2(0.1234), "f-real with zero-10.2"},

		{"%s", "5 m/s^2", mps2(5), "s-integer"},
		{"%s", "5.1234 m/s^2", mps2(5.1234), "s-real"},
		{"%v", "5 m/s^2", mps2(5), "v-integer"},
		{"%v", "5.1234 m/s^2", mps2(5.1234), "v-real"},
		{"%q", "\"5 m/s^2\"", mps2(5), "q-integer"},
		{"%q", "\"5.1234 m/s^2\"", mps2(5.1234), "q-real"},
	} {
		t.Run(c.desc, func(t *testing.T) {
			actual := fmt.Sprintf(c.format, c.actual)
			if actual != c.expected {
				t.Errorf("expected %q, got %q", c.expected, actual)
			}
		})
	}
}
