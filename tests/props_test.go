package main_test

import (
	"strings"
	"testing"

	"github.com/phillip-england/gtml/pkg/gtml"
)

func TestProps_BasicStringProp(t *testing.T) {
	state := createTestState(map[string]string{
		"ThatButton": `<button props='text string'>{text}</button>`,
	})

	input := `<ThatButton text='some title' />`
	expected := `<button>some title</button>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestProps_MultipleProps(t *testing.T) {
	state := createTestState(map[string]string{
		"Card": `<div props='heading string, subheading string'><h1>{heading}</h1><p>{subheading}</p></div>`,
	})

	input := `<Card heading='Hello' subheading='World' />`
	expected := `<div><h1>Hello</h1><p>World</p></div>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestProps_IntType(t *testing.T) {
	state := createTestState(map[string]string{
		"Counter": `<span props='count int'>Count: {count}</span>`,
	})

	input := `<Counter count={42} />`
	expected := `<span>Count: 42</span>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestProps_IntWithArithmetic(t *testing.T) {
	state := createTestState(map[string]string{
		"Counter": `<span props='count int'>Count: {count}</span>`,
	})

	input := `<Counter count={40 + 2} />`
	expected := `<span>Count: 42</span>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestProps_DoubleQuotes(t *testing.T) {
	state := createTestState(map[string]string{
		"Message": `<p props='text string'>{text}</p>`,
	})

	input := `<Message text="Hello World" />`
	expected := `<p>Hello World</p>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestProps_NonSelfClosingComponent(t *testing.T) {
	state := createTestState(map[string]string{
		"Wrapper": `<div props='title string' class="wrap">{title}</div>`,
	})

	input := `<Wrapper title="Test Title"></Wrapper>`
	expected := `<div class="wrap">Test Title</div>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestProps_NestedComponents(t *testing.T) {
	state := createTestState(map[string]string{
		"Inner": `<span props='value string'>{value}</span>`,
		"Outer": `<div props='text string'><Inner value='nested' /></div>`,
	})

	input := `<Outer text='unused' />`
	expected := `<div><span>nested</span></div>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestProps_PropsRemovedFromOutput(t *testing.T) {
	state := createTestState(map[string]string{
		"MyComp": `<div props='title string'><h1>{title}</h1></div>`,
	})

	input := `<MyComp title='Hello' />`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(result, "props=") {
		t.Errorf("props attribute should be removed from output, got: %s", result)
	}
}
