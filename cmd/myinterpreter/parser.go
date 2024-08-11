package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Parser struct {
	tokens           []Token
	current          int
	stackParentheses Stack
	stackASTNodes    StackASTNode
}

type ASTNode interface {
	String() string
	Evaluate() Value
	// Type()
}
type NumberValue struct {
	value float64
}

func (n NumberValue) String() string {
	number := strconv.FormatFloat(n.value, 'f', -1, 64)

	if strings.Contains(number, ".") {
		number = strings.TrimRight(number, "0") // Удаляем нули справа
		number = strings.TrimRight(number, ".") // Удаляем точку, если все числа после нее были нулями
	}
	return number
}

type NumberExpression struct {
	token Token
}

func (n NumberExpression) String() string {
	return n.token.literal
}
func (n NumberExpression) Evaluate() Value {
	number, err := strconv.ParseFloat(n.token.literal, 64)

	if err != nil {
		fmt.Println("Ошибка преобразования:", err)
	}

	return NumberValue{number}
}

type StringValue struct {
	value string
}

func (n StringValue) String() string {
	return n.value
}

type StringExpression struct {
	token Token
}

func (n StringExpression) String() string {
	return n.token.literal
}
func (n StringExpression) Evaluate() Value {
	return StringValue{n.token.literal}
}

type Value interface {
	String() string
}
type BooleanValue struct {
	value bool
}

func (n BooleanValue) String() string {
	return strconv.FormatBool(n.value)
}

type TrueExpression struct {
	token Token
}

func (n TrueExpression) String() string {
	return n.token.lexeme
}
func (n TrueExpression) Evaluate() Value {
	return BooleanValue{value: true}
}

type FalseExpression struct {
	token Token
}

func (n FalseExpression) String() string {
	return n.token.lexeme
}
func (n FalseExpression) Evaluate() Value {
	return BooleanValue{value: false}
}
func (n NilValue) String() string {
	return "nil"
}

type NilValue struct {
}
type NilExpression struct {
	token Token
}

func (n NilExpression) String() string {
	return n.token.lexeme
}
func (n NilExpression) Evaluate() Value {
	return NilValue{}
}

type GroupExpression struct {
	expression ASTNode
}

func (n GroupExpression) String() string {
	return "(group " + n.expression.String() + ")"
}
func (n GroupExpression) Evaluate() Value {
	return n.expression.Evaluate()
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

func calc(l Value, r Value, operator Token) Value {

	left, okLeft := l.(NumberValue)
	right, okRight := r.(NumberValue)

	if okLeft && okRight {
		switch operator.typeToken {
		case PLUS:
			return NumberValue{left.value + right.value}
		case MINUS:
			return NumberValue{left.value - right.value}
		case SLASH:
			return NumberValue{left.value / right.value}
		case STAR:
			return NumberValue{left.value * right.value}
		}
	}

	leftString, okLeftString := l.(StringValue)
	rightString, okRightString := r.(StringValue)

	if okLeftString && okRightString {
		return StringValue{leftString.value + rightString.value}
	}

	return NilValue{}

}

func (n BinaryExpression) Evaluate() Value {
	left := n.left.Evaluate()
	right := n.right.Evaluate()
	left = calc(left, right, n.operator)
	return left // TODO: написать реализацию
}

func (n UnaryExpression) String() string {
	return "(" + n.operator.lexeme + " " + n.expression.String() + ")"
}
func (n UnaryExpression) Evaluate() Value {
	right := n.expression.Evaluate()

	switch n.operator.typeToken {
	case MINUS:
		switch v := right.(type) {
		case NumberValue:
			right = NumberValue{-v.value}
		}
	case BANG:
		switch v := right.(type) {
		case BooleanValue:
			right = BooleanValue{!v.value}
		case NilValue:
			right = BooleanValue{true}
		case NumberValue:
			right = BooleanValue{!(v.value != 0)}
		}

	default:

	}

	return right

	//switch n.operator.typeToken {
	//case MINUS:
	//	if number, ok := right.(NumberValue); ok {
	//		number.value = -number.value
	//	}
	//
	//case BANG:
	//}
}

func initParser() Parser {
	return Parser{
		current: 0,
	}
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens)
}

func (p *Parser) parseBinary() (ASTNode, string) {
	left, leftErr := p.parseUnary()

	for p.match(PLUS, MINUS, SLASH, STAR, LESS, GREATER, LESS_EQUAL, GREATER_EQUAL, BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()

		if p.isAtEnd() {
			return nil, "Error: Expect expression"
		}
		right, errorRight := p.parseUnary()
		left = BinaryExpression{left, operator, right}
		leftErr += errorRight
	}
	return left, leftErr
}

func (p *Parser) parseUnary() (ASTNode, string) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, errorRight := p.parseUnary()
		return UnaryExpression{operator, right}, errorRight
	}
	return p.parseLiteral()
}

func (p *Parser) parseLiteral() (ASTNode, string) {
	token := p.peek()

	if p.match(TRUE) {
		return TrueExpression{token: token}, ""
	}
	if p.match(FALSE) {
		return FalseExpression{token: token}, ""
	}
	if p.match(NIL) {
		return NilExpression{token: token}, ""
	}
	if p.match(NUMBER) {
		return NumberExpression{token: token}, ""
	}
	if p.match(STRING) {
		return StringExpression{token: token}, ""
	}
	if p.match(LEFT_PAREN) {
		p.stackParentheses.Push('(')
		token = p.peek()
		if token.typeToken == RIGHT_PAREN {
			//TODO: нужен ли тут advance, чтобы consume правую скобку
			return nil, "Error: Empty expression."
		}
		expr, err := p.parseTokens()
		p.consume(RIGHT_PAREN)
		return GroupExpression{expr}, err
	}
	return nil, "Error: default case in parseLiteral()"
}

func (p *Parser) parseTokens() (ASTNode, string) {
	return p.parseBinary()
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func Parse(tokens []Token) (ASTNode, string) {
	parser := initParser()
	parser.tokens = tokens
	exp, err := parser.parseTokens()

	//expression := parser.generateASTTree()

	if !parser.stackParentheses.IsEmpty() {
		return nil, "Error: Unmatched parentheses."
	}
	return exp, err
}
