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

func (p *Parser) parseLiteral() (ASTNode, string) {
	token := p.peek()
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
	default:
		return nil, "Error: default case in parseLiteral()"
	}
}

func (p *Parser) parseBinaryNew() (ASTNode, string) {
	left, leftErr := p.parseUnaryNew()

	for p.match(PLUS, MINUS, SLASH, STAR) {
		operator := p.previous()
		right, errorRight := p.parseUnaryNew()
		left = BinaryExpression{left, operator, right}
		leftErr += errorRight
	}
	return left, leftErr
}

func (p *Parser) parseUnaryNew() (ASTNode, string) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, errorRight := p.parseUnaryNew()
		return UnaryExpression{operator, right}, errorRight
	}
	return p.parseLiteralNew()
}

func (p *Parser) parseLiteralNew() (ASTNode, string) {
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
		expr, err := p.parseTokens()
		p.consume(RIGHT_PAREN)
		return GroupExpression{expr}, err
	}
	return nil, "Error: default case in parseLiteral()"
}

func (p *Parser) parseTokens() (ASTNode, string) {
	return p.parseBinaryNew()
}

//func (p *Parser) parseTokens() (ASTNode, string) {
//	var left ASTNode
//	errMain := ""
//
//main:
//	for !p.isAtEnd() {
//		token := p.peek()
//		fmt.Fprintf(os.Stderr, "[DEBUG] [parseTokens()]. Token: %s\n", token.lexeme)
//		switch token.typeToken {
//		case BANG, MINUS:
//			if left != nil {
//				right, err := p.parseBinary()
//				left = BinaryExpression{left, token, right}
//				errMain += err
//			} else {
//				expr, err := p.parseUnary()
//				left = expr
//				errMain += err
//
//			}
//		case TRUE, FALSE, NIL, NUMBER, STRING:
//			expr, err := p.parseLiteral()
//			left = expr
//			errMain += err
//		case LEFT_PAREN:
//			p.stackParentheses.Push('(')
//			expr, err := p.parseGroup()
//
//			if err != "" {
//				errMain += err
//				return nil, errMain
//			}
//
//			if expr != nil {
//				left = expr
//			} else {
//				errMain += err
//				return nil, errMain
//			}
//		case RIGHT_PAREN:
//			p.stackParentheses.Pop()
//			break main
//		default:
//			fmt.Fprintf(os.Stderr, "[DEBUG] [Parser parseToken()] default. Token %s has not implemented\n", token.lexeme)
//		}
//
//		p.advance()
//		if p.isAtEnd() {
//			p.stackASTNodes.Push(left)
//			return left, errMain
//		}
//
//		operator := p.peek()
//		fmt.Fprintf(os.Stderr, "[DEBUG] [parseTokens()]. Operator %s\n", operator.lexeme)
//
//		switch operator.typeToken {
//		case STAR, SLASH, PLUS, MINUS:
//			expr, err := p.parseBinary()
//			if err != "" {
//				errMain += err
//				return nil, errMain
//			}
//			left = BinaryExpression{left: left, operator: operator, right: expr}
//			p.advance() //TODO:
//		}
//	}
//	p.stackASTNodes.Push(left)
//	return left, errMain
//}

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
	var right ASTNode
	operator := p.peek()
	mainErr := ""
	p.advance()
	token := p.peek()
	if token.typeToken == LEFT_PAREN {
		right, mainErr = p.parseGroup()
	} else {
		right, mainErr = p.parseLiteral()
	}
	return UnaryExpression{expression: right, operator: operator}, mainErr
}

func (p *Parser) parseBinary() (ASTNode, string) {
	var right ASTNode
	err := ""
	p.advance()
	token := p.peek()
	switch token.typeToken {
	case LEFT_PAREN:

		right, err = p.parseGroup()
	case BANG, MINUS:
		right, err = p.parseUnary()
	case NUMBER, STRING, NIL, TRUE, FALSE:
		right, err = p.parseLiteral()
	}
	//right, errRight := p.parseTokens()
	return right, err
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

	return GroupExpression{expression: expr}, err
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
