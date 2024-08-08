package main

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

func (p *Parser) parseBinary() (ASTNode, string) {
	left, leftErr := p.parseUnary()

	for p.match(PLUS, MINUS, SLASH, STAR) {
		operator := p.previous()
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

func (p *Parser) advance() {
	p.current = p.current + 1
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
