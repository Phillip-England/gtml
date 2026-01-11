package gtml

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	PropTypeString  = "string"
	PropTypeInt     = "int"
	PropTypeBoolean = "boolean"
)

var (
	reComponentTag    = regexp.MustCompile(`</?([A-Z][a-zA-Z0-9]*)`)
	reStyleBlock      = regexp.MustCompile(`(?s)<style>(.*?)</style>`)
	rePropsAttr       = regexp.MustCompile(`\s+props\s*=\s*['"]([^'"]+)['"]`)
	reExpression      = regexp.MustCompile(`\{([^{}]+)\}`)
	reSlotPlaceholder = regexp.MustCompile(`(?s)<slot\s+name=['"](\w+)['"]\s*/?>`)
	reSlotUsage       = regexp.MustCompile(`(?s)<slot\s+([^>]+)>(.*?)</slot>`)
	reSlotDef         = regexp.MustCompile(`\{\{\s*slot:\s*(\w+)\s*\}\}`)
)

type PropDef struct {
	Name string
	Type string
}

type Component struct {
	Name        string
	RawContent  string
	Template    string
	Styles      string
	ScopedStyle string
	ScopeID     string
	Path        string
	PropDefs    map[string]PropDef
}

type GlobalState struct {
	Components map[string]*Component
	CSSOutput  strings.Builder
}

type Value struct {
	Type    string
	StrVal  string
	IntVal  int
	BoolVal bool
}

func (v Value) String() string {
	switch v.Type {
	case PropTypeString:
		return v.StrVal
	case PropTypeInt:
		return strconv.Itoa(v.IntVal)
	case PropTypeBoolean:
		if v.BoolVal {
			return "true"
		}
		return "false"
	}
	return ""
}

func CompileHTML(html string, state *GlobalState, scopeProps map[string]Value) (string, error) {
	var err error
	html, err = evaluateTernaries(html, scopeProps)
	if err != nil {
		return "", err
	}

	for {
		startIdx, endIdx, tagName, _, attrsStr, innerContent := findFirstComponent(html)
		if startIdx == -1 {
			break
		}

		compDef, exists := state.Components[tagName]
		if !exists {
			return "", fmt.Errorf("component '%s' not found", tagName)
		}

		props, err := parseComponentAttributes(attrsStr, scopeProps, compDef.PropDefs)
		if err != nil {
			return "", fmt.Errorf("error parsing attributes for %s: %v", tagName, err)
		}

		compiledChildren, err := CompileHTML(innerContent, state, scopeProps)
		if err != nil {
			return "", err
		}

		slotsMap := extractSlots(compiledChildren)

		renderedComp := compDef.Template

		renderedComp = reSlotDef.ReplaceAllStringFunc(renderedComp, func(match string) string {
			subMatch := reSlotDef.FindStringSubmatch(match)
			if len(subMatch) < 2 {
				return match
			}
			slotName := subMatch[1]
			return "<slot name='" + slotName + "' />"
		})

		renderedComp, err = EvaluateExpressions(renderedComp, props)
		if err != nil {
			return "", fmt.Errorf("error evaluating expressions in %s: %v", tagName, err)
		}

		renderedComp = reSlotPlaceholder.ReplaceAllStringFunc(renderedComp, func(match string) string {
			subMatch := reSlotPlaceholder.FindStringSubmatch(match)
			if len(subMatch) < 2 {
				return ""
			}
			slotName := subMatch[1]
			if content, ok := slotsMap[slotName]; ok {
				return content
			}
			return ""
		})

		finalRendered, err := CompileHTML(renderedComp, state, props)
		if err != nil {
			return "", err
		}

		html = html[:startIdx] + finalRendered + html[endIdx:]
	}

	html, err = EvaluateExpressions(html, scopeProps)
	if err != nil {
		return "", err
	}

	return html, nil
}

