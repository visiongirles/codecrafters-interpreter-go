package main

import (
	"fmt"
	"testing"
)

func TestSyntaxBinaryError(t *testing.T) {
	source := `(76 +)`
	tokens, errScan := Scan(source)
	expectedErr := " [line 1] Error: Unterminated string."

	if errScan != expectedErr {
		t.Errorf("Scanner Failed: %s\n", errScan)
	}

	result, errParse := Parse(tokens)
	expected := ""

	if errParse != expectedErr {
		t.Errorf("Parser Failed: %s\n", errParse)
	}
	fmt.Println(result)
	if result.String() != expected {
		t.Errorf("Result Failed! %s\n", result.String())
	}
}
func TestSyntaxError(t *testing.T) {
	source := `"foo`
	tokens, errScan := Scan(source)
	expectedErr := "[line 1] Error: Unterminated string."

	if errScan != expectedErr {
		t.Errorf("Scanner Failed: %s\n", errScan)
	}

	result, errParse := Parse(tokens)
	expected := ""

	if errParse != expectedErr {
		t.Errorf("Parser Failed: %s\n", errParse)
	}
	fmt.Println(result)
	if result.String() != expected {
		t.Errorf("Result Failed! %s\n", result.String())
	}
}
func TestParseBinaryMinus(t *testing.T) {
	source := "7 - 5"
	tokens, errScan := Scan(source)
	expectedErr := ""

	if errScan != expectedErr {
		t.Errorf("Scanner Failed: %s\n", errScan)
	}

	result, errParse := Parse(tokens)
	expected := "(- 7.0 5.0)"

	if errParse != expectedErr {
		t.Errorf("Parser Failed: %s\n", errParse)
	}
	fmt.Println(result)
	if result.String() != expected {
		t.Errorf("Result Failed! %s\n", result.String())
	}
}
func TestParseBinaryMinusNEW(t *testing.T) {
	source := "7 - 5"
	tokens, errScan := Scan(source)
	expectedErr := ""

	if errScan != expectedErr {
		t.Errorf("Scanner Failed: %s", errScan)
	}

	parser := initParser()
	parser.tokens = tokens

	result := parser.expression()
	errParse := ""
	expected := "(- 7.0 5.0)"
	fmt.Println(result)

	if errParse != expectedErr {
		t.Errorf("Parser Failed: %s", errParse)
	}

	if result.String() != expected {
		t.Errorf("Result Failed! %s\n", result.String())
	}
}
func TestParseBinaryAndUnary(t *testing.T) {
	source := "-5 + 9"
	tokens, errScan := Scan(source)
	expectedErr := ""

	if errScan != expectedErr {
		t.Errorf("Scanner Failed: %s", errScan)
	}

	result, errParse := Parse(tokens)
	expected := "(+ (- 5.0) 9.0)"
	fmt.Println(result)

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
	fmt.Println(result)

	if errParse != expectedErr {
		t.Errorf("Parser Failed: %s", errParse)
	}

	if result.String() != expected {
		t.Errorf("Result Failed! %s\n", result.String())
	}
}
func TestParseDoubleUnaryMinus(t *testing.T) {
	source := "-(-23 + 5)"
	tokens, errScan := Scan(source)
	expectedErr := ""

	if errScan != expectedErr {
		t.Errorf("Scanner Failed: %s", errScan)
	}

	result, errParse := Parse(tokens)
	expected := "(- (group (+ (- 23.0) 5.0)))"
	fmt.Println(result)

	if errParse != expectedErr {
		t.Errorf("Parser Failed: %s", errParse)
	}

	if result.String() != expected {
		t.Errorf("Result Failed! %s\n", result.String())
	}
}
func TestParseDoubleUnaryMinusNEW(t *testing.T) {
	source := "-(-23 + 5)"
	tokens, errScan := Scan(source)
	expectedErr := ""

	if errScan != expectedErr {
		t.Errorf("Scanner Failed: %s", errScan)
	}

	parser := initParser()
	parser.tokens = tokens

	result := parser.expression()

	fmt.Println(result)

	errParse := ""
	expected := "(- (group (+ (- 23.0) 5.0)))"

	if errParse != expectedErr {
		t.Errorf("Parser Failed: %s", errParse)
	}

	if result.String() != expected {
		t.Errorf("Result Failed! %s\n", result.String())
	}
}
func TestParseTwoGroupsAndUnary(t *testing.T) {
	source := `-("small") / ("talk")`
	tokens, errScan := Scan(source)
	expectedErr := ""

	if errScan != expectedErr {
		t.Errorf("Scanner Failed: %s", errScan)
	}

	result, errParse := Parse(tokens)
	// TODO:
	expected := "(/ (- (group small)) (group talk))"
	fmt.Println(result)

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
	fmt.Println(result)

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
	fmt.Println(result)

	expected := "(/ (* (- (group (+ (- 23.0) 54.0))) (group (* 87.0 34.0))) (group (+ 73.0 62.0)))"

	if errParse != expectedErr {
		t.Errorf("Parser Failed: %s", errParse)
	}

	if result.String() != expected {
		t.Errorf("Result Failed! Result.String(): %s\n", result.String())
	}
}
func TestParseBinaryMultipleOperandsAndOperatorsNEW(t *testing.T) {
	source := `-(-23 + 54) * (87 * 34) / (73 + 62)`
	tokens, errScan := Scan(source)
	expectedErr := ""

	if errScan != expectedErr {
		t.Errorf("Scanner Failed: %s", errScan)
	}

	parser := initParser()
	parser.tokens = tokens

	result := parser.expression()
	errParse := ""
	fmt.Println(result)
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
	fmt.Println(result)

	if errParse != expectedErr {
		t.Errorf("Parser Failed: %s", errParse)
	}

	if result.String() != expected {
		t.Errorf("Result Failed! %s\n", result.String())
	}
}
func TestParseBinaryMinusSeveralOperandsNEW(t *testing.T) {
	source := "90 - 94 - 33"
	tokens, errScan := Scan(source)
	expectedErr := ""

	if errScan != expectedErr {
		t.Errorf("Scanner Failed: %s", errScan)
	}

	parser := initParser()
	parser.tokens = tokens

	result := parser.expression()
	errParse := ""
	fmt.Println(result)
	expected := "(- (- 90.0 94.0) 33.0)"
	fmt.Println(result)

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
	fmt.Println(result)

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
	fmt.Println(result)

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
	fmt.Println(result)

	if errParse != expectedParseErr {
		t.Errorf("Parser Failed: %s", errParse)
	}

	if result != nil {
		t.Errorf("Result Failed! %s", result.String())
	}

}
