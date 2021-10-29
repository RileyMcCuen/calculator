package main

import (
	"fmt"
	"strings"
)

type (
	tokenNode struct {
		token      Token
		prev, next *tokenNode
	}
	TokenList struct {
		head, tail *tokenNode
		size       int
	}
)

func newTokenNode(tok Token) *tokenNode {
	return &tokenNode{tok, nil, nil}
}

func (tn *tokenNode) setPrev(newPrev *tokenNode) *tokenNode {
	if tn == nil {
		return nil
	}
	tn.prev = newPrev
	return tn
}

func (tn *tokenNode) setNext(newNext *tokenNode) *tokenNode {
	if tn == nil {
		return nil
	}
	tn.next = newNext
	return tn
}

func NewTokenList() *TokenList {
	return &TokenList{}
}

func (tl *TokenList) String() string {
	sb := strings.Builder{}
	sb.WriteString("[")
	for node := tl.head; node != nil; node = node.next {
		sb.WriteString(fmt.Sprintf("%s, ", node.token))
	}
	sb.WriteString("]")
	return sb.String()
}

func (tl *TokenList) Empty() bool {
	return tl.size == 0
}

// func (tl *TokenList) Enqueue(token Token) *TokenList {
// 	newHead := newTokenNode(token)
// 	tl.size += 1
// 	if tl.size == 0 {
// 		tl.head, tl.tail = newHead, newHead
// 	} else {
// 		tl.head = newHead.setNext(tl.head.setPrev(newHead))
// 	}
// 	return tl
// }

func (tl *TokenList) Push(tok Token) *TokenList {
	newTail := newTokenNode(tok)
	if tl.size == 0 {
		tl.head, tl.tail = newTail, newTail
	} else {
		tl.tail = newTail.setPrev(tl.tail.setNext(newTail))
	}
	tl.size += 1

	return tl
}

func (tl *TokenList) Peek() Token {
	if tl.Empty() {
		return Token{-1, -1, "TokenList.Poll was called when TokenList was empty", ErrTok}
	}

	return tl.tail.token
}

func (tl *TokenList) Poll() Token {
	switch tl.size {
	case 0:
		return Token{-1, -1, "TokenList.Poll was called when TokenList was empty", ErrTok}

	case 1:
		token := tl.tail.token
		tl.head, tl.tail, tl.size = nil, nil, tl.size-1
		return token

	default:
		token := tl.tail.token
		tl.tail, tl.size = tl.tail.prev, tl.size-1
		tl.tail.next.setPrev(nil)
		tl.tail.setNext(nil)
		return token
	}
}
