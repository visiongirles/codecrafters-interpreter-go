package main

var ReservedWordTokens = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

const (
	OK          = 0
	InputError  = 1
	SyntaxError = 65
)

const (
	tokenize = "tokenize"
	parse    = "parse"
	evaluate = "evaluate"
)
