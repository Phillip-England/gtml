package main_test

import (
	"strings"
	"testing"

	"github.com/phillip-england/gtml/pkg/gtml"
)

func TestError_ComponentNotFound(t *testing.T) {
	state := createTestState(map[string]string{})

	input := `<NonExistentComponent />`

	_, err := gtml.CompileHTML(input, state, map[string]gtml.Value{})
	if err == nil {
		t.Error("expected error for non-existent component, got nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' error, got: %v", err)
	}
}

func TestError_UndefinedVariable(t *testing.T) {
	state := createTestState(map[string]string{
		"Test": `<span props='x int'>{undefinedVar}</span>`,
	})

	input := `<Test x={1} />`

	_, err := gtml.CompileHTML(input, state, map[string]gtml.Value{})
	if err == nil {
		t.Error("expected error for undefined variable, got nil")
	}
	if !strings.Contains(err.Error(), "undefined") {
		t.Errorf("expected 'undefined' error, got: %v", err)
	}
}

func TestError_TypeMismatch(t *testing.T) {
	state := createTestState(map[string]string{
		"Counter": `<span props='count int'>{count}</span>`,
	})

	input := `<Counter count='not a number' />`

	_, err := gtml.CompileHTML(input, state, map[string]gtml.Value{})
	if err == nil {
		t.Error("expected error for type mismatch, got nil")
	}
}
