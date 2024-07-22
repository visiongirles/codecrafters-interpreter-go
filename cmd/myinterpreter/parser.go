package main

import "fmt"

type Parser struct {
	tokens  []Token
	current int
}

type NodeType string

// https://gobyexample.com/interfaces
type ASTNode interface {
	String() string
	// Type()
}

type Expression struct {
}

type LiteralExpression struct {
	value Primitive
}

func (l LiteralExpression) String() string {
	return l.value.String()
}

type NumberExpression struct {
	token Token
}

func (l NumberExpression) String() string {
	return l.token.literal
}

type Primitive int

const (
	nilPrimitive Primitive = iota
	truePrimitive
	falsePrimitive
)

func (p Primitive) String() string {
	return [...]string{
		"nil",
		"true",
		"false",
	}[p]
}

func initParser() Parser {
	return Parser{
		current: 0,
	}
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens)
}

func (p *Parser) parseTokens() ASTNode {
	if !p.isAtEnd() {
		return p.parseToken()
	}
	return nil
}

func (p *Parser) parseToken() ASTNode {
	token := p.advance()
	// fmt.Printf("[Parser parseToken()]. Token %s has not implemented\n", token.typeToken.String())
	switch token.typeToken {
	case TRUE:
		return LiteralExpression{value: truePrimitive}
	case FALSE:
		return LiteralExpression{value: falsePrimitive}
	case NIL:
		return LiteralExpression{value: nilPrimitive}
	case NUMBER:
		return NumberExpression{token: token}
	default:
		fmt.Printf("[Parser parseToken()] default. Token %s has not implemented", token.lexeme)
	}
	return nil
}

func (p *Parser) advance() Token {
	index := p.current
	p.current++
	return p.tokens[index]
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

// func (p *Parser) previous() Token {
// 	return p.tokens[p.current-1]
// }

func Parse(tokens []Token) (ASTNode, string) {
	parser := initParser()
	parser.tokens = tokens
	expression := parser.parseTokens()
	return expression, ""
}

// type BinaryExpression struct {
// 	Expression
// 	left     Expression
// 	operator Token
// 	right    Expression
// }

// type UnaryExpression struct {
// 	Expression
// 	operator Token
// 	terminal Primitive
// }

func (p *Parser) match(tokensTypes []TokenType) bool {
	for _, tokenType := range tokensTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().typeToken == tokenType
}
