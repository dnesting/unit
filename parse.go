package unit

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

type parser struct {
	value     string
	n         int
	ch        rune
	eof       bool
	reg       *Registry
	mustExist bool
}

func newParser(v string, reg *Registry, mustExist bool) *parser {
	p := &parser{
		value:     v,
		n:         -1,
		reg:       reg,
		mustExist: mustExist,
	}
	p.next()
	return p
}

func (p *parser) next() {
	if p.n+1 == len(p.value) {
		p.n++
		p.ch = 0
		p.eof = true
	} else if p.n+1 < len(p.value) {
		p.n++
		p.ch = rune(p.value[p.n])
	}
}

func (p *parser) skipSpaces() {
	for unicode.IsSpace(p.ch) {
		p.next()
	}
}

func (p *parser) skipDigits() {
	for unicode.IsDigit(p.ch) {
		p.next()
	}
}

func (p *parser) consume(str string) bool {
	start := p.n
	var i int
	for i < len(str) && p.ch == rune(str[i]) {
		p.next()
	}
	if i == len(str) {
		return true
	}
	p.n = start
	return false
}

func (p *parser) isDigit() bool {
	return unicode.IsDigit(p.ch)
}

func (p *parser) parseInt() (int64, error) {
	start := p.n
	for p.isDigit() {
		p.next()
	}
	return strconv.ParseInt(p.value[start:p.n], 10, 64)
}

func (p *parser) parseFloat() (value float64, ok bool, err error) {
	start := p.n
	var commit bool
	if p.ch == '-' || p.ch == '+' {
		p.next()
		commit = true
	}
	if !p.consume("Inf") && !p.consume("NaN") {
		if commit && !p.isDigit() {
			return 0, false, fmt.Errorf("offset %d: expected number after %q, got %q", p.n, p.value[:p.n], p.value[p.n:])
		}
		p.skipDigits()
		if p.ch == '.' {
			p.next()
			p.skipDigits()
		}
		if p.ch == 'e' {
			p.next()
			if !p.isDigit() {
				return 0, false, fmt.Errorf("offset %d: expected number after %q, got %q", p.n, p.value[:p.n], p.value[p.n:])
			}
			p.skipDigits()
		}
	}
	if start == p.n {
		return 0, false, nil
	}
	value, err = strconv.ParseFloat(p.value[start:p.n], 64)
	if err != nil {
		err = fmt.Errorf("parse float %q: %w", p.value[start:p.n], err)
		return value, false, err
	}
	return value, true, nil
}

func (p *parser) isUnitDiv() bool {
	return p.ch == '/'
}

func (p *parser) parseUnits() (u Units, err error) {
	start := p.n
	u.N, err = p.parseUnitLine()
	if err != nil {
		return Units{}, err
	}
	if p.isUnitDiv() {
		p.next()
		p.skipSpaces()
		denom, err := p.parseUnitLine()
		if err != nil {
			return Units{}, err
		}
		if denom == nil {
			return Units{}, fmt.Errorf("missing units after slash at offset %d", start)
		}
		u.D = denom
	}
	return u, nil
}

func fromSuper(r rune) int {
	for i, o := range supers {
		if r == o {
			return i
		}
	}
	return -1
}

func (p *parser) isExponent() bool {
	return p.ch == '^' || fromSuper(p.ch) >= 0
}

func (p *parser) parseExponent() (exp int64, err error) {
	if p.ch == '^' {
		p.next()
		return p.parseInt()
	}
	for {
		n := fromSuper(p.ch)
		if n >= 0 {
			exp = exp*10 + int64(n)
			p.next()
		} else {
			break
		}
	}
	if exp == 0 {
		return 0, errors.New("exponent must be > 0")
	}
	return exp, nil
}

var (
	unitFirst = []*unicode.RangeTable{unicode.Letter, unicode.So}
	unitAfter = []*unicode.RangeTable{unicode.Letter, unicode.So, unicode.Number}
)

func (p *parser) parseUnitLine() (us []Unit, err error) {
	for unicode.IsOneOf(unitFirst, p.ch) {
		start := p.n
		p.next()
		for unicode.IsOneOf(unitAfter, p.ch) {
			p.next()
		}
		name := p.value[start:p.n]
		if found := p.reg.Find(name); found != nil {
			us = append(us, found.Unit())
		} else if p.mustExist {
			return nil, fmt.Errorf("unknown unit %q", name)
		} else {
			us = append(us, p.reg.Register(Primitive(name)).Unit())
		}
		if p.isExponent() {
			exp, err := p.parseExponent()
			if err != nil {
				return nil, err
			}
			for i := 1; i < int(exp); i++ {
				us = append(us, us[len(us)-1])
			}
		}
		p.skipSpaces()
	}
	return
}

func Parse(str string, reg *Registry, mustExist bool) (val Value, err error) {
	p := newParser(str, reg, mustExist)

	scalar, gotScalar, err := p.parseFloat()
	if err != nil {
		return Value{}, fmt.Errorf("units parse: %w", err)
	}
	if !gotScalar {
		scalar = 1
	}
	p.skipSpaces()

	units, err := p.parseUnits()
	if err != nil {
		return Value{}, fmt.Errorf("units parse: %w", err)
	}
	p.skipSpaces()
	if !p.eof {
		return Value{}, fmt.Errorf("units parse: extra text after units at %q here-> %q", p.value[:p.n], p.value[p.n:])
	}
	return Value{S: scalar, U: units}, nil
}
