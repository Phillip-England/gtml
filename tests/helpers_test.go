package main_test

import (
	"strings"
	"testing"

	"github.com/phillip-england/gtml/pkg/gtml"
)

func createTestState(components map[string]string) *gtml.GlobalState {
	state := &gtml.GlobalState{
		Components: make(map[string]*gtml.Component),
	}
	for name, template := range components {
		scopeID := "data-" + strings.ToLower(name)

		propDefs, cleanTemplate, _ := gtml.ParsePropsAttribute(template)

		state.Components[name] = &gtml.Component{
			Name:     name,
			Template: cleanTemplate,
			ScopeID:  scopeID,
			PropDefs: propDefs,
		}
	}
	return state
}

func normalizeHTML(s string) string {
	s = strings.TrimSpace(s)
	var result strings.Builder
	inWhitespace := false
	for _, r := range s {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			if !inWhitespace {
				result.WriteRune(' ')
				inWhitespace = true
			}
		} else {
			result.WriteRune(r)
			inWhitespace = false
		}
	}
	return result.String()
}

func TestIsPascalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"BasicButton", true},
		{"A", true},
		{"Ab", true},
		{"AB", true},
		{"Button123", true},
		{"basicButton", false},
		{"basic_button", false},
		{"basic-button", false},
		{"", false},
		{"123Button", false},
	}

	for _, tt := range tests {
		result := gtml.IsPascalCase(tt.input)
		if result != tt.expected {
			t.Errorf("IsPascalCase(%q) = %v, expected %v", tt.input, result, tt.expected)
		}
	}
}

func TestIsKebabCase(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"index", true},
		{"about-us", true},
		{"my-page-123", true},
		{"a", true},
		{"Index", false},
		{"aboutUs", false},
		{"my_page", false},
		{"", false},
	}

	for _, tt := range tests {
		result := gtml.IsKebabCase(tt.input)
		if result != tt.expected {
			t.Errorf("IsKebabCase(%q) = %v, expected %v", tt.input, result, tt.expected)
		}
	}
}

func TestParseAttributes(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]string
	}{
		{
			`title="Hello" name='World'`,
			map[string]string{"title": "Hello", "name": "World"},
		},
		{
			`class="my-class"`,
			map[string]string{"class": "my-class"},
		},
		{
			``,
			map[string]string{},
		},
	}

	for _, tt := range tests {
		result := gtml.ParseAttributes(tt.input)
		if len(result) != len(tt.expected) {
			t.Errorf("ParseAttributes(%q): length mismatch, got %d, expected %d", tt.input, len(result), len(tt.expected))
			continue
		}
		for k, v := range tt.expected {
			if result[k] != v {
				t.Errorf("ParseAttributes(%q): key %q = %q, expected %q", tt.input, k, result[k], v)
			}
		}
	}
}

func TestHasSingleRoot(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"<div></div>", true},
		{"<div><p>content</p></div>", true},
		{"<button />", true},
		{"<br />", true},
		{"  <div></div>  ", true},
	}

	for _, tt := range tests {
		result := gtml.HasSingleRoot(tt.input)
		if result != tt.expected {
			t.Errorf("HasSingleRoot(%q) = %v, expected %v", tt.input, result, tt.expected)
		}
	}
}
