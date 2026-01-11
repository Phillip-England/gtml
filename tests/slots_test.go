package main_test

import (
	"testing"

	"github.com/phillip-england/gtml/pkg/gtml"
)

func TestSlot_BasicSlot(t *testing.T) {
	state := createTestState(map[string]string{
		"PageLayout": `<html><body><header>Site Header</header>{{ slot: content }}<footer>Site Footer</footer></body></html>`,
	})

	input := `<PageLayout><slot name='content' tag='main'><p>Hello World</p></slot></PageLayout>`
	expected := `<html><body><header>Site Header</header><main><p>Hello World</p></main><footer>Site Footer</footer></body></html>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", normalizeHTML(expected), normalizeHTML(result))
	}
}

func TestSlot_MultipleSlots(t *testing.T) {
	state := createTestState(map[string]string{
		"TwoColumnLayout": `<div class='container'>{{ slot: sidebar }}{{ slot: main }}</div>`,
	})

	input := `<TwoColumnLayout><slot name='sidebar' tag='aside'><nav>Navigation</nav></slot><slot name='main' tag='section'><p>Main content here</p></slot></TwoColumnLayout>`
	expected := `<div class='container'><aside><nav>Navigation</nav></aside><section><p>Main content here</p></section></div>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", normalizeHTML(expected), normalizeHTML(result))
	}
}

func TestSlot_WithClass(t *testing.T) {
	state := createTestState(map[string]string{
		"Layout": `<div>{{ slot: content }}</div>`,
	})

	input := `<Layout><slot name='content' tag='div' class='my-class'><p>Content</p></slot></Layout>`
	expected := `<div><div class='my-class'><p>Content</p></div></div>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", normalizeHTML(expected), normalizeHTML(result))
	}
}

func TestSlot_EmptySlot(t *testing.T) {
	state := createTestState(map[string]string{
		"Layout": `<div>{{ slot: content }}</div>`,
	})

	input := `<Layout></Layout>`
	expected := `<div></div>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestSlot_NestedComponentInSlot(t *testing.T) {
	state := createTestState(map[string]string{
		"Button": `<button props='label string'>{label}</button>`,
		"Card":   `<div class="card">{{ slot: actions }}</div>`,
	})

	input := `<Card><slot name='actions' tag='div'><Button label='Click Me' /></slot></Card>`
	expected := `<div class="card"><div><button>Click Me</button></div></div>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", normalizeHTML(expected), normalizeHTML(result))
	}
}

func TestSlot_OrderIndependence(t *testing.T) {
	state := createTestState(map[string]string{
		"Layout": `<div>{{ slot: header }}{{ slot: footer }}</div>`,
	})

	input := `<Layout><slot name='footer' tag='footer'>Footer Content</slot><slot name='header' tag='header'>Header Content</slot></Layout>`
	expected := `<div><header>Header Content</header><footer>Footer Content</footer></div>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", normalizeHTML(expected), normalizeHTML(result))
	}
}
