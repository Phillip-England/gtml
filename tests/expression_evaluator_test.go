package main_test

import (
	"testing"

	"github.com/phillip-england/gtml/pkg/gtml"
)

func TestEvaluateExpression_IntLiteral(t *testing.T) {
	val, err := gtml.EvaluateExpression("42", map[string]gtml.Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val.Type != gtml.PropTypeInt || val.IntVal != 42 {
		t.Errorf("expected int 42, got %+v", val)
	}
}

func TestEvaluateExpression_StringLiteral(t *testing.T) {
	val, err := gtml.EvaluateExpression("'hello'", map[string]gtml.Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val.Type != gtml.PropTypeString || val.StrVal != "hello" {
		t.Errorf("expected string 'hello', got %+v", val)
	}
}

func TestEvaluateExpression_BoolLiteral(t *testing.T) {
	val, err := gtml.EvaluateExpression("true", map[string]gtml.Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val.Type != gtml.PropTypeBoolean || !val.BoolVal {
		t.Errorf("expected boolean true, got %+v", val)
	}
}

func TestEvaluateExpression_Variable(t *testing.T) {
	props := map[string]gtml.Value{
		"myVar": {Type: gtml.PropTypeString, StrVal: "test"},
	}
	val, err := gtml.EvaluateExpression("myVar", props)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val.Type != gtml.PropTypeString || val.StrVal != "test" {
		t.Errorf("expected string 'test', got %+v", val)
	}
}

func TestEvaluateExpression_Comparison(t *testing.T) {
	props := map[string]gtml.Value{
		"x": {Type: gtml.PropTypeInt, IntVal: 5},
	}
	val, err := gtml.EvaluateExpression("x > 3", props)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val.Type != gtml.PropTypeBoolean || !val.BoolVal {
		t.Errorf("expected boolean true, got %+v", val)
	}
}
