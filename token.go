package main

import "fmt"

type (
	TokenType byte
	Token     struct {
		start, end int
		val        string
		t          TokenType
	}
)

const (
	Number        TokenType = 'N'
	LeftParen     TokenType = '('
	RightParen    TokenType = ')'
	Addition      TokenType = '+'
	Subtraction   TokenType = '-'
	Multiplicaton TokenType = '*'
	Division      TokenType = '/'
	Exponent      TokenType = '^'
	Variable      TokenType = 'V'
	Equal         TokenType = '='
	NilTok        TokenType = 0
	EOFTok        TokenType = 1
	ErrTok        TokenType = 2
)

func (t Token) String() string {
	switch t.t {
	case EOFTok:
		return "EOF"
	case ErrTok:
		return fmt.Sprintf("Invalid token at cols: [%d:%d). Token: {{%s}}", t.start, t.end, t.val)
	}
	return fmt.Sprintf("%s : %d", t.val, t.t)
}
