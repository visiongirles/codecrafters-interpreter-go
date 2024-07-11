package main

import (
	"fmt"
	"os"
)

// type OperatorType int

// const (
// 	StringType OperatorType = iota
// 	RuneType
// )

// type OperatorValue struct {
// 	Type        OperatorType
// 	StringValue string
// 	RuneValue   rune
// }

func buildResponse(token rune, tokenT string, content string, index int) (tokenType string, lexeme string) {
	tokenType = tokenT + "_EQUAL"
	lexeme = string(token) + string(content[index+1])
	return
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
		'/': "SLASH",
	}

	operators := map[rune]string{
		'=': "EQUAL",
		// "==": "EQUAL_EQUAL",
		'!': "BANG",
		// "!=": "BANG_EQUAL",
		'<': "LESS",
		// "<=": "LESS_EQUAL",
		'>': "GREATER",
		// ">=": "GREATER_EQUAL",
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
			token := rune(content[index])
			if token == '\n' {
				count++
			}
			if value, ok := singleTokens[token]; ok {
				fmt.Printf("%s %c null\n", value, token)
			} else if _, ok := operators[token]; ok {
				var tokenType string
				var lexeme string
				isNextTokenAlsoEqual := index+1 < len(content) && content[index+1] == '='

				switch token {
				case '=':
					{
						tokenType = "EQUAL"
						lexeme = string(token)
						if isNextTokenAlsoEqual {
							tokenType, lexeme = buildResponse(token, tokenType, content, index)
							index++
						}
					}
				case '!':
					{
						tokenType = "BANG"
						lexeme = string(token)
						if isNextTokenAlsoEqual {
							tokenType, lexeme = buildResponse(token, tokenType, content, index)
							index++
						}
					}
				case '<':
					{
						tokenType = "LESS"
						lexeme = string(token)
						if isNextTokenAlsoEqual {
							tokenType, lexeme = buildResponse(token, tokenType, content, index)
							index++
						}
					}
				case '>':
					{
						tokenType = "GREATER"
						lexeme = string(token)
						if isNextTokenAlsoEqual {
							tokenType, lexeme = buildResponse(token, tokenType, content, index)
							index++
						}
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