func extractSlots(content string) map[string]string {
	slots := make(map[string]string)
	matches := reSlotUsage.FindAllStringSubmatchIndex(content, -1)

	for _, loc := range matches {
		attrsStr := content[loc[2]:loc[3]]
		innerContent := content[loc[4]:loc[5]]

		attrs := ParseAttributes(attrsStr)
		name := attrs["name"]
		tag := attrs["tag"]

		if name != "" && tag != "" {
			wrapper := "<" + tag
			for k, v := range attrs {
				if k != "name" && k != "tag" {
					wrapper += fmt.Sprintf(" %s='%s'", k, v)
				}
			}
			wrapper += ">" + innerContent + "</" + tag + ">"
			slots[name] = wrapper
		}
	}
	return slots
}

func EvaluateExpressions(html string, props map[string]Value) (string, error) {
	result := html
	offset := 0
	for {
		match := reExpression.FindStringSubmatchIndex(result[offset:])
		if match == nil {
			break
		}

		fullStart := offset + match[0]
		fullEnd := offset + match[1]
		exprStart := offset + match[2]
		exprEnd := offset + match[3]

		expr := result[exprStart:exprEnd]

		if isInsideComponentTag(result, fullStart) {
			offset = fullEnd
			continue
		}

		if strings.Contains(expr, "?") && strings.Contains(expr, "(") {
			offset = fullEnd
			continue
		}

		value, err := EvaluateExpression(expr, props)
		if err != nil {
			return "", err
		}

		result = result[:fullStart] + value.String() + result[fullEnd:]
		offset = fullStart + len(value.String())
	}
	return result, nil
}

func isInsideComponentTag(html string, pos int) bool {
	for i := pos - 1; i >= 0; i-- {
		if html[i] == '>' {
			return false
		}
		if html[i] == '<' {
			if i+1 < len(html) {
				nextChar := html[i+1]
				if nextChar >= 'A' && nextChar <= 'Z' {
					return true
				}
			}
			return false
		}
	}
	return false
}

func evaluateTernaries(html string, props map[string]Value) (string, error) {
	result := html
	for {
		ternaryStart := findTernaryStart(result)
		if ternaryStart == -1 {
			break
		}

		ternaryEnd, condition, truthy, falsy, err := parseTernary(result, ternaryStart)
		if err != nil {
			return "", err
		}

		condValue, err := EvaluateExpression(condition, props)
		if err != nil {
			return "", err
		}

		var replacement string
		if condValue.Type == PropTypeBoolean {
			if condValue.BoolVal {
				replacement = truthy
			} else {
				replacement = falsy
			}
		} else {
			return "", fmt.Errorf("ternary condition '%s' must evaluate to boolean, got %s", condition, condValue.Type)
		}

		replacement, err = evaluateTernaries(replacement, props)
		if err != nil {
			return "", err
		}

		result = result[:ternaryStart] + replacement + result[ternaryEnd:]
	}
	return result, nil
}

func findTernaryStart(s string) int {
	for i := 0; i < len(s); i++ {
		if s[i] == '{' {
			depth := 1
			hasQuestion := false
			hasParen := false
			for j := i + 1; j < len(s); j++ {
				if s[j] == '{' {
					depth++
				} else if s[j] == '}' {
					depth--
					if depth == 0 {
						break
					}
				} else if depth == 1 {
					if s[j] == '?' {
						hasQuestion = true
					} else if s[j] == '(' && hasQuestion {
						hasParen = true
						return i
					}
				}
			}
			if !hasQuestion || !hasParen {
				continue
			}
		}
	}
	return -1
}

