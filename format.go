package unit

import (
	"fmt"
	"runtime/debug"
	"strings"
)

// Formatter allows for formatting qualified values and units.  It's configured
// with one or more FormatOpts.  A DefaultFormatter is provided which formats
// units according to a standard "1.234 kg m/s^2" style.  Formatter values should
// be created with NewFormatter.
type Formatter struct {
	unitFn         func(u Unit) string
	powerFn        func(u string, p int) string
	valueFn        func(tmpl string, v float64) string
	fracFn         func(num, denom string) string
	negativePowers bool
	valueFmt       string
	noGapFor       []Units

	beforeUnits   string
	beforeUnitRow string
	unitSep       string
	afterUnitRow  string
	afterUnits    string
}

// FormatOpt is an option passed to NewFormatter or Formatter.Config which
// configures the Formatter.
type FormatOpt func(f *Formatter)

/*
// WithUnitFunc provides a way to customize the string value of a Unit.
// Fn will be called at least once for every unit to be formatted, and
// the implementation is expected to return a string representation of
// the unit.  By default, the string returned from the Unit's Symbol()
// method will be used.
func WithUnitFunc(fn func(u Unit) string) FormatOpt {
	return func(f *Formatter) { f.unitFn = fn }
}
*/

/*
// WithPowerFunc provides a way to customize the rendering of units.
// Fn will be called at least once for every unit to be formatted, along
// with the power associated with the unit.  U will be set to the return
// value of the function specified in WithUnitFunc, or the return value
// of the Unit's Symbol() method.  P will be the power associated with
// the Unit.  This may be 1 but will never be 0.  If WithNegativeExponent
// is specified as an option, exponents will be negative to designate
// they are in the denominator.  By default, this will produce strings
// of the form "u", "u^2", or "u^-1".
func WithPowerFunc(fn func(u string, p int) string) FormatOpt {
	return func(f *Formatter) { f.powerFn = fn }
}
*/

// WithFraction specifies a string to use to separate the numerator and
// denominator of a value's units.  It will not be used if the
// denominator is empty.  By default, "/" will be used.
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

// WithUnitSep specifies a string to use to separate units.  By default,
// a space will be used.
func WithUnitSep(sep string) FormatOpt {
	return func(f *Formatter) { f.unitSep = sep }
}

// WithNoGap eliminates the default space between a value's scalar
// component and its units.
func WithNoGap() FormatOpt {
	return func(f *Formatter) { f.beforeUnits = "" }
}

func WithNoGapFor(ms ...Maker) FormatOpt {
	return func(f *Formatter) {
		for _, m := range ms {
			f.noGapFor = append(f.noGapFor, m.Units())
		}
	}
}

/*
// WithBeforeUnits specifies a string to use to separate a qualified
// value's scalar component and its units.  It will not be used if the
// value has no units.  By default, a space will be used.
func WithBeforeUnits(sep string) FormatOpt {
	return func(f *Formatter) { f.beforeUnits = sep }
}
*/

// WithFmt specifies the default Sprintf-style format to use for the qualified
// value's scalar component.  By default, "%g" is used.
func WithFmt(tmp string) FormatOpt {
	return func(f *Formatter) { f.valueFmt = tmp }
}

// WithNoFraction specifies that the units should not be rendered as
// a fraction.  Units in the denominator will be rendered with a negative
// exponent instead.
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

// WithUnicode uses a DotOperator as a unit separator, FractionSlash to
// separate the numerator and denominator of a unit, and uses superscript
// numbers to format the exponents.
func WithUnicode() FormatOpt {
	return func(f *Formatter) {
		f.powerFn = formatUnicodePower
		f.unitSep = string(DotOperator)
		f.fracFn = func(a, b string) string {
			if b != "" {
				return a + string(FractionSlash) + b
			}
			return a
		}
	}
}

// WithMathML specifies that MathML tags should be used to lay out components
// of the qualified value.  Enclosing <math> tags are not generated.
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
		f.unitSep = "<mo>" + string(DotOperator) + "</mo>"
		f.beforeUnitRow = "<mrow>"
		f.afterUnitRow = "</mrow>"
	}
}

// WithLaTeX specifies that LaTeX markup should be used to lay out components
// of the qualified value.
func WithLaTeX() FormatOpt {
	return func(f *Formatter) {
		f.unitSep = string(DotOperator)
		f.fracFn = func(n, d string) string {
			if d != "" {
				return fmt.Sprintf("\frac{%s}{%s}", n, d)
			}
			return n
		}
	}
}

const (
	DotOperator   = '\u22C5'
	FractionSlash = '\u2044'
	Nbsp          = '\u00A0'
)

func defaultFormatValue(tmp string, v float64) string { return fmt.Sprintf(tmp, v) }
func defaultFormatUnit(u Unit) string                 { return u.Symbol() }
func defaultFormatPower(u string, p int) string {
	if p == 1 {
		return u
	}
	return fmt.Sprintf("%s^%d", u, p)
}
func defaultFormatFraction(n, d string) string {
	if d != "" {
		return n + "/" + d
	}
	return n
}

var defaultFormatter = Formatter{
	unitFn:         defaultFormatUnit,
	powerFn:        defaultFormatPower,
	valueFn:        defaultFormatValue,
	fracFn:         defaultFormatFraction,
	negativePowers: false,
	unitSep:        " ",
	beforeUnits:    " ",
	valueFmt:       "%g",
}

// DefaultFormatter is the Formatter that is used by Value.String() and
// Units.String() to specify default formatting configuration.  It may be
// configured by passing FormatOpts to its Config method.
var DefaultFormatter = defaultFormatter

// NewFormatter creates a new Formatter configured by the given opts.
func NewFormatter(opts ...FormatOpt) *Formatter {
	f := defaultFormatter
	f.Config(opts...)
	return &f
}

// Config configures the Formatter through the provideed opts.
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
	add := func() {
		ustr := f.powerFn(f.unitFn(prev), pow*mult)
		r = append(r, ustr)
	}
	for i := 1; i < len(us); i++ {
		if us[i] == nil {
			continue
		}
		if us[i].Equal(prev) {
			pow++
			continue
		}
		if prev != nil {
			add()
		}
		prev = us[i]
	}
	add()
	return
}

// FormatUnits formats the Units according the Formatter's configuration.
func (f *Formatter) FormatUnits(us Units) string {
	var sb strings.Builder
	if us.N != nil || us.D != nil {
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

// Format formats v's scalar value and units according to the Formatter's
// configuration.
func (f *Formatter) Format(v Value) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
			panic(r)
		}
	}()
	return f.Sprintf(f.valueFmt, v)
}

func contains(u Units, list []Units) bool {
	for _, x := range list {
		if u.Equal(x) {
			return true
		}
	}
	return false
}

// Sprintf formats the qualified value's scalar value using the given
// Sprintf-style formatting template, and then formats the units according
// to the Formatter's configuration.
func (f *Formatter) Sprintf(tmpl string, v Qualified) string {
	var sb strings.Builder
	sb.WriteString(f.valueFn(tmpl, v.Value()))
	u := v.Units()
	if !u.Empty() {
		if !contains(u, f.noGapFor) {
			sb.WriteString(f.beforeUnits)
		}
		sb.WriteString(f.FormatUnits(u))
	}
	return sb.String()
}
