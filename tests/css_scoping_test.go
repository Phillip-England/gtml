package main_test

import (
	"strings"
	"testing"

	"github.com/phillip-england/gtml/pkg/gtml"
)

func TestProcessComponentStyles_ExtractsCSS(t *testing.T) {
	input := `<style>
p { color: red; }
</style>
<div><p>Hello</p></div>`

	html, css, err := processComponentStyles(input, "data-test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(html, "<style>") {
		t.Error("HTML should not contain style block")
	}

	if !strings.Contains(html, "<div>") {
		t.Error("HTML should contain the div element")
	}

	if !strings.Contains(css, "p[data-test]") {
		t.Errorf("CSS should be scoped, got: %s", css)
	}
}

func TestProcessComponentStyles_NoStyle(t *testing.T) {
	input := `<div><p>Hello</p></div>`

	html, css, err := processComponentStyles(input, "data-test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if html != input {
		t.Errorf("HTML should be unchanged, got: %s", html)
	}

	if css != "" {
		t.Errorf("CSS should be empty, got: %s", css)
	}
}

func TestInjectScopeID(t *testing.T) {
	tests := []struct {
		input    string
		scopeID  string
		expected string
	}{
		{
			"<div></div>",
			"data-test",
			`<div data-test=""></div>`,
		},
		{
			"<div class='foo'></div>",
			"data-component",
			`<div data-component="" class='foo'></div>`,
		},
	}

	for _, tt := range tests {
		result := injectScopeID(tt.input, tt.scopeID)
		if result != tt.expected {
			t.Errorf("injectScopeID(%q, %q) = %q, expected %q", tt.input, tt.scopeID, result, tt.expected)
		}
	}
}

func processComponentStyles(raw string, scopeID string) (string, string, error) {
	return gtml.ProcessComponentStyles(raw, scopeID)
}

func injectScopeID(html string, scopeID string) string {
	return gtml.InjectScopeID(html, scopeID)
}
