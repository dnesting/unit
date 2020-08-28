package us_test

import (
	"testing"

	"github.com/dnesting/unit/si"
	"github.com/dnesting/unit/us"
)

func TestFtoC(t *testing.T) {
	freezing := us.DegFahrenheit(32)
	boiling := us.DegFahrenheit(212)

	e := si.DegCelsius(0)
	r, ok := us.ToCelsius(freezing)
	if !ok || !e.Approx(r, 0.0000000001) {
		t.Errorf("FToC should convert %q, expected %q, got %q (ok=%v)", freezing, e, r, ok)
	}

	e = si.DegCelsius(100)
	r, ok = us.ToCelsius(boiling)
	if !ok || !e.Approx(r, 0.0000000001) {
		t.Errorf("FToC should convert %q, expected %q, got %q (ok=%v)", boiling, e, r, ok)
	}

	freezing = si.DegCelsius(0)
	boiling = si.DegCelsius(100)

	e = us.DegFahrenheit(32)
	r, ok = us.ToFahrenheit(freezing)
	if !ok || !e.Approx(r, 0.0000000001) {
		t.Errorf("CToF should convert %q, expected %q, got %q (ok=%v)", freezing, e, r, ok)
	}

	e = us.DegFahrenheit(212)
	r, ok = us.ToFahrenheit(boiling)
	if !ok || !e.Approx(r, 0.0000000001) {
		t.Errorf("CToF should convert %q, expected %q, got %q (ok=%v)", boiling, e, r, ok)
	}

	a := si.Meter(5)
	r, ok = us.ToCelsius(a)
	if ok {
		t.Errorf("FToC should not convert %q, got %q (ok=%v)", a, r, ok)
	}
}

func TestDMS(t *testing.T) {
	orig := us.Degree(100.2625)
	ed := us.Degree(100)
	em := us.DegMinute(15)
	es := us.DegSecond(45)
	d, m, s, ok := us.ToDMS(orig)
	if !ok || !ed.Equal(d) || !em.Equal(m) || !es.Approx(s, 0.000001) {
		t.Errorf("us.DMS(%q) should return (%q, %q, %q, true), got (%q, %q, %q, %v)",
			orig, ed, em, es, d, m, s, ok)
	}
}
