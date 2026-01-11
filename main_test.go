package main

import (
	"strings"
	"testing"
)

//============================
// TEST HELPERS
//============================

// createTestState creates a GlobalState with the given components for testing
func createTestState(components map[string]string) *GlobalState {
	state := &GlobalState{
		Components: make(map[string]*Component),
	}
	for name, template := range components {
		scopeID := "data-" + strings.ToLower(name)

		// Parse props from template
		propDefs, cleanTemplate, _ := parsePropsAttribute(template)

		state.Components[name] = &Component{
			Name:     name,
			Template: cleanTemplate,
			ScopeID:  scopeID,
			PropDefs: propDefs,
		}
	}
	return state
}

// normalizeHTML removes extra whitespace for comparison
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

//============================
// PROPS TESTS
//============================

func TestProps_BasicStringProp(t *testing.T) {
	state := createTestState(map[string]string{
		"ThatButton": `<button props='text string'>{text}</button>`,
	})

	input := `<ThatButton text='some title' />`
	expected := `<button>some title</button>`

	result, err := compileHTML(input, state, map[string]Value{})
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

	result, err := compileHTML(input, state, map[string]Value{})
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

	result, err := compileHTML(input, state, map[string]Value{})
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

	result, err := compileHTML(input, state, map[string]Value{})
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

	result, err := compileHTML(input, state, map[string]Value{})
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

	result, err := compileHTML(input, state, map[string]Value{})
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

	result, err := compileHTML(input, state, map[string]Value{})
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

	result, err := compileHTML(input, state, map[string]Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(result, "props=") {
		t.Errorf("props attribute should be removed from output, got: %s", result)
	}
}

//============================
// PROP DRILLING TESTS
//============================

func TestDrill_BasicDrilling(t *testing.T) {
	state := createTestState(map[string]string{
		"SomeComponent": `<div props='text string'><p>{text}</p></div>`,
		"SomeLayout":    `<div props='heading string'><h1>{heading}</h1><SomeComponent text={heading} /></div>`,
	})

	input := `<SomeLayout heading='some heading' />`
	expected := `<div><h1>some heading</h1><div><p>some heading</p></div></div>`

	result, err := compileHTML(input, state, map[string]Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestDrill_MultiLevelDrilling(t *testing.T) {
	state := createTestState(map[string]string{
		"DeepComponent":   `<span props='value string'>{value}</span>`,
		"MiddleComponent": `<div props='data string'><DeepComponent value={data} /></div>`,
		"TopComponent":    `<section props='info string'><MiddleComponent data={info} /></section>`,
	})

	input := `<TopComponent info='drilled value' />`
	expected := `<section><div><span>drilled value</span></div></section>`

	result, err := compileHTML(input, state, map[string]Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestDrill_MultipleDrills(t *testing.T) {
	state := createTestState(map[string]string{
		"Display":   `<div props='title string, desc string'><h1>{title}</h1><p>{desc}</p></div>`,
		"Container": `<article props='heading string, description string'><Display title={heading} desc={description} /></article>`,
	})

	input := `<Container heading='My Title' description='My Description' />`
	expected := `<article><div><h1>My Title</h1><p>My Description</p></div></article>`

	result, err := compileHTML(input, state, map[string]Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestDrill_WithTransformation(t *testing.T) {
	state := createTestState(map[string]string{
		"Display": `<span props='num int'>{num}</span>`,
		"Wrapper": `<div props='base int'><Display num={base * 2} /></div>`,
	})

	input := `<Wrapper base={5} />`
	expected := `<div><span>10</span></div>`

	result, err := compileHTML(input, state, map[string]Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

//============================
// SLOT TESTS
//============================

func TestSlot_BasicSlot(t *testing.T) {
	state := createTestState(map[string]string{
		"PageLayout": `<html><body><header>Site Header</header>{{ slot: content }}<footer>Site Footer</footer></body></html>`,
	})

	input := `<PageLayout><slot name='content' tag='main'><p>Hello World</p></slot></PageLayout>`
	expected := `<html><body><header>Site Header</header><main><p>Hello World</p></main><footer>Site Footer</footer></body></html>`

	result, err := compileHTML(input, state, map[string]Value{})
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

	result, err := compileHTML(input, state, map[string]Value{})
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

	result, err := compileHTML(input, state, map[string]Value{})
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

	result, err := compileHTML(input, state, map[string]Value{})
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

	result, err := compileHTML(input, state, map[string]Value{})
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

	result, err := compileHTML(input, state, map[string]Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", normalizeHTML(expected), normalizeHTML(result))
	}
}

//============================
// TERNARY CONDITIONAL TESTS
//============================

func TestTernary_Basic(t *testing.T) {
	state := createTestState(map[string]string{
		"ShowIfActive": `<div props='active boolean'>{ active == true ? (<p>Active</p>) : (<p>Inactive</p>) }</div>`,
	})

	// Test true case
	input := `<ShowIfActive active={true} />`
	expected := `<div><p>Active</p></div>`

	result, err := compileHTML(input, state, map[string]Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}

	// Test false case
	input2 := `<ShowIfActive active={false} />`
	expected2 := `<div><p>Inactive</p></div>`

	result2, err := compileHTML(input2, state, map[string]Value{})
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

	result, err := compileHTML(input, state, map[string]Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}

	input2 := `<AgeCheck age={15} />`
	expected2 := `<div><p>Minor</p></div>`

	result2, err := compileHTML(input2, state, map[string]Value{})
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

	result, err := compileHTML(input, state, map[string]Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}

	input2 := `<ColorCheck color='red' />`
	expected2 := `<div><p>not blue</p></div>`

	result2, err := compileHTML(input2, state, map[string]Value{})
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

	result, err := compileHTML(input, state, map[string]Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}

	input2 := `<AccessCheck role='admin' active={false} />`
	expected2 := `<div><p>Limited Access</p></div>`

	result2, err := compileHTML(input2, state, map[string]Value{})
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

	result, err := compileHTML(input, state, map[string]Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}

	input2 := `<PriorityUser role='guest' />`
	expected2 := `<div><p>Standard User</p></div>`

	result2, err := compileHTML(input2, state, map[string]Value{})
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

	result, err := compileHTML(input, state, map[string]Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}

	input2 := `<NotBanned status='banned' />`
	expected2 := `<div><p>Account suspended</p></div>`

	result2, err := compileHTML(input2, state, map[string]Value{})
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

	result, err := compileHTML(input, state, map[string]Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

//============================
// EXPRESSION TESTS
//============================

func TestExpression_Addition(t *testing.T) {
	state := createTestState(map[string]string{
		"Math": `<span props='a int, b int'>{a + b}</span>`,
	})

	input := `<Math a={5} b={3} />`
	expected := `<span>8</span>`

	result, err := compileHTML(input, state, map[string]Value{})
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

	result, err := compileHTML(input, state, map[string]Value{})
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

	result, err := compileHTML(input, state, map[string]Value{})
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

	result, err := compileHTML(input, state, map[string]Value{})
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

	result, err := compileHTML(input, state, map[string]Value{})
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

	result, err := compileHTML(input, state, map[string]Value{})
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

	result, err := compileHTML(input, state, map[string]Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

//============================
// COMPONENT NOT FOUND ERROR TESTS
//============================

func TestError_ComponentNotFound(t *testing.T) {
	state := createTestState(map[string]string{})

	input := `<NonExistentComponent />`

	_, err := compileHTML(input, state, map[string]Value{})
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

	_, err := compileHTML(input, state, map[string]Value{})
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

	_, err := compileHTML(input, state, map[string]Value{})
	if err == nil {
		t.Error("expected error for type mismatch, got nil")
	}
}

//============================
// COMBINED FEATURES TESTS
//============================

func TestCombined_SlotWithDrilledProp(t *testing.T) {
	state := createTestState(map[string]string{
		"Button":    `<button props='text string'>{text}</button>`,
		"Container": `<div props='title string'>{title}{{ slot: content }}</div>`,
	})

	input := `<Container title='My Page'><slot name='content' tag='main'><Button text='Click' /></slot></Container>`
	expected := `<div>My Page<main><button>Click</button></main></div>`

	result, err := compileHTML(input, state, map[string]Value{})
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

	result, err := compileHTML(input, state, map[string]Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", normalizeHTML(expected), normalizeHTML(result))
	}
}

//============================
// HELPER FUNCTION TESTS
//============================

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
		{"Index", false},
		{"aboutUs", false},
		{"my_page", false},
		{"", false},
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
		{"  <div></div>  ", true},
	}

	for _, tt := range tests {
		result := hasSingleRoot(tt.input)
		if result != tt.expected {
			t.Errorf("hasSingleRoot(%q) = %v, expected %v", tt.input, result, tt.expected)
		}
	}
}

//============================
// CSS SCOPING TESTS
//============================

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

//============================
// PROPS PARSING TESTS
//============================

func TestParsePropsAttribute_Basic(t *testing.T) {
	input := `<div props='name string'>{name}</div>`
	propDefs, template, err := parsePropsAttribute(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(propDefs) != 1 {
		t.Errorf("expected 1 prop def, got %d", len(propDefs))
	}
	if propDefs["name"].Type != "string" {
		t.Errorf("expected type 'string', got '%s'", propDefs["name"].Type)
	}
	if strings.Contains(template, "props=") {
		t.Errorf("template should not contain props attribute, got: %s", template)
	}
}

func TestParsePropsAttribute_Multiple(t *testing.T) {
	input := `<div props='name string, age int, active boolean'>{name} is {age}</div>`
	propDefs, _, err := parsePropsAttribute(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(propDefs) != 3 {
		t.Errorf("expected 3 prop defs, got %d", len(propDefs))
	}
	if propDefs["name"].Type != "string" {
		t.Errorf("expected 'name' type 'string', got '%s'", propDefs["name"].Type)
	}
	if propDefs["age"].Type != "int" {
		t.Errorf("expected 'age' type 'int', got '%s'", propDefs["age"].Type)
	}
	if propDefs["active"].Type != "boolean" {
		t.Errorf("expected 'active' type 'boolean', got '%s'", propDefs["active"].Type)
	}
}

func TestParsePropsAttribute_InvalidType(t *testing.T) {
	input := `<div props='data unknown'>{data}</div>`
	_, _, err := parsePropsAttribute(input)
	if err == nil {
		t.Error("expected error for invalid prop type, got nil")
	}
}

func TestParsePropsAttribute_DuplicateName(t *testing.T) {
	input := `<div props='name string, name int'>{name}</div>`
	_, _, err := parsePropsAttribute(input)
	if err == nil {
		t.Error("expected error for duplicate prop name, got nil")
	}
}

//============================
// EDGE CASE TESTS
//============================

func TestEdge_SelfClosingWithSpaces(t *testing.T) {
	state := createTestState(map[string]string{
		"Icon": `<svg></svg>`,
	})

	input := `<Icon   />`
	expected := `<svg></svg>`

	result, err := compileHTML(input, state, map[string]Value{})
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

	result, err := compileHTML(input, state, map[string]Value{})
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

	result, err := compileHTML(input, state, map[string]Value{})
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

	result, err := compileHTML(input, state, map[string]Value{})
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

	result, err := compileHTML(input, state, map[string]Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

//============================
// EXPRESSION EVALUATOR UNIT TESTS
//============================

func TestEvaluateExpression_IntLiteral(t *testing.T) {
	val, err := evaluateExpression("42", map[string]Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val.Type != PropTypeInt || val.IntVal != 42 {
		t.Errorf("expected int 42, got %+v", val)
	}
}

func TestEvaluateExpression_StringLiteral(t *testing.T) {
	val, err := evaluateExpression("'hello'", map[string]Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val.Type != PropTypeString || val.StrVal != "hello" {
		t.Errorf("expected string 'hello', got %+v", val)
	}
}

func TestEvaluateExpression_BoolLiteral(t *testing.T) {
	val, err := evaluateExpression("true", map[string]Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val.Type != PropTypeBoolean || !val.BoolVal {
		t.Errorf("expected boolean true, got %+v", val)
	}
}

func TestEvaluateExpression_Variable(t *testing.T) {
	props := map[string]Value{
		"myVar": {Type: PropTypeString, StrVal: "test"},
	}
	val, err := evaluateExpression("myVar", props)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val.Type != PropTypeString || val.StrVal != "test" {
		t.Errorf("expected string 'test', got %+v", val)
	}
}

func TestEvaluateExpression_Comparison(t *testing.T) {
	props := map[string]Value{
		"x": {Type: PropTypeInt, IntVal: 5},
	}
	val, err := evaluateExpression("x > 3", props)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val.Type != PropTypeBoolean || !val.BoolVal {
		t.Errorf("expected boolean true, got %+v", val)
	}
}
