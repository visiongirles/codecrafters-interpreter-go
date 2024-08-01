package main

import (
	"fmt"
	"os"
)

type Parser struct {
	tokens           []Token
	current          int
	stackParentheses Stack
	stackASTNodes    StackASTNode
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

main:
	for !p.isAtEnd() {
		token := p.peek()
		fmt.Fprintf(os.Stderr, "[DEBUG] [parseTokens()]. Token %s\n", token.lexeme)
		switch token.typeToken {
		case TRUE:
			left = TrueExpression{token: token}
		case FALSE:
			left = FalseExpression{token: token}
		case NIL:
			left = NilExpression{token: token}
		case NUMBER:
			left = NumberExpression{token: token}
		case STRING:
			left = StringExpression{token: token}
		case BANG, MINUS:
			expr, err := p.parseUnary()
			if err != "" {
				errMain += err
				return nil, errMain
			}
			left = expr
			// errMain += err
		// 	return UnaryExpression{expression: expr, operator: token}, err
		case LEFT_PAREN:
			p.stackParentheses.Push('(')
			expr, err := p.parseGroup()

			if err != "" {
				errMain += err
				return nil, errMain
			}

			if expr != nil {
				left = GroupExpression{expression: expr}
			} else {
				// TODO: p.stuckASTNode ????
				errMain += err
				return nil, errMain
			}

		case RIGHT_PAREN:
			p.stackParentheses.Pop()
			break main
		default:
			fmt.Fprintf(os.Stderr, "[DEBUG] [Parser parseToken()] default. Token %s has not implemented\n", token.lexeme)
		}

		p.advance()
		if p.isAtEnd() {
			p.stackASTNodes.Push(left)
			return left, errMain
		}

		operator := p.peek()
		fmt.Fprintf(os.Stderr, "[DEBUG] [parseTokens()]. Token %s\n", operator.lexeme)
		// p.stackASTNodes.Push(operator)

		// TODO: проверка на оператор
		switch operator.typeToken {
		case STAR, SLASH, PLUS, MINUS:
			expr, err := p.parseBinary()
			//TODO:

			if err != "" {
				errMain += err
				return nil, errMain
			}
			//TODO: second push?!
			// p.stackASTNodes.Push(expr)
			// left = BinaryExpression{left: left, operator: operator, right: expr}
			// left = BinaryExpression{operator: operator, right: expr}
			p.stackASTNodes.Push(BinaryExpression{operator: operator, right: expr})
		}

	}
	p.stackASTNodes.Push(left)
	return left, errMain
}

func (p *Parser) generateASTTree() ASTNode {
	left := p.stackASTNodes.Pop()
	for !p.stackASTNodes.IsEmpty() {
		node := p.stackASTNodes.Pop()

		if binaryExpr, ok := node.(BinaryExpression); ok {
			left = BinaryExpression{left: left, operator: binaryExpr.operator, right: binaryExpr.right}
			p.stackASTNodes.Pop()
		}

	}
	return left
}

func (p *Parser) parseUnary() (ASTNode, string) {
	p.advance()
	token := p.peek()
	expr, err := p.parseTokens()
	return UnaryExpression{expression: expr, operator: token}, err
}

func (p *Parser) parseBinary() (ASTNode, string) {
	//TODO:
	p.advance()
	right, errRight := p.parseTokens()
	// fmt.Fprintf(os.Stderr, "[DEBUG] Right %s\n", right.String())
	return right, errRight
}

func (p *Parser) parseGroup() (ASTNode, string) {
	p.advance()

	var err string

	// Сразу конец выражения, есть только левая скобка
	if p.isAtEnd() {
		err += "[parseGroup(), condition: p.isAtEnd()]Error: Unmatched parentheses."
		return nil, err
	}
	token := p.peek()

	// сразу следом за левой скобкой идет правая скобка - пустое выражение
	if token.typeToken == RIGHT_PAREN {
		//TODO: нужен ли тут advance, чтобы consume правую скобку
		return nil, "Error: Empty expression."
	}

	expr, errParseTokens := p.parseTokens()
	err += errParseTokens
	// if p.isAtEnd() {
	// 	err += "Error: Unmatched parentheses."
	// 	return nil, err
	// }
	return expr, err
}

func (p *Parser) advance() {
	p.current = p.current + 1
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func Parse(tokens []Token) (ASTNode, string) {
	parser := initParser()
	parser.tokens = tokens
	_, err := parser.parseTokens()

	expression := parser.generateASTTree()

	if !parser.stackParentheses.IsEmpty() {
		//TODO: nil или что-то распарсенное?
		return nil, "Error: Unmatched parentheses."
	}
	return expression, err
}
