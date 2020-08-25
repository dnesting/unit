package unit_test

import (
	"reflect"
	"testing"

	"github.com/dnesting/unit"
)

func TestUnit(t *testing.T) {
	meter := unit.Primitive("m")
	second := unit.Primitive("s")
	mdivs := meter.Div(second)

	if !unit.IsPrimitive(meter.Unit()) {
		t.Errorf("unit.Primitive() should satisfy IsPrimitive")
	}
	if unit.IsPrimitive(mdivs.Unit()) {
		t.Errorf("Anonymous derived unit should not satisfy IsPrimitive")
	}

	if !reflect.DeepEqual(mdivs.Units().N, []unit.Unit{meter.Unit()}) || !reflect.DeepEqual(mdivs.Units().D, []unit.Unit{second.Unit()}) {
		t.Errorf("meter.Div(second) should return Units with meter over second, got %v", mdivs.Units())
	}

	mps := unit.Derive("mps", mdivs)
	if unit.IsPrimitive(mps.Unit()) {
		t.Errorf("Named derived unit should not satisfy IsPrimitive")
	}

	mps12 := mps(1.2)
	expected := "1.2 mps"
	if mps12.String() != "1.2 mps" {
		t.Errorf("mps(1.2) should stringify %q, got %q", expected, mps12.String())
	}

	mdivs2kg := mdivs.Div(second).Div(unit.Primitive("kg"))
	mdivs2kg34 := mdivs2kg(3.4)
	expected = "3.4 m/kg s^2"
	if mdivs2kg34.String() != expected {
		t.Errorf("mdivs2kg(3.4) should stringify %q, got %q", expected, mdivs2kg34.String())
	}
}
