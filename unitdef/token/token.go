package token

import (
	"strconv"
	"unicode"
)

type Token int

const (
	ILLEGAL Token = iota
	EOF
	COMMENT
	IDENT
	NUMBER
	DIRECTIVE
	PRIM
	PRIM_DIMENSIONLESS

	BACKSLASH
	CARET
	COMMA
	EQUAL
	LBRACKET
	RBRACKET
	LPAREN
	RPAREN
	MINUS
	PIPE
	PLUS
	SEMICOLON
	SLASH
	NEWLINE

	keyword_beg
	DOMAIN
	RANGE
	UNITS
	keyword_end
)

var tokens = [...]string{
	EOF:                "EOF",
	COMMENT:            "COMMENT",
	IDENT:              "IDENT",
	NUMBER:             "NUMBER",
	DIRECTIVE:          "DIRECTIVE",
	PRIM:               "!",
	PRIM_DIMENSIONLESS: "!dimensionless",

	BACKSLASH: "\\",
	CARET:     "^",
	COMMA:     ",",
	MINUS:     "-",
	EQUAL:     "=",
	LBRACKET:  "[",
	RBRACKET:  "]",
	LPAREN:    "(",
	RPAREN:    ")",
	PIPE:      "|",
	PLUS:      "+",
	SEMICOLON: ";",
	SLASH:     "/",
	NEWLINE:   "\\n",

	DOMAIN: "domain",
	RANGE:  "range",
	UNITS:  "units",
}

func (tok Token) String() string {
	var s string
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

func Lookup(ident string) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	return IDENT
}

func IsKeyword(name string) bool {
	_, ok := keywords[name]
	return ok
}

func IsIdentifier(name string) bool {
	for i, c := range name {
		if !unicode.IsLetter(c) && c != '_' && (i == 0 || !unicode.IsDigit(c)) {
			return false
		}
	}
	return name != "" && !IsKeyword(name)
}
