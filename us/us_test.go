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
	eq, err := e.Approx(r, 0.0000000001)
	if !ok || !eq || err != nil {
		t.Errorf("FToC should convert %q, expected %q, got %q (ok=%v, err=%v)", freezing, e, r, ok, err)
	}

	e = si.DegCelsius(100)
	r, ok = us.ToCelsius(boiling)
	eq, err = e.Approx(r, 0.0000000001)
	if !ok || !eq || err != nil {
		t.Errorf("FToC should convert %q, expected %q, got %q (ok=%v, err=%v)", boiling, e, r, ok, err)
	}

	freezing = si.DegCelsius(0)
	boiling = si.DegCelsius(100)

	e = us.DegFahrenheit(32)
	r, ok = us.ToFahrenheit(freezing)
	eq, err = e.Approx(r, 0.0000000001)
	if !ok || !eq || err != nil {
		t.Errorf("CToF should convert %q, expected %q, got %q (ok=%v, err=%v)", freezing, e, r, ok, err)
	}

	e = us.DegFahrenheit(212)
	r, ok = us.ToFahrenheit(boiling)
	eq, err = e.Approx(r, 0.0000000001)
	if !ok || !eq || err != nil {
		t.Errorf("CToF should convert %q, expected %q, got %q (ok=%v, err=%v)", boiling, e, r, ok, err)
	}

	a := si.Meter(5)
	r, ok = us.ToCelsius(a)
	if ok {
		t.Errorf("FToC should not convert %q, got %q (ok=%v)", a, r, ok)
	}
}
