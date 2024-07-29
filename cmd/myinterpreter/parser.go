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
	operator   Token
	expression ASTNode
}

type BinaryExpression struct {
	left     ASTNode
	operator Token
	right    ASTNode
}

func (n BinaryExpression) String() string {
	return "(" + n.operator.lexeme + " " + n.left.String() + " " + n.right.String() + ")"
}

func (n UnaryExpression) String() string {
	return "(" + n.operator.lexeme + " " + n.expression.String() + ")"
}

type Primitive int

// const (
// 	nilPrimitive Primitive = iota
// 	truePrimitive
// 	falsePrimitive
// )

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
	var left ASTNode
	errMain := ""

	for !p.isAtEnd() {
		result, err := p.parseToken()
		if result == nil {
			//TODO:
			// if p.tokens[p.current-1].typeToken == RIGHT_PAREN {
			if !p.isAtEnd() && p.tokens[p.current].typeToken == RIGHT_PAREN {
				break
			} else {
				return result, err
			}
		}

		if binaryExpr, ok := result.(BinaryExpression); ok {
			fmt.Fprintf(os.Stderr, "[DEBUG] binaryExpr: %s %s\n", binaryExpr.operator.lexeme, binaryExpr.right.String())
			result = BinaryExpression{left: left, operator: binaryExpr.operator, right: binaryExpr.right}
			// left = BinaryExpression{left: result, right: nil}
		}
		left = result
		fmt.Fprintf(os.Stderr, "[DEBUG] left: %s\n", left.String())
		errMain += err
		p.advance()
	}

	return left, errMain
}

func (p *Parser) parseToken() (ASTNode, string) {
	token := p.peek()
	//TODO:
	// p.advance()
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
		fmt.Fprint(os.Stderr, "[DEBUG] Left_Parent in stack:")
		fmt.Fprintln(os.Stderr, p.stack.items)
		expr, err := p.parseGroup()
		if expr != nil {
			return GroupExpression{expression: expr}, err
		} else {
			return nil, err
		}
	case RIGHT_PAREN:
		p.stack.Pop()
		return nil, "[RIGHT PAREN]Error: Unmatched parentheses.\n"
	case BANG, MINUS:
		expr, err := p.parseUnary()
		return UnaryExpression{expression: expr, operator: token}, err
	case STAR, SLASH:
		expr, err := p.parseBinary()

		return BinaryExpression{left: nil, operator: token, right: expr}, err
	default:
		fmt.Fprintf(os.Stderr, "[DEBUG] [Parser parseToken()] default. Token %s has not implemented\n", token.lexeme)
	}
	return nil, ""
}

func (p *Parser) parseUnary() (ASTNode, string) {
	//TODO:
	p.advance()
	expr, err := p.parseToken()
	return expr, err
}

func (p *Parser) parseBinary() (ASTNode, string) {
	//TODO:
	p.advance()
	right, errRight := p.parseToken()
	// fmt.Fprintf(os.Stderr, "[DEBUG] Right %s\n", right.String())
	return right, errRight
}

func (p *Parser) parseGroup() (ASTNode, string) {
	p.advance()

	var err string
	if p.isAtEnd() {
		err += "Error: Unmatched parentheses."
		return nil, err
	}
	token := p.peek()
	if token.typeToken == RIGHT_PAREN {
		//TODO:
		// p.advance()
		return nil, "Error: Empty expression."
	}

	expr, errParseTokens := p.parseTokens()
	err += errParseTokens
	if p.isAtEnd() {
		err += "Error: Unmatched parentheses."
		return nil, err
	}
	// token = p.peek()
	// if token.typeToken == RIGHT_PAREN {
	// 	top := p.stack.Pop()
	// 	if top != '(' {
	// 		err += "[p.stack.Pop() != '(']Error: Unmatched parentheses."
	// 		return nil, err
	// 	}

	// }
	// fmt.Fprint(os.Stderr, "[DEBUG] AFTER POP() Left_Parent in stack:")
	// fmt.Fprintln(os.Stderr, p.stack.items)
	//TODO:
	// p.advance()
	return expr, err
}

// func (p *Parser) parseBinary(left ASTNode) (ASTNode, string) {

// p.previous()
// var operator Token = p.peek()
// p.previous()
// left, errLeft := p.parseToken()
// p.advance()
// p.advance()
// fmt.Fprintf(os.Stderr, "[DEBUG] Left %s\n", left.String())
// fmt.Fprintf(os.Stderr, "[DEBUG] Operator %s\n", operator.String())
// return nil, ""

// right, errRight := p.parseToken()
// fmt.Fprintf(os.Stderr, "[DEBUG] Right %s\n", right.String())

// return right, errRight

// }

func (p *Parser) advance() {
	p.current = p.current + 1
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

// func (p *Parser) previous() Token {
// 	return p.tokens[p.current-1]
// }

func (p *Parser) previous() {
	p.current = p.current - 1
}

// func (p *Parser) previous() Token {
// 	return p.tokens[p.current-1]
// }

func Parse(tokens []Token) (ASTNode, string) {
	parser := initParser()
	parser.tokens = tokens
	expression, err := parser.parseTokens()

	if !parser.stack.IsEmpty() {
		return expression, "[p.stack.IsEmpty() ]Error: Unmatched parentheses."
	}
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
