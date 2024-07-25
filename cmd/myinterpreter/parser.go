package main

import (
	"fmt"
	"os"
)

type Parser struct {
	tokens  []Token
	current int
	stack   Stack
}

type NodeType string

// https://gobyexample.com/interfaces
type ASTNode interface {
	String() string
	// Type()
}

type LiteralExpression struct {
	value Primitive
}

func (n LiteralExpression) String() string {
	return n.value.String()
}

type NumberExpression struct {
	token Token
}

func (n NumberExpression) String() string {
	return n.token.literal
}

type StringExpression struct {
	token Token
}

func (n StringExpression) String() string {
	return n.token.literal
}

type TrueExpression struct {
	token Token
}

func (n TrueExpression) String() string {
	return n.token.lexeme
}

type FalseExpression struct {
	token Token
}

func (n FalseExpression) String() string {
	return n.token.lexeme
}

type NilExpression struct {
	token Token
}

func (n NilExpression) String() string {
	return n.token.lexeme
}

type GroupExpression struct {
	expression ASTNode
}

func (n GroupExpression) String() string {
	return "(group " + n.expression.String() + ")"
}

type UnaryExpression struct {
	operator   TokenType
	expression ASTNode
}

func (n UnaryExpression) String() string {
	if n.operator == BANG {
		return "(! " + n.expression.String() + ")"
	} else {
		return "(- " + n.expression.String() + ")"

	}
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

func (p *Parser) parseTokens() (ASTNode, string) {
	var node ASTNode
	err := ""

	fmt.Fprintln(os.Stderr, "[DEBUG] Tokens: ", p.tokens)

	for !p.isAtEnd() {
		node, err = p.parseToken()
	}

	if !p.stack.IsEmpty() {
		return node, "[p.stack.IsEmpty() ]Error: Unmatched parentheses."
	}
	return node, err
}

func (p *Parser) parseGroup() (ASTNode, string) {

	token := p.peek()
	if token.typeToken == RIGHT_PAREN {
		p.advance()
		return nil, "Error: Empty expression."
	}

	expr, err := p.parseToken()
	if p.isAtEnd() {
		err += "Error: Unmatched parentheses."
		return nil, err
	}
	token = p.peek()
	if token.typeToken == RIGHT_PAREN {
		if p.stack.Pop() != '(' {
			err += "[p.stack.Pop() != '(']Error: Unmatched parentheses."
			return nil, err
		}
	}
	p.advance()

	return expr, err

}

func (p *Parser) parseToken() (ASTNode, string) {
	token := p.peek()
	p.advance()
	fmt.Fprintf(os.Stderr, "[DEBUG] [parseToken()]. Token %s\n", token.lexeme)
	switch token.typeToken {
	case TRUE:
		return TrueExpression{token: token}, ""
	case FALSE:
		return FalseExpression{token: token}, ""
	case NIL:
		return NilExpression{token: token}, ""
	case NUMBER:
		return NumberExpression{token: token}, ""
	case STRING:
		return StringExpression{token: token}, ""
	case LEFT_PAREN:
		p.stack.Push('(')
		expr, err := p.parseGroup()
		if expr != nil {
			return GroupExpression{expression: expr}, err
		} else {
			return nil, err
		}
		// return nil, ""
	case RIGHT_PAREN:
		return nil, "[RIGHT PAREN]Error: Unmatched parentheses."

	case BANG, MINUS:
		expr, err := p.parseUnary()
		return UnaryExpression{expression: expr, operator: token.typeToken}, err
	default:
		fmt.Fprintf(os.Stderr, "[DEBUG] [Parser parseToken()] default. Token %s has not implemented", token.lexeme)
	}
	return nil, ""
}

func (p *Parser) parseUnary() (ASTNode, string) {
	var expr ASTNode
	var err string

	expr, err = p.parseToken()

	return expr, err
}

func (p *Parser) advance() {
	p.current = p.current + 1
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

// func (p *Parser) previous() Token {
// 	return p.tokens[p.current-1]
// }

func Parse(tokens []Token) (ASTNode, string) {
	parser := initParser()
	parser.tokens = tokens
	expression, err := parser.parseTokens()
	return expression, err
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

func (p *Parser) match(tokensTypes ...TokenType) bool {
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

// func (p *Parser) consume(tokenType TokenType, message string) {
// 	if p.check(tokenType) {
// 		p.advance()
// 		return
// 	}
// }

// func (p *Parser) primary() (ASTNode, string) {
// 	token := p.advance()
// 	if p.match(FALSE) {
// 		return FalseExpression{token: token}, ""
// 	}
// 	if p.match(TRUE) {
// 		return TrueExpression{token: token}, ""
// 	}
// 	if p.match(NIL) {
// 		return NilExpression{token: token}, ""
// 	}
// 	if p.match(NUMBER) {
// 		return NumberExpression{token: token}, ""
// 	}
// 	if p.match(STRING) {
// 		return StringExpression{token: token}, ""
// 	}

// 	if p.match(LEFT_PAREN) {
// 		expr := expression()
// 		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
// 		return GroupExpression{expression: expr}, err
// 	}
// }
