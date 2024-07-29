package main

import (
	"testing"
)

func TestParse(t *testing.T) {
	// token1 := Token{typeToken: LEFT_PAREN}
	// token2 := Token{typeToken: STRING, lexeme: "foo", literal: "foo"}
	// token3 := Token{typeToken: RIGHT_PAREN}
	// tokens := []Token{token1, token2, token3}

	source := "(76 * -30 / (84 * 39))"
	// source := "(\"foo\")"
	tokens, errScan := Scan(source)
	expectedErr := ""

	if errScan != expectedErr {
		t.Errorf("Scanner Failed: %s", errScan)
	}

	result, errParse := Parse(tokens)
	expected := "(group (/ (* 76.0 (- 30.0)) (group (* 84.0 39.0))))"

	if errParse != expectedErr {
		t.Errorf("Parser Failed: %s", errParse)
	}

	if result.String() != expected {
		t.Errorf("Result Failed!")
	}

}

func TestParseUnmatched(t *testing.T) {

	source := "(\"foo\")"
	tokens, errScan := Scan(source)
	expectedScanErr := "Error: Unmatched parentheses."
	expectedParseErr := "Error: Unmatched parentheses."

	if errScan != expectedScanErr {
		t.Errorf("Scanner Failed: %s", errScan)
	}

	result, errParse := Parse(tokens)

	if errParse != expectedParseErr {
		t.Errorf("Parser Failed: %s", errParse)
	}

	if result != nil {
		t.Errorf("Result Failed!")
	}

}

// func TestScan(t *testing.T) {
// 	text := "(76 * -30 / (84 * 39))"

// 	result, err := Scan(text)
// 	var expected []Token
// 	expectedErr := ""

// 	if !reflect.DeepEqual(result, expected) || err != expectedErr {
// 		t.Errorf("ERROR!")
// 	}
// }
