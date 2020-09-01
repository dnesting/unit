package unitdef

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/dnesting/unit/unitdef/scanner"
	"github.com/dnesting/unit/unitdef/token"

	"github.com/dnesting/unit"
)

type parser struct {
	sc   *scanner.Scanner
	defs unit.Registry

	pos scanner.Pos
	tok token.Token
	lit string
}

func From(r io.Reader) (*unit.Registry, error) {
	return parseFrom("-", r)
}

func FromFile(fname string) (*unit.Registry, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	return parseFrom(fname, f)
}

func parseFrom(fname string, r io.Reader) (*unit.Registry, error) {
	p, err := newParser(fname, r)
	if err != nil {
		return nil, err
	}
	if err := p.parse(); err != nil {
		return nil, err
	}
	return &p.defs, nil
}

func newParser(fname string, r io.Reader) (*parser, error) {
	var sc scanner.Scanner
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	sc.Init(fname, data)
	return &parser{sc: &sc}, nil
}

/*
type posError struct {
	fname      string
	line       int
	lineOffset int
	wrapped    error
}

func (p posError) Unwrap() error { return p.wrapped }
func (p posError) Error() string {
	return fmt.Sprintf("%s (at %s:%d)", p.wrapped.Error(), p.fname, p.line, p.lineOffset)
}

func (p *Scanner) error(s string, args ...interface{}) error {
	e := posError{
		fname:      p.fname,
		line:       p.line,
		lineOffset: p.n - p.lineStart,
		wrapped:    fmt.Errorf(s, args...),
	}
	p.errs = append(p.errs, e)
	return e
}
*/

func (p *parser) error(pos scanner.Pos, msg string, args ...interface{}) {
	p.errors.Add(pos, fmt.Sprintf(msg, args...))
}

func (p *parser) next() {
	for {
		p.pos, p.tok, p.lit = p.sc.Scan()
		if p.tok != token.COMMENT {
			break
		}
	}
}

func isLineStart(t token.Token) bool {
	switch t {
	case token.DIRECTIVE:
		return true
	case token.IDENT:
		return true
	}
	return false
}

func eq(wanted ...token.Token) func(token.Token) bool {
	return func(have token.Token) bool {
		for _, t := range wanted {
			if have == t {
				return true
			}
		}
		return false
	}
}

func (p *parser) advance(fn func(token.Token) bool) {
	for ; p.tok != token.EOF; p.next() {
		if !fn(p.tok) {
			return
		}
	}
}

func (p *parser) parse() error {
	for {
	}
	return nil
}
