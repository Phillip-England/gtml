package main_test

import (
	"testing"

	"github.com/phillip-england/gtml/pkg/gtml"
)

func TestTernary_Basic(t *testing.T) {
	state := createTestState(map[string]string{
		"ShowIfActive": `<div props='active boolean'>{ active == true ? (<p>Active</p>) : (<p>Inactive</p>) }</div>`,
	})

	input := `<ShowIfActive active={true} />`
	expected := `<div><p>Active</p></div>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}

	input2 := `<ShowIfActive active={false} />`
	expected2 := `<div><p>Inactive</p></div>`

	result2, err := gtml.CompileHTML(input2, state, map[string]gtml.Value{}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result2) != normalizeHTML(expected2) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected2, result2)
	}
}

func TestTernary_NumericComparison(t *testing.T) {
	state := createTestState(map[string]string{
		"AgeCheck": `<div props='age int'>{ age >= 18 ? (<p>Adult</p>) : (<p>Minor</p>) }</div>`,
	})

	input := `<AgeCheck age={21} />`
	expected := `<div><p>Adult</p></div>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}

	input2 := `<AgeCheck age={15} />`
	expected2 := `<div><p>Minor</p></div>`

	result2, err := gtml.CompileHTML(input2, state, map[string]gtml.Value{}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result2) != normalizeHTML(expected2) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected2, result2)
	}
}

func TestTernary_StringComparison(t *testing.T) {
	state := createTestState(map[string]string{
		"ColorCheck": `<div props='color string'>{ color == "blue" ? (<p>blue</p>) : (<p>not blue</p>) }</div>`,
	})

	input := `<ColorCheck color='blue' />`
	expected := `<div><p>blue</p></div>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}

	input2 := `<ColorCheck color='red' />`
	expected2 := `<div><p>not blue</p></div>`

	result2, err := gtml.CompileHTML(input2, state, map[string]gtml.Value{}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result2) != normalizeHTML(expected2) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected2, result2)
	}
}

func TestTernary_LogicalAnd(t *testing.T) {
	state := createTestState(map[string]string{
		"AccessCheck": `<div props='role string, active boolean'>{ role == "admin" && active == true ? (<p>Full Access</p>) : (<p>Limited Access</p>) }</div>`,
	})

	input := `<AccessCheck role='admin' active={true} />`
	expected := `<div><p>Full Access</p></div>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}

	input2 := `<AccessCheck role='admin' active={false} />`
	expected2 := `<div><p>Limited Access</p></div>`

	result2, err := gtml.CompileHTML(input2, state, map[string]gtml.Value{}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result2) != normalizeHTML(expected2) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected2, result2)
	}
}

func TestTernary_LogicalOr(t *testing.T) {
	state := createTestState(map[string]string{
		"PriorityUser": `<div props='role string'>{ role == "admin" || role == "moderator" ? (<p>Priority User</p>) : (<p>Standard User</p>) }</div>`,
	})

	input := `<PriorityUser role='admin' />`
	expected := `<div><p>Priority User</p></div>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}

	input2 := `<PriorityUser role='guest' />`
	expected2 := `<div><p>Standard User</p></div>`

	result2, err := gtml.CompileHTML(input2, state, map[string]gtml.Value{}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result2) != normalizeHTML(expected2) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected2, result2)
	}
}

func TestTernary_NotEqual(t *testing.T) {
	state := createTestState(map[string]string{
		"NotBanned": `<div props='status string'>{ status != "banned" ? (<p>Welcome back!</p>) : (<p>Account suspended</p>) }</div>`,
	})

	input := `<NotBanned status='active' />`
	expected := `<div><p>Welcome back!</p></div>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}

	input2 := `<NotBanned status='banned' />`
	expected2 := `<div><p>Account suspended</p></div>`

	result2, err := gtml.CompileHTML(input2, state, map[string]gtml.Value{}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result2) != normalizeHTML(expected2) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected2, result2)
	}
}

func TestTernary_WithPropExpressions(t *testing.T) {
	state := createTestState(map[string]string{
		"ScoreCard": `<div props='score int'><h2>Your score: {score}</h2>{ score >= 50 ? (<p>You passed!</p>) : (<p>You failed.</p>) }</div>`,
	})

	input := `<ScoreCard score={75} />`
	expected := `<div><h2>Your score: 75</h2><p>You passed!</p></div>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}
