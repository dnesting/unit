package unit_test

import (
	"reflect"
	"testing"

	"github.com/dnesting/unit"
)

func TestSingular(t *testing.T) {
	prim := unit.Primitive("prim").Unit()
	us := unit.Units{
		N: []unit.Unit{prim},
	}
	u := us.Singular()
	if u == nil {
		t.Errorf("Units with a single unit should have non-nil Singular()")
	}
	if u.Symbol() != "prim" {
		t.Errorf("Units with single unit should return same unit")
	}

	us = unit.Units{
		N: []unit.Unit{prim, prim},
	}
	u = us.Singular()
	if u != nil {
		t.Errorf("Units with multiple units should have nil Singular(), got %v", u)
	}

	us = unit.Units{
		D: []unit.Unit{prim},
	}
	u = us.Singular()
	if u != nil {
		t.Errorf("Units with units only in denominator should have nil Singular(), got %v", u)
	}
}

func TestEmpty(t *testing.T) {
	prim := unit.Primitive("prim").Unit()
	var u unit.Units
	if !u.IsEmpty() {
		t.Errorf("empty units should be empty")
	}
	u = unit.Units{
		N: []unit.Unit{prim},
	}
	if u.IsEmpty() {
		t.Errorf("non-empty units should not be empty")
	}
	u = unit.Units{
		D: []unit.Unit{prim},
	}
	if u.IsEmpty() {
		t.Errorf("non-empty units should not be empty")
	}
}

func TestRecip(t *testing.T) {
	a := unit.Primitive("a").Unit()
	b := unit.Primitive("b").Unit()

	u := unit.Units{N: []unit.Unit{a}, D: []unit.Unit{b}}
	expected := unit.Units{N: []unit.Unit{b}, D: []unit.Unit{a}}
	actual := u.Recip()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("reciprocal of %v should be %v, got %v", u, expected, actual)
	}

	u = unit.Units{D: []unit.Unit{a}}
	expected = unit.Units{N: []unit.Unit{a}}
	actual = u.Recip()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("reciprocal of %v should be %v, got %v", u, expected, actual)
	}
}

func TestEqual(t *testing.T) {
	a := unit.Primitive("a")
	b := unit.Primitive("b")
	ab := a.Div(b)

	var us unit.Units
	us = unit.Units{N: []unit.Unit{a.Unit()}, D: []unit.Unit{b.Unit()}}
	if !us.Equal(ab.Units()) {
		t.Errorf("%v should equal %v", us, ab.Units())
	}

	us = unit.Units{N: []unit.Unit{a.Unit()}, D: []unit.Unit{b.Unit(), b.Unit()}}
	if us.Equal(ab.Units()) {
		t.Errorf("%v should not equal %v", us, ab.Units())
	}
}

func TestEquiv(t *testing.T) {
	a := unit.Primitive("a")
	b := unit.Primitive("b")
	c := unit.Derive("c", a.Div(b))

	var us unit.Units
	us = unit.Units{N: []unit.Unit{a.Unit()}, D: []unit.Unit{b.Unit()}}
	if us.Equal(c.Units()) {
		t.Errorf("%v should not be equal to %v", us, c.Units())
	}
	if !us.Equiv(c.Units()) {
		t.Errorf("%v should be equivalent to %v", us, c.Units())
	}
}

/*
func TestCancel(t *testing.T) {
	a := unit.Primitive("a").Unit()
	b := unit.Primitive("b").Unit()
	c := unit.Primitive("c").Unit()

	us := unit.Units{
		N: []unit.Unit{a, a, b},
		D: []unit.Unit{a, b, c},
	}
	expected := unit.Units{
		N: []unit.Unit{a},
		D: []unit.Unit{c},
	}
	us.Cancel()
	if !reflect.DeepEqual(expected, us) {
		t.Errorf("Cancel expected %v, got %v", expected, us)
	}
}
*/

func TestReduce(t *testing.T) {
	a := unit.Primitive("a")
	b := unit.Primitive("b")
	c := unit.Derive("c", a.Div(b))

	us := c.Units()
	expected := a.Div(b).Units()
	r := us.Reduce()
	if r.S != 1 {
		t.Errorf("Reduce multipler should be 1, got %v", r.S)
	}
	if !reflect.DeepEqual(expected, r.U) {
		t.Errorf("Reduce expected %v, got %v", expected, r.U)
	}
}

func TestMakeWithMultiply(t *testing.T) {
	a := unit.Primitive("a")
	ka := a.Mul(unit.Scalar(1000))

	e := a(2000)
	r := ka(2)
	if eq, err := e.Equal(r); !eq || err != nil {
		t.Errorf("ka(2) should equal a(2000), expected %v, got %v (err=%v)", e, r, err)
	}

	ka = unit.Derive("ka", a.Mul(unit.Scalar(1000)))
	r = ka(2)
	if eq, err := e.Equal(r); !eq || err != nil {
		t.Errorf("ka(2) should equal a(2000), expected %v, got %v (err=%v)", e, r, err)
	}
	if eq, err := r.Equal(e); !eq || err != nil {
		t.Errorf("ka(2) should equal a(2000), expected %v, got %v (err=%v)", e, r, err)
	}
}
