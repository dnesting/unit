package unit

import (
	"fmt"
	"strings"
)

type Formatter struct {
	unitFn         func(u Unit) string
	powerFn        func(u string, p int) string
	valueFn        func(tmpl string, v float64) string
	fracFn         func(num, denom string) string
	negativePowers bool
	valueFmt       string

	beforeUnits   string
	beforeUnitRow string
	unitSep       string
	afterUnitRow  string
	afterUnits    string
}

type FormatOpt func(f *Formatter)

func WithUnitFunc(fn func(u Unit) string) FormatOpt {
	return func(f *Formatter) { f.unitFn = fn }
}

func WithPowerFunc(fn func(u string, p int) string) FormatOpt {
	return func(f *Formatter) { f.powerFn = fn }
}

func WithFraction(frac string) FormatOpt {
	return func(f *Formatter) {
		f.fracFn = func(n, d string) string {
			if d != "" {
				return n + frac + d
			}
			return n
		}
	}
}

func WithUnitSep(sep string) FormatOpt {
	return func(f *Formatter) { f.unitSep = sep }
}

func WithNoGap() FormatOpt { return WithBeforeUnits("") }

func WithBeforeUnits(sep string) FormatOpt {
	return func(f *Formatter) { f.beforeUnits = sep }
}

func WithFmt(tmp string) FormatOpt {
	return func(f *Formatter) { f.valueFmt = tmp }
}

func WithNoFraction() FormatOpt {
	return func(f *Formatter) {
		f.fracFn = func(n, d string) string {
			if d != "" {
				return n + f.unitSep + d
			}
			return n
		}
		f.negativePowers = true
	}
}

func formatUnitSymbol(u Unit) string { return u.Symbol() }
func formatAsciiPower(u string, p int) string {
	if p == 1 {
		return u
	}
	return fmt.Sprintf("%s^%d", u, p)
}

// SUPERSCRIPT ZERO through SUPERSCRIPT NINE
var supers = []rune("\u2070\u00B9\u00B2\u00B3\u2074\u2075\u2076\u2077\u2078\u2079")

const superMinus = '\u207B'

func formatUnicodePower(u string, pow int) string {
	var minus bool
	if pow == 1 {
		return u
	}
	if pow < 0 {
		minus = true
		pow *= -1
	}
	var s [64]rune
	i := len(s)
	for pow >= 10 {
		i--
		s[i] = rune(supers[pow%10])
		pow /= 10
	}
	i--
	s[i] = rune(supers[pow])
	if minus {
		i--
		s[i] = superMinus
	}
	return u + string(s[i:])
}

func WithUnicode() FormatOpt {
	return func(f *Formatter) {
		f.powerFn = formatUnicodePower
		f.unitSep = MiddleDot
		f.fracFn = func(a, b string) string {
			if b != "" {
				return a + FractionSlash + b
			}
			return a
		}
	}
}

func WithMathML() FormatOpt {
	return func(f *Formatter) {
		f.powerFn = func(u string, p int) string {
			if p == 1 {
				return fmt.Sprintf("<mi>%s</mi>", u)
			}
			return fmt.Sprintf("<msup><mi>%s</mi><mn>%d</mn></msup>", u, p)
		}
		f.valueFn = func(tmp string, v float64) string {
			return fmt.Sprintf("<mn>%s</mn>", fmt.Sprintf(tmp, v))
		}
		f.fracFn = func(n, d string) string {
			if d != "" {
				return fmt.Sprintf("<mfrac>%s%s</mfrac>", n, d)
			}
			return n
		}
		f.unitSep = "<mo>" + MiddleDot + "</mo>"
		f.beforeUnitRow = "<mrow>"
		f.afterUnitRow = "</mrow>"
	}
}

func WithLaTeX() FormatOpt {
	return func(f *Formatter) {
		f.unitSep = MiddleDot
		f.fracFn = func(n, d string) string {
			if d != "" {
				return fmt.Sprintf("\frac{%s}{%s}", n, d)
			}
			return n
		}
	}
}

const (
	MiddleDot     = "\u22C5"
	FractionSlash = "\u2044"
	Nbsp          = "\u00A0"
)

func formatValue(tmp string, v float64) string { return fmt.Sprintf(tmp, v) }

var defaultFormatter = Formatter{
	unitFn:  func(u Unit) string { return u.Symbol() },
	powerFn: formatAsciiPower,
	valueFn: formatValue,
	fracFn: func(n, d string) string {
		if d != "" {
			return n + "/" + d
		}
		return n
	},
	negativePowers: false,
	unitSep:        " ",
	beforeUnits:    " ",
	valueFmt:       "%g",
}

var DefaultFormatter = defaultFormatter

func NewFormatter(opts ...FormatOpt) *Formatter {
	f := defaultFormatter
	f.Config(opts...)
	return &f
}

func (f *Formatter) Config(opts ...FormatOpt) {
	for _, opt := range opts {
		opt(f)
	}
}

func (f *Formatter) formatUnitList(us []Unit, mult int) (r []string) {
	if len(us) == 0 {
		return nil
	}
	prev := us[0]
	pow := 1
	write := func() {
		ustr := f.powerFn(f.unitFn(prev), pow*mult)
		r = append(r, ustr)
	}
	for i := 1; i < len(us); i++ {
		u := us[i]
		if u == nil {
			continue
		}
		if u == prev {
			pow++
			continue
		}
		if prev != nil {
			write()
		}
		prev = u
	}
	write()
	return
}

func (f *Formatter) FormatUnits(us Units) string {
	var sb strings.Builder
	if us.N != nil || us.D != nil {
		sb.WriteString(f.beforeUnits)

		var num strings.Builder
		var mult int = 1
		if us.N != nil {
			strs := f.formatUnitList(us.N, mult)
			num.WriteString(f.beforeUnitRow)
			num.WriteString(strings.Join(strs, f.unitSep))
			num.WriteString(f.afterUnitRow)
		}
		var denom strings.Builder
		if us.D != nil {
			if f.negativePowers {
				mult = -1
			}
			strs := f.formatUnitList(us.D, mult)
			denom.WriteString(f.beforeUnitRow)
			denom.WriteString(strings.Join(strs, f.unitSep))
			denom.WriteString(f.afterUnitRow)
		}
		sb.WriteString(f.fracFn(num.String(), denom.String()))
	}
	return sb.String()
}

func (f *Formatter) Format(v Value) string {
	return f.Sprintf(f.valueFmt, v)
}

func (f *Formatter) Sprintf(tmpl string, v Qualified) string {
	var sb strings.Builder
	sb.WriteString(f.valueFn(tmpl, v.Value()))
	u := v.Units()
	if !u.IsEmpty() {
		sb.WriteString(f.FormatUnits(u))
	}
	return sb.String()
}
