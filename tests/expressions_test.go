package main_test

import (
	"testing"

	"github.com/phillip-england/gtml/pkg/gtml"
)

func TestExpression_Addition(t *testing.T) {
	state := createTestState(map[string]string{
		"Math": `<span props='a int, b int'>{a + b}</span>`,
	})

	input := `<Math a={5} b={3} />`
	expected := `<span>8</span>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestExpression_Subtraction(t *testing.T) {
	state := createTestState(map[string]string{
		"Math": `<span props='a int, b int'>{a - b}</span>`,
	})

	input := `<Math a={10} b={4} />`
	expected := `<span>6</span>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestExpression_Multiplication(t *testing.T) {
	state := createTestState(map[string]string{
		"Math": `<span props='a int, b int'>{a * b}</span>`,
	})

	input := `<Math a={6} b={7} />`
	expected := `<span>42</span>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestExpression_Division(t *testing.T) {
	state := createTestState(map[string]string{
		"Math": `<span props='a int, b int'>{a / b}</span>`,
	})

	input := `<Math a={20} b={4} />`
	expected := `<span>5</span>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestExpression_Modulo(t *testing.T) {
	state := createTestState(map[string]string{
		"Math": `<span props='a int, b int'>{a % b}</span>`,
	})

	input := `<Math a={17} b={5} />`
	expected := `<span>2</span>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestExpression_StringConcatenation(t *testing.T) {
	state := createTestState(map[string]string{
		"FullName": `<p props='first string, last string'>{first + " " + last}</p>`,
	})

	input := `<FullName first='John' last='Doe' />`
	expected := `<p>John Doe</p>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestExpression_NegativeInt(t *testing.T) {
	state := createTestState(map[string]string{
		"Counter": `<span props='count int'>Count: {count}</span>`,
	})

	input := `<Counter count={-5} />`
	expected := `<span>Count: -5</span>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}
