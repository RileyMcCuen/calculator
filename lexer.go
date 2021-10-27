package main

import (
	"strings"
)

type (
	StateFN func() StateFN
	Lexer   struct {
		out        chan Token
		input      string
		start, pos int
	}
)

func NewLexer(input string, tokenStreamSize int) *Lexer {
	ret := new(Lexer)
	ret.out = make(chan Token, tokenStreamSize)
	ret.input = input
	return ret
}

func (l *Lexer) Run() *Lexer {
	go func() {
		for state := l.lexInit; state != nil; {
			state = state()
		}
		close(l.out)
	}()

	return l
}

func (l *Lexer) Out() chan Token {
	return l.out
}

func (l *Lexer) prefixed(pre string) bool {
	return strings.HasPrefix(l.input[l.start:], pre)
}

func (l *Lexer) prefixedOneOf(pres ...string) bool {
	for _, pre := range pres {
		if l.prefixed(pre) {
			return true
		}
	}
	return false
}

func (l *Lexer) curChar() uint8 {
	if l.pos >= len(l.input) {
		return NULL
	}
	return l.input[l.pos]
}

func (l *Lexer) eof() bool {
	return l.curChar() == NULL
}

func (l *Lexer) skip(size string) *Lexer {
	l.pos += len(size)
	l.start = l.pos
	return l
}

func (l *Lexer) step(size string) *Lexer {
	l.pos += len(size)
	return l
}

func (l *Lexer) stepOne() *Lexer {
	l.pos += 1
	return l
}

func (l *Lexer) backup() *Lexer {
	l.pos -= 1
	return l
}

func (l *Lexer) lexInit() StateFN {
	switch {
	case l.eof():
		l.emit(EOFTok)
		return nil
	case l.prefixed(space):
		return l.skip(space).lexInit()
	case l.prefixed(leftParen):
		return l.step(leftParen).emit(LeftParen).lexInit()
	case l.prefixed(rightParen):
		return l.step(rightParen).emit(RightParen).lexInit()
	case l.prefixed(equal):
		return l.step(equal).emit(Equal).lexInit()
	case l.prefixed(plus):
		return l.step(plus).emit(Addition).lexInit()
	case l.prefixed(minus):
		return l.step(minus).emit(Subtraction).lexInit()
	case l.prefixed(multiply):
		return l.step(multiply).emit(Multiplicaton).lexInit()
	case l.prefixed(divide):
		return l.step(divide).emit(Division).lexInit()
	case l.prefixed(exponent):
		return l.step(exponent).emit(Exponent).lexInit()
	case isNumeric(l.curChar()):
		return l.stepOne().lexNum()
	case isLetter(l.curChar()):
		return l.stepOne().lexVariable()
	default:
		return l.lexError()
	}
}

func (l *Lexer) validSeq() ErrorStatus {
	switch {
	case l.eof():
		return end
	case l.prefixedOneOf(space, leftParen, rightParen, equal, minus, plus, multiply, divide, exponent),
		isNumeric(l.curChar()), isLetter(l.curChar()):
		return valid
	default:
		return invalid
	}
}

func (l *Lexer) lexError() StateFN {
	for {
		switch l.validSeq() {
		case end:
			if l.start != l.pos {
				l.emit(ErrTok)
			} else {
				l.emit(EOFTok)
			}
			return nil
		case invalid:
			// this finds the length of the invalid string sequence
			l.stepOne()
		case valid:
			// this can only happen if lexError has already seen an invalid sequence at least once
			// find the entire invalid portion and return it for good error reporting
			l.emit(ErrTok)
			return nil
		}
	}
}

func (l *Lexer) lexNum() StateFN {
	for char := l.curChar(); char != NULL && isNumeric(char); {
		char = l.stepOne().curChar()
	}
	return l.emit(Number).lexInit()
}

func (l *Lexer) lexVariable() StateFN {
	for char := l.curChar(); char != NULL && isLetter(char); {
		char = l.stepOne().curChar()
	}
	return l.emit(Variable).lexInit()
}

func (l *Lexer) send(t Token) *Lexer {
	l.out <- t
	l.start = l.pos
	return l
}

func (l *Lexer) emit(t TokenType) *Lexer {
	return l.send(Token{
		start: l.start,
		end:   l.pos,
		val:   l.input[l.start:l.pos],
		t:     t,
	})
}