func parseTernary(s string, pos int) (int, string, string, string, error) {
	if s[pos] != '{' {
		return 0, "", "", "", fmt.Errorf("expected '{' at position %d", pos)
	}

	conditionEnd := -1
	depth := 0
	for i := pos + 1; i < len(s); i++ {
		if s[i] == '(' {
			depth++
		} else if s[i] == ')' {
			depth--
		} else if s[i] == '?' && depth == 0 {
			conditionEnd = i
			break
		} else if s[i] == '}' && depth == 0 {
			return 0, "", "", "", fmt.Errorf("malformed ternary, missing '?'")
		}
	}

	if conditionEnd == -1 {
		return 0, "", "", "", fmt.Errorf("missing '?' in ternary expression")
	}

	condition := strings.TrimSpace(s[pos+1 : conditionEnd])

	truthyStart := -1
	for i := conditionEnd + 1; i < len(s); i++ {
		if s[i] == '(' {
			truthyStart = i
			break
		} else if !unicode.IsSpace(rune(s[i])) {
			return 0, "", "", "", fmt.Errorf("expected '(' after '?' in ternary")
		}
	}
	truthyEnd := findMatchingParen(s, truthyStart)
	if truthyEnd == -1 {
		return 0, "", "", "", fmt.Errorf("unbalanced parentheses in truthy branch")
	}
	truthy := s[truthyStart+1 : truthyEnd]

	colonPos := -1
	for i := truthyEnd + 1; i < len(s); i++ {
		if s[i] == ':' {
			colonPos = i
			break
		} else if !unicode.IsSpace(rune(s[i])) {
			return 0, "", "", "", fmt.Errorf("expected ':' after truthy branch")
		}
	}
	if colonPos == -1 {
		return 0, "", "", "", fmt.Errorf("missing ':' in ternary expression")
	}

	falsyStart := -1
	for i := colonPos + 1; i < len(s); i++ {
		if s[i] == '(' {
			falsyStart = i
			break
		} else if !unicode.IsSpace(rune(s[i])) {
			return 0, "", "", "", fmt.Errorf("expected '(' after ':' in ternary")
		}
	}
	falsyEnd := findMatchingParen(s, falsyStart)
	if falsyEnd == -1 {
		return 0, "", "", "", fmt.Errorf("unbalanced parentheses in falsy branch")
	}
	falsy := s[falsyStart+1 : falsyEnd]

	closingBrace := -1
	for i := falsyEnd + 1; i < len(s); i++ {
		if s[i] == '}' {
			closingBrace = i
			break
		} else if !unicode.IsSpace(rune(s[i])) {
			return 0, "", "", "", fmt.Errorf("expected '}' after falsy branch")
		}
	}

	return closingBrace + 1, condition, truthy, falsy, nil
}

