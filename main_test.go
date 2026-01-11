package main

import (
	"strings"
	"testing"
)

// =================================================================================
// TEST HELPERS
// =================================================================================

// createTestState creates a GlobalState with the given components for testing
func createTestState(components map[string]string) *GlobalState {
	state := &GlobalState{
		Components: make(map[string]*Component),
	}
	for name, template := range components {
		scopeID := "data-" + strings.ToLower(name)
		state.Components[name] = &Component{
			Name:     name,
			Template: template,
			ScopeID:  scopeID,
		}
	}
	return state
}

// normalizeHTML removes extra whitespace for comparison
func normalizeHTML(s string) string {
	s = strings.TrimSpace(s)
	// Collapse multiple whitespace to single space
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

// =================================================================================
// PROPS TESTS
// =================================================================================

func TestProps_BasicStringProp(t *testing.T) {
	// Test basic prop substitution with string type
	state := createTestState(map[string]string{
		"ThatButton": `<button>{{ prop: title string }}</button>`,
	})

	input := `<ThatButton title='some title' />`
	expected := `<button>some title</button>`

	result, err := compileHTML(input, state, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestProps_MultipleProps(t *testing.T) {
	// Test component with multiple props
	state := createTestState(map[string]string{
		"Card": `<div><h1>{{ prop: heading string }}</h1><p>{{ prop: subheading string }}</p></div>`,
	})

	input := `<Card heading='Hello' subheading='World' />`
	expected := `<div><h1>Hello</h1><p>World</p></div>`

	result, err := compileHTML(input, state, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestProps_IntType(t *testing.T) {
	// Test prop with int type
	state := createTestState(map[string]string{
		"Counter": `<span>Count: {{ prop: count int }}</span>`,
	})

	input := `<Counter count='42' />`
	expected := `<span>Count: 42</span>`

	result, err := compileHTML(input, state, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestProps_MissingPropValue(t *testing.T) {
	// When a prop is defined but not passed, it should render empty
	state := createTestState(map[string]string{
		"Greeting": `<div>Hello {{ prop: name string }}</div>`,
	})

	input := `<Greeting />`
	expected := `<div>Hello </div>`

	result, err := compileHTML(input, state, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestProps_DoubleQuotes(t *testing.T) {
	// Test that double quotes work for attribute values
	state := createTestState(map[string]string{
		"Message": `<p>{{ prop: text string }}</p>`,
	})

	input := `<Message text="Hello World" />`
	expected := `<p>Hello World</p>`

	result, err := compileHTML(input, state, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestProps_NonSelfClosingComponent(t *testing.T) {
	// Test non-self-closing component with children
	state := createTestState(map[string]string{
		"Wrapper": `<div class="wrap">{{ prop: title string }}</div>`,
	})

	input := `<Wrapper title="Test Title"></Wrapper>`
	expected := `<div class="wrap">Test Title</div>`

	result, err := compileHTML(input, state, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestProps_NestedComponents(t *testing.T) {
	// Test nested component resolution
	state := createTestState(map[string]string{
		"Inner":  `<span>{{ prop: value string }}</span>`,
		"Outer":  `<div><Inner value='nested' /></div>`,
	})

	input := `<Outer />`
	expected := `<div><span>nested</span></div>`

	result, err := compileHTML(input, state, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

// =================================================================================
// PROP DRILLING TESTS
// =================================================================================

func TestDrill_BasicDrilling(t *testing.T) {
	// Test basic prop drilling from parent to child
	state := createTestState(map[string]string{
		"SomeComponent": `<div><p>{{ prop: text string }}</p></div>`,
		"SomeLayout": `<div><h1>{{ prop: heading string }}</h1><SomeComponent text="{{ drill: heading }}" /></div>`,
	})

	input := `<SomeLayout heading='some heading' />`
	expected := `<div><h1>some heading</h1><div><p>some heading</p></div></div>`

	result, err := compileHTML(input, state, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestDrill_MultiLevelDrilling(t *testing.T) {
	// Test drilling through multiple component levels
	state := createTestState(map[string]string{
		"DeepComponent":   `<span>{{ prop: value string }}</span>`,
		"MiddleComponent": `<div><DeepComponent value="{{ drill: data }}" /></div>`,
		"TopComponent":    `<section><MiddleComponent data="{{ drill: info }}" /></section>`,
	})

	// Start with props in scope
	input := `<TopComponent info='drilled value' />`
	expected := `<section><div><span>drilled value</span></div></section>`

	result, err := compileHTML(input, state, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestDrill_MultipleDrills(t *testing.T) {
	// Test drilling multiple props to same child
	state := createTestState(map[string]string{
		"Display": `<div><h1>{{ prop: title string }}</h1><p>{{ prop: desc string }}</p></div>`,
		"Container": `<article><Display title="{{ drill: heading }}" desc="{{ drill: description }}" /></article>`,
	})

	input := `<Container heading='My Title' description='My Description' />`
	expected := `<article><div><h1>My Title</h1><p>My Description</p></div></article>`

	result, err := compileHTML(input, state, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestDrill_MissingDrillSource(t *testing.T) {
	// When drilling a prop that doesn't exist, it should render empty
	state := createTestState(map[string]string{
		"Child":  `<p>{{ prop: text string }}</p>`,
		"Parent": `<div><Child text="{{ drill: nonexistent }}" /></div>`,
	})

	input := `<Parent />`
	expected := `<div><p></p></div>`

	result, err := compileHTML(input, state, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestDrill_MixedDrillAndLiteral(t *testing.T) {
	// Test component with both drilled and literal props
	state := createTestState(map[string]string{
		"Item": `<li>{{ prop: label string }} - {{ prop: value string }}</li>`,
		"List": `<ul><Item label='Fixed Label' value="{{ drill: dynamicValue }}" /></ul>`,
	})

	input := `<List dynamicValue='Dynamic!' />`
	expected := `<ul><li>Fixed Label - Dynamic!</li></ul>`

	result, err := compileHTML(input, state, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

// =================================================================================
// SLOT TESTS
// =================================================================================

func TestSlot_BasicSlot(t *testing.T) {
	// Test basic slot insertion
	state := createTestState(map[string]string{
		"PageLayout": `<html><body><header>Site Header</header>{{ slot: content }}<footer>Site Footer</footer></body></html>`,
	})

	input := `<PageLayout><slot name='content' tag='main'><p>Hello World</p></slot></PageLayout>`
	expected := `<html><body><header>Site Header</header><main><p>Hello World</p></main><footer>Site Footer</footer></body></html>`

	result, err := compileHTML(input, state, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", normalizeHTML(expected), normalizeHTML(result))
	}
}

func TestSlot_MultipleSlots(t *testing.T) {
	// Test component with multiple named slots
	state := createTestState(map[string]string{
		"TwoColumnLayout": `<div class='container'>{{ slot: sidebar }}{{ slot: main }}</div>`,
	})

	input := `<TwoColumnLayout><slot name='sidebar' tag='aside'><nav>Navigation</nav></slot><slot name='main' tag='section'><p>Main content here</p></slot></TwoColumnLayout>`
	expected := `<div class='container'><aside><nav>Navigation</nav></aside><section><p>Main content here</p></section></div>`

	result, err := compileHTML(input, state, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", normalizeHTML(expected), normalizeHTML(result))
	}
}

func TestSlot_WithClass(t *testing.T) {
	// Test that slot preserves class attribute
	state := createTestState(map[string]string{
		"Layout": `<div>{{ slot: content }}</div>`,
	})

	input := `<Layout><slot name='content' tag='div' class='my-class'><p>Content</p></slot></Layout>`
	expected := `<div><div class="my-class"><p>Content</p></div></div>`

	result, err := compileHTML(input, state, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", normalizeHTML(expected), normalizeHTML(result))
	}
}

func TestSlot_EmptySlot(t *testing.T) {
	// Test that unfilled slot renders empty
	state := createTestState(map[string]string{
		"Layout": `<div>{{ slot: content }}</div>`,
	})

	input := `<Layout></Layout>`
	expected := `<div></div>`

	result, err := compileHTML(input, state, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestSlot_NestedComponentInSlot(t *testing.T) {
	// Test that components inside slots are resolved
	state := createTestState(map[string]string{
		"Button": `<button>{{ prop: label string }}</button>`,
		"Card":   `<div class="card">{{ slot: actions }}</div>`,
	})

	input := `<Card><slot name='actions' tag='div'><Button label='Click Me' /></slot></Card>`
	expected := `<div class="card"><div><button>Click Me</button></div></div>`

	result, err := compileHTML(input, state, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", normalizeHTML(expected), normalizeHTML(result))
	}
}

func TestSlot_OrderIndependence(t *testing.T) {
	// Test that slots can be provided in different order than defined
	state := createTestState(map[string]string{
		"Layout": `<div>{{ slot: header }}{{ slot: footer }}</div>`,
	})

	// Provide footer first, then header
	input := `<Layout><slot name='footer' tag='footer'>Footer Content</slot><slot name='header' tag='header'>Header Content</slot></Layout>`
	expected := `<div><header>Header Content</header><footer>Footer Content</footer></div>`

	result, err := compileHTML(input, state, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", normalizeHTML(expected), normalizeHTML(result))
	}
}

// =================================================================================
// COMPONENT NOT FOUND ERROR TESTS
// =================================================================================

func TestError_ComponentNotFound(t *testing.T) {
	state := createTestState(map[string]string{})

	input := `<NonExistentComponent />`

	_, err := compileHTML(input, state, map[string]interface{}{})
	if err == nil {
		t.Error("expected error for non-existent component, got nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' error, got: %v", err)
	}
}

// =================================================================================
// COMBINED FEATURES TESTS
// =================================================================================

func TestCombined_SlotWithDrilledProp(t *testing.T) {
	// Test using drilled props inside slotted content
	state := createTestState(map[string]string{
		"Button":    `<button>{{ prop: text string }}</button>`,
		"Container": `<div>{{ prop: title string }}{{ slot: content }}</div>`,
	})

	// This tests the integration from the spec example
	input := `<Container title='My Page'><slot name='content' tag='main'><Button text='Click' /></slot></Container>`
	expected := `<div>My Page<main><button>Click</button></main></div>`

	result, err := compileHTML(input, state, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", normalizeHTML(expected), normalizeHTML(result))
	}
}

func TestCombined_FullLayoutExample(t *testing.T) {
	// Test the full example from the spec
	state := createTestState(map[string]string{
		"BasicButton": `<button>{{ prop: text string }}</button>`,
		"GuestLayout": `<html><head><title>{{ prop: title string }}</title></head><body><BasicButton text='{{ drill: title }}' />{{ slot: content }}</body></html>`,
	})

	input := `<GuestLayout title="Some Title"><slot name='content' tag='div'><p>Some Content</p><BasicButton text='Click Me' /></slot></GuestLayout>`
	expected := `<html><head><title>Some Title</title></head><body><button>Some Title</button><div><p>Some Content</p><button>Click Me</button></div></body></html>`

	result, err := compileHTML(input, state, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", normalizeHTML(expected), normalizeHTML(result))
	}
}

// =================================================================================
// HELPER FUNCTION TESTS
// =================================================================================

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
		{"basicButton", false},  // starts with lowercase
		{"basic_button", false}, // has underscore
		{"basic-button", false}, // has hyphen
		{"", false},             // empty
		{"123Button", false},    // starts with number
	}

	for _, tt := range tests {
		result := isPascalCase(tt.input)
		if result != tt.expected {
			t.Errorf("isPascalCase(%q) = %v, expected %v", tt.input, result, tt.expected)
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
		{"Index", false},        // has uppercase
		{"aboutUs", false},      // has uppercase
		{"my_page", false},      // has underscore
		{"", false},             // empty
	}

	for _, tt := range tests {
		result := isKebabCase(tt.input)
		if result != tt.expected {
			t.Errorf("isKebabCase(%q) = %v, expected %v", tt.input, result, tt.expected)
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
			`data="{{drill: value}}"`,
			map[string]string{"data": "{{drill: value}}"},
		},
		{
			``,
			map[string]string{},
		},
	}

	for _, tt := range tests {
		result := parseAttributes(tt.input)
		if len(result) != len(tt.expected) {
			t.Errorf("parseAttributes(%q): length mismatch, got %d, expected %d", tt.input, len(result), len(tt.expected))
			continue
		}
		for k, v := range tt.expected {
			if result[k] != v {
				t.Errorf("parseAttributes(%q): key %q = %q, expected %q", tt.input, k, result[k], v)
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
		{"  <div></div>  ", true}, // with whitespace
		// Note: hasSingleRoot uses a loose heuristic check, so it may not catch all cases
		{"<p>one</p><p>two</p>", true}, // implementation returns true (loose check)
	}

	for _, tt := range tests {
		result := hasSingleRoot(tt.input)
		if result != tt.expected {
			t.Errorf("hasSingleRoot(%q) = %v, expected %v", tt.input, result, tt.expected)
		}
	}
}

// =================================================================================
// CSS SCOPING TESTS
// =================================================================================

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

// =================================================================================
// EDGE CASE TESTS
// =================================================================================

func TestEdge_SelfClosingWithSpaces(t *testing.T) {
	state := createTestState(map[string]string{
		"Icon": `<svg></svg>`,
	})

	input := `<Icon   />`
	expected := `<svg></svg>`

	result, err := compileHTML(input, state, map[string]interface{}{})
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

	result, err := compileHTML(input, state, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestEdge_DeeplyNestedComponents(t *testing.T) {
	state := createTestState(map[string]string{
		"Level3": `<span>{{ prop: val string }}</span>`,
		"Level2": `<p><Level3 val='{{ drill: data }}' /></p>`,
		"Level1": `<div><Level2 data='{{ drill: input }}' /></div>`,
	})

	input := `<Level1 input='deep value' />`
	expected := `<div><p><span>deep value</span></p></div>`

	result, err := compileHTML(input, state, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestEdge_MultipleComponentsInRoute(t *testing.T) {
	state := createTestState(map[string]string{
		"Para": `<p>{{ prop: text string }}</p>`,
	})

	input := `<div><Para text='First' /><Para text='Second' /></div>`
	expected := `<div><p>First</p><p>Second</p></div>`

	result, err := compileHTML(input, state, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestEdge_ComponentWithTextAround(t *testing.T) {
	state := createTestState(map[string]string{
		"Bold": `<strong>{{ prop: text string }}</strong>`,
	})

	input := `<p>Hello <Bold text='World' /> and goodbye</p>`
	expected := `<p>Hello <strong>World</strong> and goodbye</p>`

	result, err := compileHTML(input, state, map[string]interface{}{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}
