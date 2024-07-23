package main

import (
	"testing"
)

func TestParse(t *testing.T) {
	token1 := Token{typeToken: LEFT_PAREN}
	token2 := Token{typeToken: STRING, lexeme: "foo", literal: "foo"}
	token3 := Token{typeToken: RIGHT_PAREN}
	tokens := []Token{token1, token2, token3}
	result, err := Parse(tokens)
	expected := "(group foo)"
	expectedErr := ""

	if result.String() != expected && err != expectedErr {
		t.Errorf("ERROR!")
	}

	if 
}
