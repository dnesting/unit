package unit_test

import (
	"fmt"
	"testing"

	"github.com/dnesting/unit"
)

func TestValueEqual(t *testing.T) {
	meter := unit.Primitive("m")
	second := unit.Primitive("s")
	a := meter(5)
	if !a.Equal(a) {
		t.Errorf("%v and itself should be equal", a)
	}
	b := meter(5)
	if !a.Equal(b) {
		t.Errorf("%v and %v should be equal", a, b)
	}
	b = meter(6)
	if a.Equal(b) {
		t.Errorf("%v and %v should NOT be equal but have compatible units", a, b)
	}
	b = second(5)
	if a.Equal(b) {
		t.Errorf("%v and %v should NOT be equal with incompatible units", a, b)
	}

	mps := unit.Derive("mps", meter.Div(second))
	a = mps(2)
	b = meter.Div(second)(2)
	if !a.Equal(b) {
		t.Errorf("%v and %v should be equal because their units are equivalent", a, b)
	}
}

/*
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
*/

func TestMath(t *testing.T) {
	meter := unit.Primitive("m")
	second := unit.Primitive("s")
	mps := meter.Div(second)

	a := mps(2)
	e := mps(6)
	r := a.MulN(3)
	if !e.Equal(r) {
		t.Errorf("%v.MulN(3) should give us %v, got %v", a, e, r)
	}

	b := mps(3)
	e = meter.Pow(2).Div(second.Pow(2))(6)
	r = a.Mul(b)
	if !e.Equal(r) {
		t.Errorf("%v.Mul(%v) should give us %v, got %v", a, b, e, r)
	}

	a = meter.Pow(2).Div(second.Pow(2))(6)
	b = mps(3)
	e = mps(2)
	r = a.Div(b)
	if !e.Equal(r) {
		t.Errorf("%v.Div(%v) should give us %v, got %v", a, b, e, r)
	}

	a = mps(6)
	e = mps(2)
	r = a.DivN(3)
	if !e.Equal(r) {
		t.Errorf("%v.DivN(3) should give us %v, got %v", a, e, r)
	}

	a = mps(2)
	e = meter.Pow(3).Div(second.Pow(3))(8)
	r = a.Pow(3) // 8 m^3/s^3
	if !e.Equal(r) {
		t.Errorf("%v.Pow(3) should give us %v, got %v", a, e, r)
	}

	a = meter.Pow(2).Div(second.Pow(2))(4)
	e = meter.Pow(4).Div(second.Pow(4))(float64(1) / 16)
	r = a.Pow(-2)
	if !e.Equal(r) {
		t.Errorf("%v.Pow(-2) should give us %v, got %v", a, e, r)
	}

	a = mps(4)
	e = unit.Scalar(1)(1)
	r = a.Pow(0)
	if !e.Equal(r) {
		t.Errorf("%v.Pow(0) should give us %v, got %v", a, e, r)
	}

	e = a
	r = a.Pow(1)
	if !e.Equal(r) {
		t.Errorf("%v.Pow(1) should give us %v, got %v", a, e, r)
	}

	a = mps(3)
	e = mps(5)
	r = a.AddN(2)
	if !e.Equal(r) {
		t.Errorf("%v.AddN(2) should give us %v, got %v", a, e, r)
	}

	a = mps(2)
	b = mps(3)
	e = mps(5)
	r, ok := a.Add(b)
	if !e.Equal(r) || !ok {
		t.Errorf("%v.Add(%v) should give us %v, got %v (ok=%v)", a, b, e, r, ok)
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
	if !e.Equal(r) || !ok {
		t.Errorf("%v.Sub(%v) should give us %v, got %v (ok=%v)", a, b, e, r, ok)
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
	if !r.Approx(e, 0.00001) || !extra.Empty() {
		t.Errorf("%q.Convert(%v) should give %v, got %v (extra=%q)", a.String(), m, e, r, extra)
	}

	a = m(2.54 / 100 * 2)
	e = in(2)
	r, extra = a.Convert(in)
	if !r.Approx(e, 0.00001) || !extra.Empty() {
		t.Errorf("%q.Convert(%v) should give %v, got %v (extra=%q)", a.String(), in, e, r, extra)
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
		{"/s", hz, "unit /s only"},
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

func ExampleValue_Convert() {
	// Set up some initial units
	kg := unit.Primitive("kg")
	m := unit.Primitive("m")
	s := unit.Primitive("s")
	n := unit.Derive("N", kg.Mul(m).Div(s.Pow(2)))
	j := unit.Derive("J", n.Mul(m))

	x := j(1.234)
	fmt.Println("x =", x)

	// Convert "1.234 J" to "N m", which should be a conforming conversion.
	xnm, extra := x.Convert(n.Mul(m))
	fmt.Println("as N-m =", xnm, "conforming?", extra.Empty())

	// Convert "1.234 J" to "N", which should be non-conforming, with extra units
	// returned in extra.
	xn, extra := x.Convert(n)
	fmt.Println("as N =", xn, "conforming?", extra.Empty())

	// Output:
	// x = 1.234 J
	// as N-m = 1.234 N m conforming? true
	// as N = 1.234 N conforming? false
}

func ExampleValue_Reduce() {
	// Set up some initial units
	kg := unit.Primitive("kg")
	m := unit.Primitive("m")
	s := unit.Primitive("s")
	n := unit.Derive("N", kg.Mul(m).Div(s.Pow(2)))
	j := unit.Derive("J", n.Mul(m))

	x := j(1.234)
	fmt.Println("x is", x)
	fmt.Println("x reduces to", x.Reduce())

	// Output:
	// x is 1.234 J
	// x reduces to 1.234 kg m^2/s^2
}
