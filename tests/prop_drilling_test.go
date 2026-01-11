package main_test

import (
	"testing"

	"github.com/phillip-england/gtml/pkg/gtml"
)

func TestDrill_BasicDrilling(t *testing.T) {
	state := createTestState(map[string]string{
		"SomeComponent": `<div props='text string'><p>{text}</p></div>`,
		"SomeLayout":    `<div props='heading string'><h1>{heading}</h1><SomeComponent text={heading} /></div>`,
	})

	input := `<SomeLayout heading='some heading' />`
	expected := `<div><h1>some heading</h1><div><p>some heading</p></div></div>`

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{})
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

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{})
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

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{})
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

	result, err := gtml.CompileHTML(input, state, map[string]gtml.Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if normalizeHTML(result) != normalizeHTML(expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}
