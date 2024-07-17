package main

import (
	"fmt"
	"strconv"
	"strings"
)

func Scan(source string) ([]Token, string) {
	scanner := initScanner()
	scanner.source = source
	scanner.scanTokens()
	return scanner.tokens, scanner.error
}

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
	error   string
}

func initScanner() Scanner {
	return Scanner{
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() {
	char := s.advance()
	switch char {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	case '\n':
		s.line++
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}
	case ' ':
	case '\r':
	case '\t':
	case '"':
		s.stringLiteral()
	default:
		if isDigit(char) {
			s.number()
		} else if isAlpha(char) {
			s.identifier()
		} else {
			s.error += fmt.Sprintf("[line %d] Error: Unexpected character: %c\n", s.line, char)
		}
	}
}

func (s *Scanner) number() {
	for isDigit((s.peek())) {
		s.advance()
	}

	// Look for a fractional part.
	if s.peek() == '.' && isDigit(s.peekNext()) {
		// Consume the "."
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}
	s.addTokenFloat(NUMBER, s.source[s.start:s.current])
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return '\n'
	}
	return rune(s.source[s.current+1])
}

func (s *Scanner) stringLiteral() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.error += fmt.Sprintf("[line %d] Error: Unterminated string.\n", s.line)
		return
	}
	// The closing ".
	s.advance()
	// Trim the surrounding quotes.
	value := s.source[s.start+1 : s.current-1]
	s.addTokenString(STRING, value)
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\n'
	}
	return rune(s.source[s.current])
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if rune(s.source[s.current]) != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) advance() rune {
	index := s.current
	s.current++
	return rune(s.source[index])
}

func (s *Scanner) scanTokens() {
	for !s.isAtEnd() {
		// We are at the beginning of the next lexeme.
		s.start = s.current
		s.scanToken()
	}

	newToken := Token{typeToken: EOF, lexeme: "", literal: "null", line: s.line}
	s.tokens = append(s.tokens, newToken)
}

func (s *Scanner) addToken(tokenType TokenType) {
	s.addTokenString(tokenType, "null")
}

func (s *Scanner) addTokenString(tokenType TokenType, literal string) {
	text := s.source[s.start:s.current]
	newToken := Token{typeToken: tokenType, lexeme: text, literal: literal, line: s.line}
	s.tokens = append(s.tokens, newToken)
}

func (s *Scanner) addTokenFloat(tokenType TokenType, lexeme string) {
	floatValue, err := strconv.ParseFloat(lexeme, 64)
	if err != nil {
		s.error += "[Error]: Float parse error"
	} else {
		str := strconv.FormatFloat(floatValue, 'f', -1, 64)
		literal := str
		if !strings.Contains(str, ".") {
			literal += ".0"
		}
		newToken := Token{typeToken: tokenType, lexeme: lexeme, literal: literal, line: s.line}
		s.tokens = append(s.tokens, newToken)
	}
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}
	str := s.source[s.start:s.current]
	if tokenType, ok := ReservedWordTokens[str]; ok {
		s.addToken(tokenType)
	} else {
		s.addToken(IDENTIFIER)
	}
}

// available in unicode package
func isDigit(token rune) bool {
	return token >= '0' && token <= '9'
}

// available in unicode package
func isAlpha(token rune) bool {
	return (token >= 'a' && token <= 'z') ||
		(token >= 'A' && token <= 'Z') ||
		token == '_'
}

// available in unicode package
func isAlphaNumeric(token rune) bool {
	return isDigit(token) || isAlpha(token)
}
