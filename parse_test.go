package unit_test

import (
	"strings"
	"testing"

	"github.com/dnesting/unit"
)

func TestParse(t *testing.T) {
	var r unit.Registry
	m := r.Primitive("m")
	s := r.Primitive("s")
	k := r.Prefix("k", 1000)
	mps := r.Derive("mps", m.Div(s))
	foo := unit.Primitive("foo") // not registered

	ohm := r.Primitive("Ω")
	degC := r.Primitive("°C")

	var cases = []struct {
		test      string
		mustExist bool
		expected  unit.Qualified
		err       string
		desc      string
	}{
		{"23", true, unit.Scalar(23), "", "unitless value"},
		{"23.", true, unit.Scalar(23), "", "unitless value."},
		{"23.000", true, unit.Scalar(23), "", "unitless value.000"},
		{".23", true, unit.Scalar(0.23), "", "unitless value.23"},
		{"23e2", true, unit.Scalar(2300), "", "unitless with exponent"},
		{"2.3e2", true, unit.Scalar(230), "", "unitless with exponent and decimal"},

		{"m", true, m(1), "", "bare unit"},
		{"m/s", true, m.Div(s)(1), "", "bare fractioned unit"},
		{"m^3/s^2", true, m.Pow(3).Div(s.Pow(2))(1), "", "bare complex unit"},

		{"23m", true, m(23), "", "basic value"},
		{"23  m", true, m(23), "", "basic value with spaces"},
		{"23  m/s", true, m.Div(s)(23), "", "fractioned"},
		{"23 m^3/s^2", true, m.Pow(3).Div(s.Pow(2))(23), "", "complex"},

		{"23 m mps", true, m.Mul(mps)(23), "", "with derived"},

		{"23 foo", true, nil, "unknown unit", "missing unit"},
		{"23 foo", false, foo(23), "", "missing unit but ok"},
		{"23 Ω", true, ohm(23), "", "ohm symbol"},
		{"23 °C", true, degC(23), "", "degrees Celsius symbol"},

		{"23 km/ks", true, k(m).Div(k(s))(23), "", "with prefix"},
		{"23 km/ks", true, m.Div(s)(23), "", "with prefix, reduced"},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			v, err := unit.Parse(c.test, &r, c.mustExist)
			if c.err == "" {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
					return
				}
			} else {
				if !strings.Contains(err.Error(), c.err) {
					t.Errorf("expected error with %q, got %v", c.err, err)
				}
			}
			if c.expected != nil {
				if eq, _ := v.Equal(c.expected); !eq {
					t.Errorf("expected %v, got %v", c.expected, v)
				}
			}
		})
	}
}
