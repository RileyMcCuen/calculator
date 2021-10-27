package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

type (
	Node interface {
		Operator() (Operator, error)
		Eval() (float64, error)
		String() string
	}
	Value    string
	Operator byte
	Var      struct {
		name   string
		values map[string]float64
	}
	BinaryOperation struct {
		left, op, right Node
	}
	NodeList struct {
		nodes []Node
	}
	MapEvaluator struct {
		name   string
		node   Node
		values map[string]float64
	}

	Parser struct {
		in     chan Token
		cur    Token
		values map[string]float64
	}
)

const (
	NilOp  Operator = 0
	AddOp           = Operator(Addition)
	SubOp           = Operator(Subtraction)
	MultOp          = Operator(Multiplicaton)
	DivOp           = Operator(Division)
	ExpOp           = Operator(Exponent)
)

var (
	Operators = []Operator{ExpOp, DivOp, MultOp, SubOp, AddOp}
)

func NewParser(c chan Token, values map[string]float64) *Parser {
	return &Parser{
		c,
		<-c,
		values,
	}
}

func (p *Parser) Parse() (Node, error) {
	root, err := p.optionalAssignment()
	if err != nil {
		return nil, err
	}
	return root, nil
}

func (p *Parser) next() *Parser {
	p.cur = <-p.in
	return p
}

func (p *Parser) match(ts ...TokenType) error {
	for _, t := range ts {
		if p.cur.t == t {
			p.next()
		} else {
			return errors.New("unexpected token: [" + p.cur.val + "]")
		}
	}
	return nil
}

func (p *Parser) optionalAssignment() (Node, error) {
	// <variable> <=> <expression> | <expression>

	if firstToken := p.cur; firstToken.t == Variable {
		// might be an assignment

		switch {
		case p.next().cur.t == Equal:
			// this is an assignment

			val, err := p.next().expression()
			if err != nil {
				return nil, err
			}

			return MapEvaluator{firstToken.val, val, p.values}, nil

		case p.cur.t == EOFTok:
			return Var{firstToken.val, p.values}, nil

		case p.cur.t == ErrTok:
			return nil, fmt.Errorf("encountered an invalid sequence in the input: '%s'", p.cur.val)

		default:
			// essentially same as parsing parenthetic expression, but we already pulled out the first token and cannot put it back
			nl, err := p.parenthetic(Var{firstToken.val, p.values})
			if err != nil {
				return nil, err
			}

			return nl.pemdas()

		}
	} else {
		// not an assignment, treat as regular expression
		return p.parenthetic()
	}
}

func (p *Parser) parenthetic(nodes ...Node) (*NodeList, error) {
	nl := &NodeList{}

	for _, node := range nodes {
		nl.Append(node)
	}

	for i := len(nodes); p.cur.t != EOFTok && p.cur.t != ErrTok && p.cur.t != RightParen; i++ {
		if i%2 == 0 {
			node, err := p.expression()
			if err != nil {
				return nil, err
			}

			nl.Append(node)

		} else {
			node, err := p.operator()
			if err != nil {
				return nil, err
			}

			nl.Append(node)
		}
	}

	if p.cur.t == ErrTok {
		return nil, fmt.Errorf("encountered an invalid sequence in the input: '%s'", p.cur.val)
	} else if len(nl.nodes) == 0 {
		return nil, errors.New("nodelist has no elements, therefore an invalid top level or parenthetic expression has been encountered")
	} else if _, err := nl.nodes[len(nl.nodes)-1].Operator(); err == nil || len(nl.nodes)%2 == 0 {
		fmt.Println(nl)
		return nil, errors.New("nodelist ended in an operator, but must end in an expression to be valid: " + err.Error())
	}

	return nl, nil
}

func (p *Parser) expression() (Node, error) {
	// ( <expression> ) | variable | number
	switch token := p.cur; token.t {
	case LeftParen:
		p.next()

		nl, err := p.parenthetic()
		if err != nil {
			return nil, err
		}

		if p.match(RightParen) != nil {
			return nil, errors.New("parenthetic expression did not end in right paren")
		}

		return nl.pemdas()

	case Variable:
		p.next()
		return Var{token.val, p.values}, nil

	case Number:
		p.next()
		return Value(token.val), nil

	default:
		return nil, fmt.Errorf("invalid token '%c' found when parsing expression", token.t)

	}
}

func (p *Parser) operator() (Node, error) {
	// + | - | * | / | ^
	node, err := Op(p.cur.t)
	if err != nil {
		return nil, err
	}

	p.next()

	return node, nil
}

func (v Value) String() string {
	switch TokenType(v[0]) {
	case 'p':
		return "PI"
	case 'e':
		return "e"
	default:
		return string(v)
	}
}

