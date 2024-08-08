package main

func (p *Parser) expression() ASTNode {
	return p.equality()
}

func (p *Parser) equality() ASTNode {
	expr := p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = BinaryExpression{expr, operator, right}
	}

	return expr
}

func (p *Parser) comparison() ASTNode {
	expr := p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = BinaryExpression{expr, operator, right}
	}

	return expr
}

func (p *Parser) term() ASTNode {
	expr := p.factor()

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = BinaryExpression{expr, operator, right}
	}

	return expr
}

func (p *Parser) factor() ASTNode {
	expr := p.unary()

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = BinaryExpression{expr, operator, right}
	}

	return expr
}

func (p *Parser) unary() ASTNode {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return UnaryExpression{operator, right}
	}

	return p.primary()
}

func (p *Parser) primary() ASTNode {
	token := p.peek()
	if p.match(FALSE) {
		return FalseExpression{token: token}
	}
	if p.match(TRUE) {
		return TrueExpression{token: token}
	}
	if p.match(NIL) {
		return NilExpression{token: token}
	}
	if p.match(NUMBER) {
		return NumberExpression{token: token}
	}
	if p.match(STRING) {
		return StringExpression{token: token}
	}

	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN)
		return GroupExpression{expression: expr}
	}
	return nil
}

func (p *Parser) match(tokensTypes ...TokenType) bool {
	for _, tokenType := range tokensTypes {
		if p.check(tokenType) {
			p.advanceAndReturnToken()
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

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(tokenType TokenType) {
	if p.check(tokenType) {
		p.advanceAndReturnToken()
		p.stackParentheses.Pop()
		return
	}
}

func (p *Parser) advanceAndReturnToken() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}