func findMatchingParen(s string, start int) int {
	if s[start] != '(' {
		return -1
	}
	depth := 1
	for i := start + 1; i < len(s); i++ {
		if s[i] == '(' {
			depth++
		} else if s[i] == ')' {
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func EvaluateExpression(expr string, props map[string]Value) (Value, error) {
	expr = strings.TrimSpace(expr)

	if expr == "true" {
		return Value{Type: PropTypeBoolean, BoolVal: true}, nil
	}
	if expr == "false" {
		return Value{Type: PropTypeBoolean, BoolVal: false}, nil
	}

	if (strings.HasPrefix(expr, "'") && strings.HasSuffix(expr, "'")) ||
		(strings.HasPrefix(expr, "\"") && strings.HasSuffix(expr, "\"")) {
		if len(expr) < 2 {
			return Value{Type: PropTypeString, StrVal: ""}, nil
		}
		return Value{Type: PropTypeString, StrVal: expr[1 : len(expr)-1]}, nil
	}

	if i, err := strconv.Atoi(expr); err == nil {
		return Value{Type: PropTypeInt, IntVal: i}, nil
	}
	if strings.HasPrefix(expr, "-") {
		rest := strings.TrimSpace(expr[1:])
		if i, err := strconv.Atoi(rest); err == nil {
			return Value{Type: PropTypeInt, IntVal: -i}, nil
		}
	}

	if idx := findOperator(expr, "||"); idx != -1 {
		left, err := EvaluateExpression(expr[:idx], props)
		if err != nil {
			return Value{}, err
		}
		right, err := EvaluateExpression(expr[idx+2:], props)
		if err != nil {
			return Value{}, err
		}
		if left.Type != PropTypeBoolean || right.Type != PropTypeBoolean {
			return Value{}, fmt.Errorf("|| operator requires boolean operands")
		}
		return Value{Type: PropTypeBoolean, BoolVal: left.BoolVal || right.BoolVal}, nil
	}
	if idx := findOperator(expr, "&&"); idx != -1 {
		left, err := EvaluateExpression(expr[:idx], props)
		if err != nil {
			return Value{}, err
		}
		right, err := EvaluateExpression(expr[idx+2:], props)
		if err != nil {
			return Value{}, err
		}
		if left.Type != PropTypeBoolean || right.Type != PropTypeBoolean {
			return Value{}, fmt.Errorf("&& operator requires boolean operands")
		}
		return Value{Type: PropTypeBoolean, BoolVal: left.BoolVal && right.BoolVal}, nil
	}

	for _, op := range []string{"==", "!=", "<=", ">=", "<", ">"} {
		if idx := findOperator(expr, op); idx != -1 {
			left, err := EvaluateExpression(expr[:idx], props)
			if err != nil {
				return Value{}, err
			}
			right, err := EvaluateExpression(expr[idx+len(op):], props)
			if err != nil {
				return Value{}, err
			}
			return compareValues(left, right, op)
		}
	}

	if idx := findOperatorRTL(expr, "+", "-"); idx != -1 {
		op := string(expr[idx])
		left, err := EvaluateExpression(expr[:idx], props)
		if err != nil {
			return Value{}, err
		}
		right, err := EvaluateExpression(expr[idx+1:], props)
		if err != nil {
			return Value{}, err
		}

		if op == "+" {
			if left.Type == PropTypeString && right.Type == PropTypeString {
				return Value{Type: PropTypeString, StrVal: left.StrVal + right.StrVal}, nil
			}
			if left.Type == PropTypeInt && right.Type == PropTypeInt {
				return Value{Type: PropTypeInt, IntVal: left.IntVal + right.IntVal}, nil
			}
			return Value{}, fmt.Errorf("+ operator requires matching types (string+string or int+int)")
		} else {
			if left.Type != PropTypeInt || right.Type != PropTypeInt {
				return Value{}, fmt.Errorf("- operator requires int operands")
			}
			return Value{Type: PropTypeInt, IntVal: left.IntVal - right.IntVal}, nil
		}
	}
	if idx := findOperatorRTL(expr, "*", "/", "%"); idx != -1 {
		op := string(expr[idx])
		left, err := EvaluateExpression(expr[:idx], props)
		if err != nil {
			return Value{}, err
		}
		right, err := EvaluateExpression(expr[idx+1:], props)
		if err != nil {
			return Value{}, err
		}

		if left.Type != PropTypeInt || right.Type != PropTypeInt {
			return Value{}, fmt.Errorf("%s operator requires int operands", op)
		}
		switch op {
		case "*":
			return Value{Type: PropTypeInt, IntVal: left.IntVal * right.IntVal}, nil
		case "/":
			if right.IntVal == 0 {
				return Value{}, fmt.Errorf("division by zero")
			}
			return Value{Type: PropTypeInt, IntVal: left.IntVal / right.IntVal}, nil
		case "%":
			if right.IntVal == 0 {
				return Value{}, fmt.Errorf("modulo by zero")
			}
			return Value{Type: PropTypeInt, IntVal: left.IntVal % right.IntVal}, nil
		}
	}

	if strings.HasPrefix(expr, "(") && strings.HasSuffix(expr, ")") {
		return EvaluateExpression(expr[1:len(expr)-1], props)
	}

	if isValidIdentifier(expr) {
		if val, ok := props[expr]; ok {
			return val, nil
		}
		return Value{}, fmt.Errorf("undefined variable: %s", expr)
	}

	return Value{}, fmt.Errorf("invalid expression: %s", expr)
}

func findOperator(expr string, op string) int {
	depth := 0
	inStr := false
	strChar := byte(0)
	for i := 0; i < len(expr)-len(op)+1; i++ {
		if !inStr && (expr[i] == '"' || expr[i] == '\'') {
			inStr = true
			strChar = expr[i]
		} else if inStr && expr[i] == strChar {
			inStr = false
		} else if !inStr {
			if expr[i] == '(' {
				depth++
			} else if expr[i] == ')' {
				depth--
			} else if depth == 0 {
				if expr[i:i+len(op)] == op {
					return i
				}
			}
		}
	}
	return -1
}

func findOperatorRTL(expr string, ops ...string) int {
	depth := 0
	inStr := false
	strChar := byte(0)
	for i := len(expr) - 1; i >= 0; i-- {
		if !inStr && (expr[i] == '"' || expr[i] == '\'') {
			inStr = true
			strChar = expr[i]
		} else if inStr && expr[i] == strChar {
			inStr = false
		} else if !inStr {
			if expr[i] == ')' {
				depth++
			} else if expr[i] == '(' {
				depth--
			} else if depth == 0 {
				for _, op := range ops {
					if string(expr[i]) == op {
						if i > 0 {
							leftSide := strings.TrimSpace(expr[:i])
							if len(leftSide) > 0 {
								prevNonSpace := leftSide[len(leftSide)-1]
								if !isBinaryOperatorChar(prevNonSpace) {
									return i
								}
							}
						}
					}
				}
			}
		}
	}
	return -1
}

func isBinaryOperatorChar(c byte) bool {
	return c == '+' || c == '-' || c == '*' || c == '/' || c == '%' ||
		c == '<' || c == '>' || c == '=' || c == '!' || c == '&' || c == '|'
}

func isValidIdentifier(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i, r := range s {
		if i == 0 {
			if !unicode.IsLetter(r) && r != '_' {
				return false
			}
		} else {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
				return false
			}
		}
	}
	return true
}

func compareValues(left, right Value, op string) (Value, error) {
	if left.Type != right.Type {
		return Value{}, fmt.Errorf("cannot compare %s with %s", left.Type, right.Type)
	}

	var result bool
	switch left.Type {
	case PropTypeInt:
		switch op {
		case "==":
			result = left.IntVal == right.IntVal
		case "!=":
			result = left.IntVal != right.IntVal
		case "<":
			result = left.IntVal < right.IntVal
		case ">":
			result = left.IntVal > right.IntVal
		case "<=":
			result = left.IntVal <= right.IntVal
		case ">=":
			result = left.IntVal >= right.IntVal
		}
	case PropTypeString:
		switch op {
		case "==":
			result = left.StrVal == right.StrVal
		case "!=":
			result = left.StrVal != right.StrVal
		case "<":
			result = left.StrVal < right.StrVal
		case ">":
			result = left.StrVal > right.StrVal
		case "<=":
			result = left.StrVal <= right.StrVal
		case ">=":
			result = left.StrVal >= right.StrVal
		}
	case PropTypeBoolean:
		switch op {
		case "==":
			result = left.BoolVal == right.BoolVal
		case "!=":
			result = left.BoolVal != right.BoolVal
		default:
			return Value{}, fmt.Errorf("boolean comparison only supports == and !=")
		}
	}
	return Value{Type: PropTypeBoolean, BoolVal: result}, nil
}

func parseComponentAttributes(attrStr string, scopeProps map[string]Value, propDefs map[string]PropDef) (map[string]Value, error) {
	result := make(map[string]Value)
	i := 0
	for i < len(attrStr) {
		for i < len(attrStr) && unicode.IsSpace(rune(attrStr[i])) {
			i++
		}
		if i >= len(attrStr) {
			break
		}

		nameStart := i
		for i < len(attrStr) && attrStr[i] != '=' && !unicode.IsSpace(rune(attrStr[i])) {
			i++
		}
		if nameStart == i {
			break
		}
		name := attrStr[nameStart:i]

		for i < len(attrStr) && unicode.IsSpace(rune(attrStr[i])) {
			i++
		}
		if i >= len(attrStr) || attrStr[i] != '=' {
			continue
		}
		i++

		for i < len(attrStr) && unicode.IsSpace(rune(attrStr[i])) {
			i++
		}
		if i >= len(attrStr) {
			break
		}

		var value Value
		var err error

		if attrStr[i] == '{' {
			exprStart := i + 1
			depth := 1
			i++
			for i < len(attrStr) && depth > 0 {
				if attrStr[i] == '{' {
					depth++
				} else if attrStr[i] == '}' {
					depth--
				}
				i++
			}
			expr := attrStr[exprStart : i-1]
			value, err = EvaluateExpression(expr, scopeProps)
			if err != nil {
				return nil, fmt.Errorf("error evaluating expression for '%s': %v", name, err)
			}
		} else if attrStr[i] == '\'' || attrStr[i] == '"' {
			quote := attrStr[i]
			i++
			valueStart := i
			for i < len(attrStr) && attrStr[i] != quote {
				i++
			}
			strVal := attrStr[valueStart:i]
			i++

			if def, ok := propDefs[name]; ok {
				if def.Type != PropTypeString {
					return nil, fmt.Errorf("prop '%s' expects type '%s', but got raw string value. Use {expression} syntax for non-string types", name, def.Type)
				}
			}
			value = Value{Type: PropTypeString, StrVal: strVal}
		} else {
			continue
		}

		if def, ok := propDefs[name]; ok {
			if value.Type != def.Type {
				return nil, fmt.Errorf("prop '%s' expects type '%s', but got '%s'", name, def.Type, value.Type)
			}
		}
		result[name] = value
	}
	return result, nil
}

func ParsePropsAttribute(template string) (map[string]PropDef, string, error) {
	propDefs := make(map[string]PropDef)
	match := rePropsAttr.FindStringSubmatch(template)
	if match == nil {
		return propDefs, template, nil
	}

	propsStr := match[1]
	template = rePropsAttr.ReplaceAllString(template, "")

	pairs := strings.Split(propsStr, ",")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		parts := strings.Fields(pair)
		if len(parts) != 2 {
			return nil, "", fmt.Errorf("invalid prop definition '%s': expected 'name type' format", pair)
		}
		name := parts[0]
		propType := parts[1]

		if propType != PropTypeString && propType != PropTypeInt && propType != PropTypeBoolean {
			return nil, "", fmt.Errorf("invalid prop type '%s' for prop '%s': must be string, int, or boolean", propType, name)
		}

		if _, exists := propDefs[name]; exists {
			return nil, "", fmt.Errorf("duplicate prop name: %s", name)
		}
		propDefs[name] = PropDef{Name: name, Type: propType}
	}
	return propDefs, template, nil
}

func ParseAttributes(attrStr string) map[string]string {
	reAttrs := regexp.MustCompile(`(\w+)=["']([^"']*)["']`)
	matches := reAttrs.FindAllStringSubmatch(attrStr, -1)
	res := make(map[string]string)
	for _, m := range matches {
		res[m[1]] = m[2]
	}
	return res
}

func findFirstComponent(html string) (int, int, string, bool, string, string) {
	loc := reComponentTag.FindStringIndex(html)
	if loc == nil {
		return -1, -1, "", false, "", ""
	}

	start := loc[0]
	tagOpenContent := html[loc[0]:loc[1]]

	if strings.HasPrefix(tagOpenContent, "</") {
		s, e, t, sc, a, i := findFirstComponent(html[loc[1]:])
		if s != -1 {
			return loc[1] + s, loc[1] + e, t, sc, a, i
		}
		return -1, -1, "", false, "", ""
	}

	tagName := tagOpenContent[1:]

	rest := html[loc[1]:]
	closeBracket := strings.IndexAny(rest, ">")
	if closeBracket == -1 {
		return -1, -1, "", false, "", ""
	}
	attrEnd := loc[1] + closeBracket

	fullTagContent := html[start : attrEnd+1]
	attrsStr := html[loc[1]:attrEnd]

	if strings.HasSuffix(fullTagContent, "/>") {
		return start, attrEnd + 1, tagName, true, strings.TrimSuffix(attrsStr, "/"), ""
	}

	nestLevel := 1
	searchStart := attrEnd + 1
	closingTag := "</" + tagName + ">"
	openingTagPrefix := "<" + tagName

	for {
		nextClose := strings.Index(html[searchStart:], closingTag)
		if nextClose == -1 {
			return -1, -1, "", false, "", ""
		}
		absClose := searchStart + nextClose

		chunk := html[searchStart:absClose]
		nestedOpens := 0
		openIdx := 0
		for {
			idx := strings.Index(chunk[openIdx:], openingTagPrefix)
			if idx == -1 {
				break
			}
			charAfter := chunk[openIdx+idx+len(openingTagPrefix)]
			if charAfter == ' ' || charAfter == '>' || charAfter == '/' {
				nestedOpens++
			}
			openIdx += idx + 1
		}

		if nestedOpens == 0 {
			return start, absClose + len(closingTag), tagName, false, attrsStr, html[attrEnd+1 : absClose]
		}

		current := attrEnd + 1
		nestLevel = 1
		for nestLevel > 0 {
			nextOpenIdx := strings.Index(html[current:], openingTagPrefix)
			nextCloseIdx := strings.Index(html[current:], closingTag)

			if nextCloseIdx == -1 {
				return -1, -1, "", false, "", ""
			}

			if nextOpenIdx != -1 && nextOpenIdx < nextCloseIdx {
				nestLevel++
				current += nextOpenIdx + 1
			} else {
				nestLevel--
				if nestLevel == 0 {
					return start, current + nextCloseIdx + len(closingTag), tagName, false, attrsStr, html[attrEnd+1 : current+nextCloseIdx]
				}
				current += nextCloseIdx + len(closingTag)
			}
		}
	}
}

func HasSingleRoot(html string) bool {
	reComment := regexp.MustCompile(``)
	clean := reComment.ReplaceAllString(html, "")
	clean = strings.TrimSpace(clean)

	if !strings.HasPrefix(clean, "<") {
		return false
	}

	idx := strings.IndexAny(clean, " >/")
	if idx == -1 {
		return false
	}
	tagName := clean[1:idx]

	if strings.HasSuffix(clean, "/>") {
		return true
	}

	suffix := "</" + tagName + ">"
	if strings.HasSuffix(clean, suffix) {
		return true
	}
	return false
}

func ProcessComponentStyles(raw string, scopeID string) (string, string, error) {
	loc := reStyleBlock.FindStringSubmatchIndex(raw)
	if loc == nil {
		return raw, "", nil
	}

	cssContent := raw[loc[2]:loc[3]]
	htmlContent := raw[:loc[0]] + raw[loc[1]:]
	htmlContent = strings.TrimSpace(htmlContent)

	var scopedCSS strings.Builder
	blocks := strings.Split(cssContent, "}")
	for _, block := range blocks {
		if strings.TrimSpace(block) == "" {
			continue
		}
		parts := strings.Split(block, "{")
		if len(parts) != 2 {
			continue
		}
		selectors := parts[0]
		body := parts[1]

		selList := strings.Split(selectors, ",")
		var newSels []string
		for _, s := range selList {
			s = strings.TrimSpace(s)
			newSels = append(newSels, fmt.Sprintf("%s[%s]", s, scopeID))
		}
		scopedCSS.WriteString(strings.Join(newSels, ", ") + " {" + body + "}\n")
	}

	return htmlContent, scopedCSS.String(), nil
}

func InjectScopeID(html string, scopeID string) string {
	clean := strings.TrimSpace(html)
	firstSpace := strings.IndexAny(clean, " />")
	if firstSpace == -1 {
		return html
	}
	return clean[:firstSpace] + " " + scopeID + "=\"\"" + clean[firstSpace:]
}

func IsPascalCase(s string) bool {
	if len(s) == 0 {
		return false
	}
	r := []rune(s)
	if !unicode.IsUpper(r[0]) {
		return false
	}
	for _, x := range r {
		if !unicode.IsLetter(x) && !unicode.IsNumber(x) {
			return false
		}
	}
	return true
}

func IsKebabCase(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if !unicode.IsLower(r) && !unicode.IsNumber(r) && r != '-' {
			return false
		}
	}
	return true
}

