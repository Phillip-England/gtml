package main_test

import (
	"testing"

	"github.com/phillip-england/gtml/pkg/gtml"
)

func TestCombined_SlotWithDrilledProp(t *testing.T) {
	state := createTestState(map[string]string{
		"Button":    `<button props='text string'>{text}</button>`,
		"Container": `<div props='title string'>{title}{{ slot: content }}</div>`,
	})

	input := `<Container title='My Page'><slot name='content' tag='main'><Button text='Click' /></slot></Container>`
	expected := `<div>My Page<main><button>Click</button></main></div>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", normalizeHTML(expected), normalizeHTML(result))
	}
}

func TestCombined_FullLayoutExample(t *testing.T) {
	state := createTestState(map[string]string{
		"BasicButton": `<button props='text string'>{text}</button>`,
		"GuestLayout": `<html props='title string'><head><title>{title}</title></head><body><BasicButton text={title} />{{ slot: content }}</body></html>`,
	})

	input := `<GuestLayout title="Some Title"><slot name='content' tag='div'><p>Some Content</p><BasicButton text='Click Me' /></slot></GuestLayout>`
	expected := `<html><head><title>Some Title</title></head><body><button>Some Title</button><div><p>Some Content</p><button>Click Me</button></div></body></html>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", normalizeHTML(expected), normalizeHTML(result))
	}
}
