package main

import (
	"fmt"
	"os"
	"strconv"
)

func buildResponse(token rune, tokenT string, content string, index int) (tokenType string, lexeme string) {
	tokenType = tokenT + "_EQUAL"
	lexeme = string(token) + string(content[index+1])
	return
}

func isDigit(token rune) bool {
	return token >= '0' && token <= '9'
}

func isAlpha(token rune) bool {
	return (token >= 'a' && token <= 'z') ||
		(token >= 'A' && token <= 'Z') ||
		token == '_'
}

func isAlphaNumeric(token rune) bool {
	return isDigit(token) || isAlpha(token)
}

func hasTrailingZeros(lexeme string) string {
	for {
		if (lexeme[len(lexeme)-1]) == '0' {
			if (lexeme[len(lexeme)-2]) == '.' {
				break
			} else {
				lexeme = lexeme[:len(lexeme)-1]
			}
		} else {
			break
		}
	}
	return lexeme
}

func main() {

	hasError := false

	singleTokens := map[rune]string{
		'(': "LEFT_PAREN",
		')': "RIGHT_PAREN",
		'{': "LEFT_BRACE",
		'}': "RIGHT_BRACE",
		',': "COMMA",
		'.': "DOT",
		'-': "MINUS",
		'+': "PLUS",
		';': "SEMICOLON",
		'*': "STAR",
	}

	operatorTokens := map[rune]string{
		'=': "EQUAL",
		'!': "BANG",
		'<': "LESS",
		'>': "GREATER",
	}

	// runeLiteralTokens:= []rune{' ', '\t', '\v'}
	runeLiteralTokens := map[rune]string{
		' ':  "WHITESPACE",
		'\t': "HORIZONTAL TAB",
		'\v': "VERTICAL TAB",
	}

	stringTokens := map[rune]string{
		'"': "STRING",
	}

	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	if len(fileContents) > 0 {
		content := string(fileContents)
		count := 1

		for index := 0; index < len(content); index++ {
			isNextTokenAlsoEqual := index+1 < len(content) && content[index+1] == '='
			isNextTokenAlsoSlash := index+1 < len(content) && content[index+1] == '/'
			var tokenType string
			var lexeme string
			token := rune(content[index])

			if token == '\n' {
				count++
				continue
			}
			if _, ok := runeLiteralTokens[token]; ok {
				continue
			}
			if value, ok := singleTokens[token]; ok {
				fmt.Printf("%s %c null\n", value, token)
			} else if _, ok := operatorTokens[token]; ok {
				switch token {
				case '=':
					{
						tokenType = operatorTokens[token]
						lexeme = string(token)
						if isNextTokenAlsoEqual {
							tokenType, lexeme = buildResponse(token, tokenType, content, index)
							index++
						}
					}
				case '!':
					{
						tokenType = operatorTokens[token]
						lexeme = string(token)
						if isNextTokenAlsoEqual {
							tokenType, lexeme = buildResponse(token, tokenType, content, index)
							index++
						}
					}
				case '<':
					{
						tokenType = operatorTokens[token]
						lexeme = string(token)
						if isNextTokenAlsoEqual {
							tokenType, lexeme = buildResponse(token, tokenType, content, index)
							index++
						}
					}
				case '>':
					{
						tokenType = operatorTokens[token]
						lexeme = string(token)
						if isNextTokenAlsoEqual {
							tokenType, lexeme = buildResponse(token, tokenType, content, index)
							index++
						}
					}
				}
				fmt.Printf("%s %s null\n", tokenType, lexeme)
			} else if token == '/' {
				if isNextTokenAlsoSlash {
					for {
						if content[index] == '\n' || index == len(content)-1 {
							count++
							break
						}
						index++
					}
				} else {
					tokenType = "SLASH"
					lexeme = string(token)
					fmt.Printf("%s %s null\n", tokenType, lexeme)

				}
			} else if tokenType, ok := stringTokens[token]; ok {
				index++
				var lexeme string
				hasStringError := false
				for {
					if content[index] == '"' {
						break
					}
					if index == len(content)-1 {
						fmt.Fprintf(os.Stderr, "[line %d] Error: Unterminated string.\n", count)
						hasStringError = true
						hasError = true
						break
					}
					lexeme += string(content[index])
					index++
				}
				if !hasStringError {
					fmt.Printf("%s \"%s\" %s\n", tokenType, lexeme, lexeme)
				}
			} else if isDigit(token) {
				tokenType = "NUMBER"
				lexeme += string(token)
				index++
				hasDot := false
				for {
					if index < len(content) {
						if isDigit(rune(content[index])) {
							lexeme += string(content[index])
							index++
							// fmt.Printf("lexeme: %s\n", lexeme)
							continue
						}
						if content[index] == '.' && !hasDot {
							hasDot = true
							lexeme += string(content[index])
							index++
							continue
						} else {
							index--
							break
						}
					} else {
						break
					}
				}
				if hasDot {
					if lexeme[len(lexeme)-1] == '.' {
						lexeme = lexeme[:len(lexeme)-1]
						floatValue, _ := strconv.ParseFloat(lexeme, 64)
						fmt.Printf("%s %s %.1f\n", tokenType, lexeme, floatValue)
						index = index - 2
					} else {
						stringLitral := hasTrailingZeros(lexeme)
						fmt.Printf("%s %s %s\n", tokenType, lexeme, stringLitral)
					}
				} else {
					floatValue, _ := strconv.ParseFloat(lexeme, 64)
					fmt.Printf("%s %s %.1f\n", tokenType, lexeme, floatValue)
				}
			} else if isAlpha(token) {
				tokenType = "IDENTIFIER"
				lexeme += string(token)
				index++
				for {
					if index < len(content) {
						if isAlphaNumeric(rune(content[index])) {
							lexeme += string(content[index])
							index++
						} else {
							index--
							break
						}
					} else {
						break
					}

				}
				fmt.Printf("%s %s null\n", tokenType, lexeme)

			} else {
				fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %c\n", count, token)
				hasError = true
			}
		}
	}
	fmt.Println("EOF  null")
	if hasError {
		os.Exit(65)
	}
}

// <token_type> <lexeme> <literal>
// token_type VAR, IDENTIFIER, STRING, EOF
// lexeme - string itself
// literal -   The literal value of the token: string/ number === string/number, other== = null