type CompileOptions struct {
	ComponentsDir string
	RoutesDir     string
	DistDir       string
	StaticDir     string
}

func CompileProject(basePath string, opts CompileOptions) error {
	state := &GlobalState{
		Components: make(map[string]*Component),
	}

	compDir := filepath.Join(basePath, opts.ComponentsDir)
	if _, err := os.Stat(compDir); os.IsNotExist(err) {
		return fmt.Errorf("missing required directory: %s", compDir)
	}

	err := filepath.Walk(compDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".html" {
			return nil
		}

		name := strings.TrimSuffix(filepath.Base(path), ".html")
		if !IsPascalCase(name) {
			return fmt.Errorf("component '%s' must be PascalCase", path)
		}
		if _, exists := state.Components[name]; exists {
			return fmt.Errorf("duplicate component name found: %s", name)
		}

		contentBytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		content := string(contentBytes)

		scopeID := "data-" + strings.ToLower(name)
		template, css, err := ProcessComponentStyles(content, scopeID)
		if err != nil {
			return err
		}

		propDefs, template, err := ParsePropsAttribute(template)
		if err != nil {
			return fmt.Errorf("error parsing props in %s: %v", path, err)
		}

		if !HasSingleRoot(template) {
			return fmt.Errorf("component '%s' must have a single root element", name)
		}

		template = InjectScopeID(template, scopeID)

		state.Components[name] = &Component{
			Name:        name,
			RawContent:  content,
			Template:    template,
			ScopedStyle: css,
			ScopeID:     scopeID,
			Path:        path,
			PropDefs:    propDefs,
		}

		if css != "" {
			state.CSSOutput.WriteString("/* " + name + " */\n")
			state.CSSOutput.WriteString(css + "\n")
		}

		return nil
	})

	if err != nil {
		return err
	}

	routesDir := filepath.Join(basePath, opts.RoutesDir)
	if _, err := os.Stat(routesDir); os.IsNotExist(err) {
		return fmt.Errorf("missing required directory: %s", routesDir)
	}

	distDir := filepath.Join(basePath, opts.DistDir)

	err = filepath.Walk(routesDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".html" {
			return nil
		}

		relPath, _ := filepath.Rel(routesDir, path)
		fileName := strings.TrimSuffix(filepath.Base(path), ".html")
		if fileName != "index" && !IsKebabCase(fileName) {
			return fmt.Errorf("route '%s' must be kebab-case", path)
		}

		contentBytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		compiledHTML, err := CompileHTML(string(contentBytes), state, map[string]Value{})
		if err != nil {
			return fmt.Errorf("error compiling %s: %v", path, err)
		}

		outPath := filepath.Join(distDir, relPath)
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return err
		}

		return os.WriteFile(outPath, []byte(compiledHTML), 0644)
	})

	if err != nil {
		return err
	}

	staticDistDir := filepath.Join(distDir, opts.StaticDir)
	if err := os.MkdirAll(staticDistDir, 0755); err != nil {
		return err
	}

	cssFile := filepath.Join(staticDistDir, "styles.css")
	if err := os.WriteFile(cssFile, []byte(state.CSSOutput.String()), 0644); err != nil {
		return err
	}

	srcStatic := filepath.Join(basePath, opts.StaticDir)
	if _, err := os.Stat(srcStatic); err == nil {
		copyDir(srcStatic, staticDistDir)
	}

	return nil
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, _ := filepath.Rel(src, path)
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(dstPath, data, info.Mode())
	})
}

func WatchProject(basePath string, opts CompileOptions) {
	fmt.Printf("Watching %s for changes...\n", basePath)
	if err := CompileProject(basePath, opts); err != nil {
		fmt.Printf("Compile Error: %v\n", err)
	} else {
		fmt.Println("Built successfully.")
	}

	lastMod := time.Now()

	for {
		time.Sleep(1 * time.Second)
		needsCompile := false

		filepath.Walk(basePath, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				if strings.Contains(path, opts.DistDir) {
					return filepath.SkipDir
				}
				return nil
			}
			if info.ModTime().After(lastMod) {
				needsCompile = true
				return errors.New("change found")
			}
			return nil
		})

		if needsCompile {
			fmt.Println("Change detected. Compiling...")
			if err := CompileProject(basePath, opts); err != nil {
				fmt.Printf("Compile Error: %v\n", err)
			} else {
				fmt.Println("Built successfully.")
			}
			lastMod = time.Now()
		}
	}
}
