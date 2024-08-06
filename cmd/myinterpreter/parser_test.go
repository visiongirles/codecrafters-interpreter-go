package main

import (
	"testing"
)

// func TestParse(t *testing.T) {
// 	source := "(76 * -30 / (84 * 39))"
// 	tokens, errScan := Scan(source)
// 	expectedErr := ""

// 	if errScan != expectedErr {
// 		t.Errorf("Scanner Failed: %s", errScan)
// 	}

// 	result, errParse := Parse(tokens)
// 	expected := "(group (/ (* 76.0 (- 30.0)) (group (* 84.0 39.0))))"

// 	if errParse != expectedErr {
// 		t.Errorf("Parser Failed: %s", errParse)
// 	}

// 	if result.String() != expected {
// 		t.Errorf("Result Failed!")
// 	}
// }

func TestParseBinaryMinus(t *testing.T) {
	source := "7 - 5"
	tokens, errScan := Scan(source)
	expectedErr := ""

	if errScan != expectedErr {
		t.Errorf("Scanner Failed: %s", errScan)
	}

	result, errParse := Parse(tokens)
	expected := "(- 7.0 5.0)"

	if errParse != expectedErr {
		t.Errorf("Parser Failed: %s", errParse)
	}

	if result.String() != expected {
		t.Errorf("Result Failed! %s\n", result.String())
	}
}
func TestParseUnaryMinus(t *testing.T) {
	source := "- 7"
	tokens, errScan := Scan(source)
	expectedErr := ""

	if errScan != expectedErr {
		t.Errorf("Scanner Failed: %s", errScan)
	}

	result, errParse := Parse(tokens)
	expected := "(- 7.0)"

	if errParse != expectedErr {
		t.Errorf("Parser Failed: %s", errParse)
	}

	if result.String() != expected {
		t.Errorf("Result Failed! %s\n", result.String())
	}
}

func TestParseTwoGroupsndUnary(t *testing.T) {
	source := `-("small") / ("talk")`
	tokens, errScan := Scan(source)
	expectedErr := ""

	if errScan != expectedErr {
		t.Errorf("Scanner Failed: %s", errScan)
	}

	result, errParse := Parse(tokens)
	// TODO:
	expected := "(/ (group small) (group talk))"

	if errParse != expectedErr {
		t.Errorf("Parser Failed: %s", errParse)
	}

	if result.String() != expected {
		t.Errorf("Result Failed! Result.String(): %s\n", result.String())
	}
}
func TestParseTwoGroups(t *testing.T) {
	source := `("small") / ("talk")`
	tokens, errScan := Scan(source)
	expectedErr := ""

	if errScan != expectedErr {
		t.Errorf("Scanner Failed: %s", errScan)
	}

	result, errParse := Parse(tokens)
	// TODO:
	expected := "(/ (group small) (group talk))"

	if errParse != expectedErr {
		t.Errorf("Parser Failed: %s", errParse)
	}

	if result.String() != expected {
		t.Errorf("Result Failed! Result.String(): %s\n", result.String())
	}
}
func TestParseBinaryMultipleOperandsAndOperators(t *testing.T) {
	source := `-(-23 + 54) * (87 * 34) / (73 + 62)`
	tokens, errScan := Scan(source)
	expectedErr := ""

	if errScan != expectedErr {
		t.Errorf("Scanner Failed: %s", errScan)
	}

	result, errParse := Parse(tokens)
	// TODO:
	expected := "(/ (* (- (group (+ (- 23.0) 54.0))) (group (* 87.0 34.0))) (group (+ 73.0 62.0)))"

	if errParse != expectedErr {
		t.Errorf("Parser Failed: %s", errParse)
	}

	if result.String() != expected {
		t.Errorf("Result Failed! Result.String(): %s\n", result.String())
	}
}
func TestParseBinaryMinusSeveralOperands(t *testing.T) {
	source := "90 - 94 - 33"
	tokens, errScan := Scan(source)
	expectedErr := ""

	if errScan != expectedErr {
		t.Errorf("Scanner Failed: %s", errScan)
	}

	result, errParse := Parse(tokens)
	expected := "(- (- 90.0 94.0) 33.0)"

	if errParse != expectedErr {
		t.Errorf("Parser Failed: %s", errParse)
	}

	if result.String() != expected {
		t.Errorf("Result Failed! %s\n", result.String())
	}
}
func TestParseBinaryMinusTwoOperands(t *testing.T) {
	source := "(90 - 94)"
	tokens, errScan := Scan(source)
	expectedErr := ""

	if errScan != expectedErr {
		t.Errorf("Scanner Failed: %s", errScan)
	}

	result, errParse := Parse(tokens)
	expected := "(group (- 90.0 94.0))"

	if errParse != expectedErr {
		t.Errorf("Parser Failed: %s", errParse)
	}

	if result.String() != expected {
		t.Errorf("Result Failed! %s\n", result.String())
	}
}
func TestParseBinaryLiteral(t *testing.T) {
	source := `"hello" + "world"`
	tokens, errScan := Scan(source)
	expectedErr := ""

	if errScan != expectedErr {
		t.Errorf("Scanner Failed: %s", errScan)
	}

	result, errParse := Parse(tokens)
	expected := "(+ hello world)"

	if errParse != expectedErr {
		t.Errorf("Parser Failed: %s", errParse)
	}

	if result.String() != expected {
		t.Errorf("Result Failed! %s\n", result.String())
	}
}

func TestParseUnmatched(t *testing.T) {

	source := `("foo"`
	tokens, errScan := Scan(source)
	expectedScanErr := ""
	expectedParseErr := "Error: Unmatched parentheses."

	if errScan != expectedScanErr {
		t.Errorf("Scanner Failed: %s", errScan)
	}

	result, errParse := Parse(tokens)

	if errParse != expectedParseErr {
		t.Errorf("Parser Failed: %s", errParse)
	}

	if result != nil {
		t.Errorf("Result Failed! %s", result.String())
	}

}

// func TestParseUnaryMinus(t *testing.T) {

// 	source := `-1`
// 	tokens, errScan := Scan(source)
// 	expectedScanErr := ""
// 	expectedParseErr := ""
// 	expected := "- 1"

// 	if errScan != expectedScanErr {
// 		t.Errorf("Scanner Failed: %s", errScan)
// 	}

// 	result, errParse := Parse(tokens)

// 	if errParse != expectedParseErr {
// 		t.Errorf("Parser Failed: %s", errParse)
// 	}

// 	if result.String() != expected {
// 		t.Errorf("Result Failed! %s", result.String())
// 	}

// }
