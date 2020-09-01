package scanner

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/dnesting/unit/unitdef/token"
)

type Scanner struct {
	src []byte

	ch       rune
	offset   int
	rdOffset int

	fname      string
	line       int
	lineOffset int
	errs       []error
}

func (s *Scanner) Init(fname string, src []byte) {
	s.fname = fname
	s.src = src
	s.ch = ' '
	s.offset = 0
	s.rdOffset = 0
	s.lineOffset = 0
	s.errs = nil
	s.next()
	if s.ch == bom {
		s.next()
	}
}

const bom = 0xFEFF

func (s *Scanner) next() {
	if s.rdOffset < len(s.src) {
		s.offset = s.rdOffset
		if s.ch == '\n' {
			s.line++
			s.lineOffset = s.offset
		}
		r, w := utf8.DecodeRune(s.src[s.rdOffset:])
		if r == utf8.RuneError && w == 1 {
			s.error("illegal UTF-8 encoding")
		} else if r == bom && s.offset > 0 {
			s.error("illegal byte order mark")
		}
		s.rdOffset += w
		s.ch = r
	} else {
		s.offset = len(s.src)
		if s.ch == '\n' {
			s.lineOffset = s.offset
			s.line++
		}
		s.ch = -1 // eof
	}
}

func (s *Scanner) peek() byte {
	if s.rdOffset < len(s.src) {
		return s.src[s.rdOffset]
	}
	return 0
}

func (s *Scanner) pos() Pos {
	return Pos{s.fname, s.line, s.offset - s.lineOffset}
}

type posError struct {
	Pos
	wrapped error
}

func (p posError) Unwrap() error { return p.wrapped }
func (p posError) Error() string {
	return fmt.Sprintf("%s (at %s)", p.wrapped.Error(), p.Pos)
}

func (p *Scanner) error(s string, args ...interface{}) error {
	e := posError{
		Pos:     p.pos(),
		wrapped: fmt.Errorf(s, args...),
	}
	p.errs = append(p.errs, e)
	return e
}

func (p *Scanner) skipSpaces() {
	for unicode.Is(unicode.Zs, p.ch) {
		p.next()
	}
}

type Pos struct {
	File string
	Line int
	Col  int
}

func (p Pos) String() string {
	return fmt.Sprintf("%s:%d:%d", p.File, p.Line, p.Col)
}

func lower(ch rune) rune     { return ('a' - 'A') | ch } // returns lower-case ch iff ch is ASCII letter
func isDecimal(ch rune) bool { return '0' <= ch && ch <= '9' }
func isLetter(ch rune) bool {
	return 'a' <= lower(ch) && lower(ch) <= 'z' || ch == '_' || ch >= utf8.RuneSelf && unicode.IsLetter(ch)
}
func isDigit(ch rune) bool {
	return isDecimal(ch) || ch >= utf8.RuneSelf && unicode.IsDigit(ch)
}

func (s *Scanner) scanIdentifier() string {
	offs := s.offset
	for isLetter(s.ch) || isDigit(s.ch) || s.ch == '.' {
		s.next()
	}
	return string(s.src[offs:s.offset])
}

func (s *Scanner) digits() {
	for isDecimal(s.ch) {
		s.next()
	}
}

func (s *Scanner) scanNumber() string {
	offs := s.offset
	s.digits()
	if s.ch == '.' {
		s.next()
		s.digits()
	}
	if lower(s.ch) == 'e' {
		s.next()
		if s.ch == '+' || s.ch == '-' {
			s.next()
		}
		s.digits()
	}
	lit := string(s.src[offs:s.offset])
	return lit
}

func (s *Scanner) scanComment() string {
	offs := s.offset
	for s.ch != 0 && s.ch != '\n' {
		s.next()
	}
	return string(s.src[offs:s.offset])
}

func (s *Scanner) scanDirective() string {
	offs := s.offset
	for isLetter(s.ch) || isDigit(s.ch) {
		s.next()
	}
	return string(s.src[offs:s.offset])
}

func (s *Scanner) Scan() (pos Pos, tok token.Token, lit string) {
	s.skipSpaces()
	pos = s.pos()

	switch ch := s.ch; {
	case isLetter(ch):
		lit = s.scanIdentifier()
		tok = token.Lookup(lit)
	case isDecimal(ch) || ch == '.' && isDecimal(rune(s.peek())):
		tok = token.NUMBER
		lit = s.scanNumber()
	default:
		s.next()
		switch ch {
		case -1:
			tok = token.EOF
		case '!':
			tok = token.DIRECTIVE
			lit = s.scanDirective()
		case '#':
			tok = token.COMMENT
			lit = s.scanComment()
		case '^':
			tok = token.CARET
		case '\\':
			tok = token.BACKSLASH
		case ',':
			tok = token.COMMA
		case '[':
			tok = token.LBRACKET
		case ']':
			tok = token.RBRACKET
		case '(':
			tok = token.LPAREN
		case ')':
			tok = token.RPAREN
		case '-':
			tok = token.MINUS
		case '|':
			tok = token.PIPE
		case '+':
			tok = token.PLUS
		case ';':
			tok = token.SEMICOLON
		case '/':
			tok = token.SLASH
		case '\n':
			tok = token.NEWLINE
		default:
			tok = token.ILLEGAL
			lit = string(ch)
		}
	}
	return
}
