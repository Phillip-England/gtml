package main_test

import (
	"testing"

	"github.com/phillip-england/gtml/pkg/gtml"
)

func TestEdge_SelfClosingWithSpaces(t *testing.T) {
	state := createTestState(map[string]string{
		"Icon": `<svg></svg>`,
	})

	input := `<Icon   />`
	expected := `<svg></svg>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestEdge_ComponentWithNoProps(t *testing.T) {
	state := createTestState(map[string]string{
		"Divider": `<hr />`,
	})

	input := `<Divider />`
	expected := `<hr />`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestEdge_DeeplyNestedComponents(t *testing.T) {
	state := createTestState(map[string]string{
		"Level3": `<span props='val string'>{val}</span>`,
		"Level2": `<p props='data string'><Level3 val={data} /></p>`,
		"Level1": `<div props='input string'><Level2 data={input} /></div>`,
	})

	input := `<Level1 input='deep value' />`
	expected := `<div><p><span>deep value</span></p></div>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestEdge_MultipleComponentsInRoute(t *testing.T) {
	state := createTestState(map[string]string{
		"Para": `<p props='text string'>{text}</p>`,
	})

	input := `<div><Para text='First' /><Para text='Second' /></div>`
	expected := `<div><p>First</p><p>Second</p></div>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestEdge_ComponentWithTextAround(t *testing.T) {
	state := createTestState(map[string]string{
		"Bold": `<strong props='text string'>{text}</strong>`,
	})

	input := `<p>Hello <Bold text='World' /> and goodbye</p>`
	expected := `<p>Hello <strong>World</strong> and goodbye</p>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}
