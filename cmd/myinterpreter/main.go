package main

import (
	"fmt"
	"os"
)

func main() {

	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(InputError)
	}

	command := os.Args[1]
	filename := os.Args[2]
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(InputError)
	}

	switch command {
	case tokenize:
		tokens, err := handleTokenizeCommand(fileContent)
		if tokens != nil {
			PrintTokens(tokens)
		}

		if err != "" {
			PrintErrors(err)
			os.Exit(SyntaxError)
		}
	case parse:
		expression, _ := handleParseCommand(fileContent)
		fmt.Printf("%s\n", expression.String())

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(InputError)
	}

	os.Exit(OK)
}

func handleTokenizeCommand(fileContent []byte) ([]Token, string) {

	if len(fileContent) >= 0 {
		tokens, err := Scan(string(fileContent))
		return tokens, err
	}
	return nil, ""
}

func handleParseCommand(fileContent []byte) (ASTNode, string) {
	tokens, _ := handleTokenizeCommand(fileContent)
	expression, err := Parse(tokens)
	// result := expressionToString(expression)
	return expression, err
}

// func expressionToString(expression LiteralExpression) string {
// 	return expression.value.String()
// }

// <token_type> <lexeme> <literal>
// token_type VAR, IDENTIFIER, STRING, EOF
// lexeme - string itself
// literal -   The literal value of the token: string/ number === string/number, other== = null
