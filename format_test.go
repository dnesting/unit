package unit_test

import (
	"fmt"
	"testing"

	"github.com/dnesting/unit"
)

func TestDefaultFormatter(t *testing.T) {
	m := unit.Primitive("m")
	kg := unit.Primitive("kg")
	s := unit.Primitive("s")
	v := kg.Mul(m).Div(s.Pow(2))(1.234)

	expected := "1.234 kg m/s^2"
	actual := unit.DefaultFormatter.Format(v)
	if expected != actual {
		t.Errorf("default format should be %q, got %q", expected, actual)
	}
}

func opts(o ...unit.FormatOpt) []unit.FormatOpt { return o }

func TestFormatter(t *testing.T) {
	m := unit.Primitive("m")
	kg := unit.Primitive("kg")
	s := unit.Primitive("s")
	v := kg.Mul(m).Div(s.Pow(2))(1.234)

	for _, c := range []struct {
		desc     string
		expected string
		opts     []unit.FormatOpt
	}{
		{"basic", "1.234 kg m/s^2", opts()},
		{"unitfn", "1.234 _kg_ _m_/_s_^2", opts(
			unit.WithUnitFunc(func(u unit.Unit) string { return fmt.Sprintf("_%s_", u.Symbol()) }))},
		{"fraction", "1.234 kg m|s^2", opts(unit.WithFraction("|"))},
		{"unitsep", "1.234 kg-m/s^2", opts(unit.WithUnitSep("-"))},
		{"valuesep", "1.234&nbsp;kg-m/s^2", opts(unit.WithBeforeUnits("&nbsp;"), unit.WithUnitSep("-"))},
		{"nogap", "1.234kg m/s^2", opts(unit.WithNoGap())},
		{"fmt", "1.23400 kg m/s^2", opts(unit.WithFmt("%.5f"))},
		{"nofraction", "1.234 kg m s^-2", opts(unit.WithNoFraction())},
		{"unicode", "1.234 kg⋅m⁄s²", opts(unit.WithUnicode())},
		{"mathml", "<mn>1.234</mn> <mfrac><mrow><mi>kg</mi><mo>⋅</mo><mi>m</mi></mrow><mrow><msup><mi>s</mi><mn>2</mn></msup></mrow></mfrac>",
			opts(unit.WithMathML())},
		{"mathml-nogap", "<mn>1.234</mn><mfrac><mrow><mi>kg</mi><mo>⋅</mo><mi>m</mi></mrow><mrow><msup><mi>s</mi><mn>2</mn></msup></mrow></mfrac>",
			opts(unit.WithMathML(), unit.WithNoGap())},
	} {
		t.Run(c.desc, func(t *testing.T) {
			f := unit.NewFormatter(c.opts...)
			actual := f.Format(v)
			if c.expected != actual {
				t.Errorf("expected %q, got %q", c.expected, actual)
			}
		})
	}
}
