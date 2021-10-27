package main

type (
	ErrorStatus int
)

// common ascii values
const (
	a    = 97
	z    = 122
	A    = 65
	Z    = 90
	zero = 48
	nine = 57
	NULL = 0
	dot  = 46
)

// common strings
const (
	space      = " "
	leftParen  = "("
	rightParen = ")"
	equal      = "="
	minus      = "-"
	plus       = "+"
	multiply   = "*"
	divide     = "/"
	exponent   = "^"
)

const (
	invalid ErrorStatus = iota
	valid
	end
)

func isLowercase(c uint8) bool {
	return c >= a && c <= z
}

func isUppercase(c uint8) bool {
	return c >= A && c <= Z
}

func isLetter(c uint8) bool {
	return (c >= A && c <= Z) || (c >= a && c <= z)
}

func isNumeric(c uint8) bool {
	return c == dot || isDigit(c)
}

func isDigit(c uint8) bool {
	return c >= zero && c <= nine
}

func isNonZeroDigit(c uint8) bool {
	return c > zero && c <= nine
}
