package si_test

import (
	"math"
	"testing"

	"github.com/dnesting/unit"
	"github.com/dnesting/unit/si"
)

func TestConstants(t *testing.T) {
	for _, c := range []struct {
		name     string
		val      unit.Value
		expVal   float64
		expUnits string
	}{
		{"Caesium", si.Caesium, 9.19263177e+09, "Hz"},
		{"c", si.C, 2.99792458e+08, "m/s"},
		{"h", si.H, 6.62607015e-34, "J s"},
		{"e", si.E, 1.602176634e-19, "C"},
		{"k", si.K, 1.380649e-23, "J/K"},
		{"Na", si.NA, 6.02214076e23, "1/mol"},
		{"Kcd", si.Kcd, 683, "cd sr/W"},
	} {
		actual := c.val
		if math.Abs(c.expVal-actual.S) > 0.00000000001 || c.expUnits != actual.U.String() {
			t.Errorf("%v: expected %g %v, got %q", c.name, c.expVal, c.expUnits, actual)
		}
	}
}

func TestCelsius(t *testing.T) {
	c := si.DegCelsius(0)
	k := si.Kelvin(273.15)
	r, ok := si.CToK(c)
	eq, err := k.Equal(r)
	if !ok || !eq || err != nil {
		t.Errorf("CToK(%q) should result in %q, got %q (eq=%v err=%v)", c, k, r, eq, err)
	}
	c = si.DegCelsius(0)
	k = si.Kelvin(273.15)
	r, ok = si.KToC(k)
	eq, err = c.Equal(r)
	if !ok || !eq || err != nil {
		t.Errorf("KToC(%q) should result in %q, got %q (eq=%v err=%v)", k, c, r, eq, err)
	}

	c = si.DegCelsius(100)
	k = si.Kelvin(373.15)
	r, ok = si.CToK(c)
	if !ok || !eq || err != nil {
		t.Errorf("CToK(%q) should result in %q, got %q (eq=%v err=%v)", c, k, r, eq, err)
	}
	c = si.DegCelsius(100)
	k = si.Kelvin(373.15)
	r, ok = si.KToC(k)
	eq, err = c.Equal(r)
	if !ok || !eq || err != nil {
		t.Errorf("KToC(%q) should result in %q, got %q (eq=%v err=%v)", k, c, r, eq, err)
	}
}

func TestPrefixes(t *testing.T) {
	a := si.Gram(2000)
	b := si.Kilogram(2)
	eq, err := a.Equal(b)
	if !eq || err != nil {
		t.Errorf("%q should equal %q, got eq=%v err=%v", a, b, eq, err)
	}

	a = si.Milli(si.Gram).Mul(si.Kilo(si.Meter)).Div(si.Second.Pow(2))(5.123)
	e := "5.123 km mg/s^2"
	r := a.String()
	if e != r {
		t.Errorf("expected %q, got %q", e, r)
	}
}
