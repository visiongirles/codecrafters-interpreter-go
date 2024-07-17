package main

import (
	"fmt"
	"os"
)

func PrintErrors(errors string) {
	fmt.Fprint(os.Stderr, errors)
}

func PrintTokens(tokens []Token) {
	for _, value := range tokens {
		fmt.Print(value.String())
	}
}