func (v Value) Eval() (float64, error) {
	switch TokenType(v[0]) {
	case 'p':
		return 3.1415927, nil
	case 'e':
		return 2.7182818, nil
	default:
		return strconv.ParseFloat(string(v), 64)
	}
}

func (Value) Operator() (Operator, error) {
	return NilOp, errors.New("Value cannot be used as an operator")
}

func Op(t TokenType) (Operator, error) {
	switch t {
	case Addition, Subtraction, Multiplicaton, Division, Exponent:
		return Operator(t), nil
	default:
		return 0, fmt.Errorf("cannot use TokenType '%c' as Operator", t)
	}
}

func (o Operator) String() string {
	return fmt.Sprintf("%c", o)
}

func (o Operator) Eval() (float64, error) {
	return 0, errors.New("Operator cannot be used as an evaluable")
}

func (o Operator) Operator() (Operator, error) {
	return o, nil
}

func (v Var) String() string {
	val, ok := v.values[v.name]
	if !ok {
		return fmt.Sprintf("{%s:[UNDEFINED_VAR]}", v.name)
	}
	return fmt.Sprintf("{%s:%f}", v.name, val)
}

func (v Var) Eval() (float64, error) {
	val, ok := v.values[v.name]
	if !ok {
		return 0, fmt.Errorf("found undefined variable '%s', all variables must be defined before being used in expressions", v.name)
	}
	return val, nil
}

func (Var) Operator() (Operator, error) {
	return NilOp, errors.New("Var cannot be used as an operator")
}

func (m MapEvaluator) String() string {
	return fmt.Sprintf("{ME-%p|%s:%s}", &m.values, m.name, m.node.String())
}

func (m MapEvaluator) Eval() (float64, error) {
	val, err := m.node.Eval()
	if err != nil {
		return 0, err
	}
	m.values[m.name] = val
	return val, nil
}

func (MapEvaluator) Operator() (Operator, error) {
	return NilOp, errors.New("MapEvaluator cannot be used as an operator")
}

func (bo BinaryOperation) String() string {
	return fmt.Sprintf("[%s %c %s]", bo.left.String(), bo.op, bo.right.String())
}

func (bo BinaryOperation) Eval() (float64, error) {
	left, lErr := bo.left.Eval()
	if lErr != nil {
		return 0, lErr
	}

	right, rErr := bo.right.Eval()
	if rErr != nil {
		return 0, rErr
	}

	op, err := bo.op.Operator()
	if err != nil {
		return 0, err
	}

	switch op {
	case AddOp:
		return left + right, nil
	case SubOp:
		return left - right, nil
	case MultOp:
		return left * right, nil
	case DivOp:
		return left / right, nil
	case ExpOp:
		return math.Pow(left, right), nil
	default:
		return 0, fmt.Errorf("invalid operator found: '%c'", bo.op)
	}
}

func (BinaryOperation) Operator() (Operator, error) {
	return NilOp, errors.New("BinaryOperator cannot be used as an operator")
}

func findOpIndex(nodes []Node, start int, target Operator) int {
	for i, node := range nodes[start:] {
		if op, err := node.Operator(); err == nil && op == target {
			return i + start
		}
	}
	return -1
}

func (l *NodeList) String() string {
	if len(l.nodes) == 0 {
		return ("(EMPTY)")
	}
	ret := "(" + l.nodes[0].String()
	for _, node := range l.nodes {
		ret += fmt.Sprintf(" %s", node.String())
	}
	return ret + ")"
}

func (l *NodeList) Append(a Node) *NodeList {
	l.nodes = append(l.nodes, a)
	return l
}

func (l *NodeList) pemdas() (Node, error) {
	// special cases for list
	if len(l.nodes) == 0 {
		return nil, errors.New("nodelist was simplified while empty, therefore it was not a valid expression")
	} else if len(l.nodes) == 1 {
		return l.nodes[0], nil
	}

	for _, op := range Operators {
		start := 1
		for i := findOpIndex(l.nodes, start, op); i != -1; {
			start = i

			l.nodes[i-1] = BinaryOperation{l.nodes[i-1], l.nodes[i], l.nodes[i+1]}

			copy(l.nodes[i:], l.nodes[i+2:])
			l.nodes = l.nodes[:len(l.nodes)-2]

			i = findOpIndex(l.nodes, start, op)
		}
	}

	return l.nodes[0], nil
}

func (l *NodeList) Eval() (float64, error) {
	n, err := l.pemdas()
	if err != nil {
		return 0, err
	}
	return n.Eval()
}

func (*NodeList) Operator() (Operator, error) {
	return NilOp, errors.New("*NodeList cannot be used as an operator")
}
