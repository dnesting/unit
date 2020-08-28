package unit_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/dnesting/unit"
)

/*
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
*/

func TestEmpty(t *testing.T) {
	prim := unit.Primitive("prim").Unit()
	var u unit.Units
	if !u.Empty() {
		t.Errorf("empty units should be empty")
	}
	u = unit.Units{
		N: []unit.Unit{prim},
	}
	if u.Empty() {
		t.Errorf("non-empty units should not be empty")
	}
	u = unit.Units{
		D: []unit.Unit{prim},
	}
	if u.Empty() {
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
	if !e.Equal(r) {
		t.Errorf("ka(2) should equal a(2000), expected %v, got %v", e, r)
	}

	ka = unit.Derive("ka", a.Mul(unit.Scalar(1000)))
	r = ka(2)
	if !e.Equal(r) {
		t.Errorf("ka(2) should equal a(2000), expected %v, got %v", e, r)
	}
	if !r.Equal(e) {
		t.Errorf("ka(2) should equal a(2000), expected %v, got %v", e, r)
	}
}

func TestSameSymbol(t *testing.T) {
	// We support two units with the same symbol, such as us.Mile and us.SurveyMile.
	m := unit.Primitive("m")
	mi := unit.Derive("mi", m(5280*12*0.0254))
	smi := unit.Derive("mi", m(5280*12*100.0/3937))

	if mi(1).Equal(smi(1)) {
		t.Errorf("%q and %q should not be equal", mi(1), smi(1))
	}

	if mi.Units().Equal(smi.Units()) {
		t.Errorf("%q and %q should not be equal", mi, smi)
	}

	a := smi(3)
	b := mi(2)
	e := mi(2)
	r := a.Mul(b).Div(a)
	if !e.Equal(r) {
		t.Errorf("%q.Div(%q) should yield %q, got %q", a.Mul(b), a, e, r)
	}

	a = smi(2)
	b = mi(2)
	p := smi(2).Pow(2)
	r = a.Mul(b)
	if p.Equal(r) {
		t.Errorf("%q and %q should not be equal", p, r)
	}
}

func ExampleUnits() {
	kg := unit.Primitive("kg")
	m := unit.Primitive("m")
	s := unit.Primitive("s")
	n := unit.Derive("N", kg.Mul(m).Div(s.Pow(2)))

	x := n(1.234)
	u := x.Units()
	r := u.Reduce().Units()
	fmt.Println(u)
	fmt.Println(r)
	fmt.Println("u.Equal(r)?", u.Equal(r))
	fmt.Println("u.Equiv(r)?", u.Equiv(r))
	// Output:
	// N
	// kg m/s^2
	// u.Equal(r)? false
	// u.Equiv(r)? true
}
