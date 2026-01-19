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

	// Fetch-related regex patterns
	reFetchAttr    = regexp.MustCompile(`\s+fetch\s*=\s*['"]([^'"]+)['"]`)
	reAsAttr       = regexp.MustCompile(`\s+as\s*=\s*['"]([^'"]+)['"]`)
	reForAttr      = regexp.MustCompile(`\s+for\s*=\s*['"]([^'"]+)['"]`)
	reSuspenseAttr = regexp.MustCompile(`\s+suspense(\s|>|/)`)
	reFallbackAttr = regexp.MustCompile(`\s+fallback(\s|>|/)`)

	// Interactivity-related regex patterns
	reGtmlScript      = regexp.MustCompile(`(?s)<script\s+type\s*=\s*['"]gtml['"]\s*>(.*?)</script>`)
	reInlineGtmlEvent = regexp.MustCompile(`(?s)\s(on[a-z]+)=\{\(\)\s*=>\s*\{([\s\S]*?)\}\}`)
	reSignalAccess    = regexp.MustCompile(`\$([a-zA-Z_][a-zA-Z0-9_]*)`)
	reSignalSet       = regexp.MustCompile(`\$([a-zA-Z_][a-zA-Z0-9_]*)\s*=\s*(.+)`)
	reElementSel      = regexp.MustCompile(`#([a-zA-Z_-][a-zA-Z0-9_-]*)(\*?)`)
	reClassSel        = regexp.MustCompile(`\.([a-zA-Z_-][a-zA-Z0-9_-]*)(\*?)`)
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
	Components      map[string]*Component
	CSSOutput       strings.Builder
	InteractivityJS strings.Builder
}

type Value struct {
	Type    string
	StrVal  string
	IntVal  int
	BoolVal bool
}

// FetchElement represents an element with client-side fetch behavior
type FetchElement struct {
	ID           string // Unique ID for JavaScript targeting
	Method       string // HTTP method (GET, POST, etc.)
	URL          string // URL to fetch from
	AsName       string // Name for the data (from 'as' attribute)
	StartIdx     int    // Start position in HTML
	EndIdx       int    // End position in HTML
	TagName      string // The HTML tag name
	InnerContent string // Content inside the element
	FullElement  string // The complete element HTML
}

// ForLoop represents a for iteration expression
type ForLoop struct {
	ItemName   string // Name of each item (e.g., 'user')
	SourceName string // Name of the source data (e.g., 'users')
	SourcePath string // Full path for nested access (e.g., 'user.colors')
	TemplateID string // Unique ID for the template element
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

func CompileHTML(html string, state *GlobalState, scopeProps map[string]Value, isTopLevel bool) (string, error) {
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

		compiledChildren, err := CompileHTML(innerContent, state, scopeProps, false)
		if err != nil {
			return "", err
		}

		slotsMap := extractSlots(compiledChildren)

		renderedComp := compDef.Template

		// Extract prop signals from gtml script BEFORE processing
		propSignals := extractPropSignals(renderedComp)

		// Process gtml scripts and mark signal expressions BEFORE expression evaluation
		renderedComp, gtmlScript, err := ProcessGtmlScripts(renderedComp, props)
		if err != nil {
			return "", fmt.Errorf("error processing gtml scripts in %s: %v", tagName, err)
		}
		if gtmlScript != "" {
			state.InteractivityJS.WriteString(gtmlScript)
		}

		// Process inline gtml events
		renderedComp, inlineScript, err := ProcessInlineEvents(renderedComp, props)
		if err != nil {
			return "", fmt.Errorf("error processing inline events in %s: %v", tagName, err)
		}
		if inlineScript != "" {
			state.InteractivityJS.WriteString(inlineScript)
		}

		// Mark signal expressions in the template for runtime rendering
		// This is done BEFORE expression evaluation
		signalNames := extractSignalNames(gtmlScript)
		inlineSignalNames := extractSignalNames(inlineScript)
		for sigName := range inlineSignalNames {
			signalNames[sigName] = true
		}
		for sigName := range signalNames {
			markerRe := regexp.MustCompile(`\{` + sigName + `\}`)
			renderedComp = markerRe.ReplaceAllString(renderedComp, fmt.Sprintf("{gtml-signal-%s}", sigName))
		}

		// Protect fetch expressions before evaluation
		renderedComp = protectFetchExpressions(renderedComp)

		renderedComp, err = EvaluateExpressions(renderedComp, props)
		if err != nil {
			return "", fmt.Errorf("error evaluating expressions in %s: %v", tagName, err)
		}

		// Restore escaped braces in event handlers
		renderedComp = strings.ReplaceAll(renderedComp, "&#123;", "{")
		renderedComp = strings.ReplaceAll(renderedComp, "&#125;", "}")

		// Replace the signal markers with spans for runtime rendering
		for sigName := range signalNames {
			markerRe := regexp.MustCompile(`\{gtml-signal-` + sigName + `\}`)
			renderedComp = markerRe.ReplaceAllString(renderedComp, fmt.Sprintf("<span data-gtml-signal-value='%s'></span>", sigName))
		}

		renderedComp, err = EvaluateExpressions(renderedComp, props)
		if err != nil {
			return "", fmt.Errorf("error evaluating expressions in %s: %v", tagName, err)
		}

		// Restore fetch expressions after all evaluations are done
		renderedComp = restoreFetchExpressions(renderedComp)

		// Add data attributes for prop signals
		for propName := range propSignals {
			attrName := fmt.Sprintf("data-gtml-prop-%s", propName)
			if propValue, ok := props[propName]; ok {
				serializedValue := propValue.String()
				// Add the attribute to the first element of the component
				if strings.HasPrefix(strings.TrimSpace(renderedComp), "<") {
					firstSpace := strings.IndexAny(renderedComp, " >")
					if firstSpace != -1 {
						if renderedComp[firstSpace] == ' ' {
							// Has attributes, insert after first space
							insertPos := firstSpace + 1
							for insertPos < len(renderedComp) && renderedComp[insertPos] == ' ' {
								insertPos++
							}
							renderedComp = renderedComp[:insertPos] + fmt.Sprintf("%s='%s' ", attrName, serializedValue) + renderedComp[insertPos:]
						} else if renderedComp[firstSpace] == '>' {
							// No attributes, insert before >
							renderedComp = renderedComp[:firstSpace] + fmt.Sprintf(" %s='%s'", attrName, serializedValue) + renderedComp[firstSpace:]
						}
					}
				}
			}
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

		finalRendered, err := CompileHTML(renderedComp, state, props, false)
		if err != nil {
			return "", err
		}

		html = html[:startIdx] + finalRendered + html[endIdx:]
	}

	// Process client-side fetch elements BEFORE evaluating remaining expressions
	// This preserves expressions like {user.name} for client-side JavaScript
	html, err = ProcessFetchElements(html)
	if err != nil {
		return "", err
	}

	// Process inline gtml events at the top level
	html, inlineScript, err := ProcessInlineEvents(html, scopeProps)
	if err != nil {
		return "", err
	}
	if inlineScript != "" {
		state.InteractivityJS.WriteString(inlineScript)
	}

	html, err = EvaluateExpressions(html, scopeProps)
	if err != nil {
		return "", err
	}

	// Restore escaped braces in event handlers
	html = strings.ReplaceAll(html, "&#123;", "{")
	html = strings.ReplaceAll(html, "&#125;", "}")

	// Append interactivity scripts to the HTML only at the top level
	if isTopLevel && state.InteractivityJS.Len() > 0 {
		html = html + "\n<script>\n" + SignalLibrary + "\n</script>\n" + state.InteractivityJS.String()
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

		// Skip expressions inside script tags (they're meant for client-side JavaScript)
		if isInsideScriptTag(result, fullStart) {
			offset = fullEnd
			continue
		}

		// Skip expressions inside style tags (they're meant for CSS)
		if isInsideStyleTag(result, fullStart) {
			offset = fullEnd
			continue
		}

		// Skip expressions that are marked as signal placeholders
		if strings.HasPrefix(expr, "gtml-signal-") {
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

// isInsideScriptTag checks if the given position is inside a script tag
func isInsideScriptTag(html string, pos int) bool {
	// Look backwards for opening script tag
	lastScriptOpen := strings.LastIndex(html[:pos], "<script")
	if lastScriptOpen == -1 {
		return false
	}

	// Check if there's a closing script tag between the opening and position
	lastScriptClose := strings.LastIndex(html[:pos], "</script>")
	return lastScriptClose < lastScriptOpen
}

func isInsideStyleTag(html string, pos int) bool {
	lastStyleOpen := strings.LastIndex(html[:pos], "<style")
	if lastStyleOpen == -1 {
		return false
	}

	lastStyleClose := strings.LastIndex(html[:pos], "</style>")
	return lastStyleClose < lastStyleOpen
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

func isInsideEventHandler(html string, pos int) bool {
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
			// Skip if inside a script tag (ternary expressions in scripts are for client-side JS)
			if isInsideScriptTag(s, i) {
				continue
			}
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
			// Output both selector forms:
			// 1. [scope] .class - for descendants of the scoped root
			// 2. [scope].class - for when the class is on the scoped root itself
			newSels = append(newSels, fmt.Sprintf("[%s] %s", scopeID, s))
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

// fetchCounter is used to generate unique IDs for fetch elements
var fetchCounter int

// Marker used to protect fetch expressions from compile-time evaluation
const fetchExprMarker = "@@GTML_FETCH_EXPR@@"

// protectFetchExpressions escapes expressions inside fetch elements so they aren't evaluated at compile time
func protectFetchExpressions(html string) string {
	result := html
	offset := 0

	for {
		// Find next fetch element
		loc := reFetchAttr.FindStringIndex(result[offset:])
		if loc == nil {
			break
		}

		// Find the start of this element
		fetchAttrPos := offset + loc[0]
		elementStart := -1
		for i := fetchAttrPos; i >= 0; i-- {
			if result[i] == '<' && i+1 < len(result) && result[i+1] != '/' {
				elementStart = i
				break
			}
		}
		if elementStart == -1 {
			offset = fetchAttrPos + 1
			continue
		}

		// Find the element boundaries
		startIdx, endIdx, _, _, _, innerContent := findElementAt(result, elementStart)
		if startIdx == -1 {
			offset = fetchAttrPos + 1
			continue
		}

		// Protect expressions in inner content
		protectedInner := protectExpressionsInContent(innerContent)

		// Find where the inner content starts and ends
		openTagEnd := strings.Index(result[startIdx:], ">") + startIdx + 1

		// Find the tag name for the closing tag
		tagName := ""
		for i := startIdx + 1; i < len(result); i++ {
			if result[i] == ' ' || result[i] == '>' || result[i] == '/' {
				tagName = result[startIdx+1 : i]
				break
			}
		}
		closeTag := "</" + tagName + ">"
		closeTagStart := endIdx - len(closeTag)

		innerStart := openTagEnd
		innerEnd := closeTagStart

		if innerStart < innerEnd && innerStart > 0 && innerEnd <= len(result) {
			result = result[:innerStart] + protectedInner + result[innerEnd:]
			offset = innerStart + len(protectedInner) + len(closeTag)
		} else {
			offset = endIdx
		}
	}

	return result
}

// protectExpressionsInContent replaces curly braces with marker to prevent evaluation
func protectExpressionsInContent(content string) string {
	// Replace { with marker so the expression regex won't match
	result := strings.ReplaceAll(content, "{", fetchExprMarker+"OPEN"+fetchExprMarker)
	result = strings.ReplaceAll(result, "}", fetchExprMarker+"CLOSE"+fetchExprMarker)
	return result
}

// restoreFetchExpressions restores protected expressions
func restoreFetchExpressions(html string) string {
	result := strings.ReplaceAll(html, fetchExprMarker+"OPEN"+fetchExprMarker, "{")
	result = strings.ReplaceAll(result, fetchExprMarker+"CLOSE"+fetchExprMarker, "}")
	return result
}

// ProcessFetchElements processes HTML to find fetch elements and generate JavaScript
func ProcessFetchElements(html string) (string, error) {
	result := html
	fetchElements := findFetchElements(result)

	if len(fetchElements) == 0 {
		return result, nil
	}

	// Process each fetch element from end to start (to preserve indices)
	for i := len(fetchElements) - 1; i >= 0; i-- {
		fe := fetchElements[i]

		// Generate a unique ID for this fetch element
		fetchCounter++
		fe.ID = fmt.Sprintf("gtml-fetch-%d", fetchCounter)

		// Process the element and generate JavaScript
		processedElement, script, err := processSingleFetchElement(fe)
		if err != nil {
			return "", fmt.Errorf("error processing fetch element: %v", err)
		}

		// Replace the original element with the processed version and script
		result = result[:fe.StartIdx] + processedElement + script + result[fe.EndIdx:]
	}

	return result, nil
}

// findFetchElements finds all elements with the fetch attribute
func findFetchElements(html string) []FetchElement {
	var elements []FetchElement
	offset := 0

	for {
		// Find next element with fetch attribute
		searchArea := html[offset:]
		loc := reFetchAttr.FindStringIndex(searchArea)
		if loc == nil {
			break
		}

		// Find the start of this element (go back to find <)
		fetchAttrPos := offset + loc[0]
		elementStart := -1
		for i := fetchAttrPos; i >= 0; i-- {
			if html[i] == '<' && i+1 < len(html) && html[i+1] != '/' {
				elementStart = i
				break
			}
		}
		if elementStart == -1 {
			offset = fetchAttrPos + 1
			continue
		}

		// Parse the element
		startIdx, endIdx, tagName, isSelfClosing, attrsStr, innerContent := findElementAt(html, elementStart)
		if startIdx == -1 {
			offset = fetchAttrPos + 1
			continue
		}

		// Extract fetch attribute value
		fetchMatch := reFetchAttr.FindStringSubmatch(attrsStr)
		if fetchMatch == nil {
			offset = endIdx
			continue
		}
		fetchValue := fetchMatch[1]

		// Parse METHOD URL format
		parts := strings.SplitN(strings.TrimSpace(fetchValue), " ", 2)
		if len(parts) != 2 {
			offset = endIdx
			continue
		}
		method := strings.ToUpper(parts[0])
		url := parts[1]

		// Extract 'as' attribute if present
		asName := ""
		if asMatch := reAsAttr.FindStringSubmatch(attrsStr); asMatch != nil {
			asName = asMatch[1]
		}

		fe := FetchElement{
			Method:       method,
			URL:          url,
			AsName:       asName,
			StartIdx:     startIdx,
			EndIdx:       endIdx,
			TagName:      tagName,
			InnerContent: innerContent,
			FullElement:  html[startIdx:endIdx],
		}

		if isSelfClosing {
			fe.InnerContent = ""
		}

		elements = append(elements, fe)
		offset = endIdx
	}

	return elements
}

// findElementAt finds the element starting at the given position
func findElementAt(html string, start int) (int, int, string, bool, string, string) {
	if start >= len(html) || html[start] != '<' {
		return -1, -1, "", false, "", ""
	}

	// Find tag name
	tagEnd := start + 1
	for tagEnd < len(html) && html[tagEnd] != ' ' && html[tagEnd] != '>' && html[tagEnd] != '/' {
		tagEnd++
	}
	tagName := html[start+1 : tagEnd]

	// Find the end of opening tag
	closeBracket := strings.IndexAny(html[tagEnd:], ">")
	if closeBracket == -1 {
		return -1, -1, "", false, "", ""
	}
	attrEnd := tagEnd + closeBracket

	fullOpenTag := html[start : attrEnd+1]
	attrsStr := html[tagEnd:attrEnd]

	// Check for self-closing
	if strings.HasSuffix(fullOpenTag, "/>") {
		return start, attrEnd + 1, tagName, true, strings.TrimSuffix(attrsStr, "/"), ""
	}

	// Find matching closing tag
	closingTag := "</" + tagName + ">"
	nestLevel := 1
	searchStart := attrEnd + 1
	openingTagPrefix := "<" + tagName

	for nestLevel > 0 {
		nextCloseIdx := strings.Index(html[searchStart:], closingTag)
		if nextCloseIdx == -1 {
			return -1, -1, "", false, "", ""
		}

		// Count nested opens in this chunk
		chunk := html[searchStart : searchStart+nextCloseIdx]
		nestedOpens := countTagOccurrences(chunk, openingTagPrefix)

		if nestedOpens == 0 {
			// Found our closing tag
			innerContent := html[attrEnd+1 : searchStart+nextCloseIdx]
			return start, searchStart + nextCloseIdx + len(closingTag), tagName, false, attrsStr, innerContent
		}

		// Need to skip past nested elements
		for j := 0; j < nestedOpens; j++ {
			nextCloseIdx = strings.Index(html[searchStart:], closingTag)
			if nextCloseIdx == -1 {
				return -1, -1, "", false, "", ""
			}
			searchStart += nextCloseIdx + len(closingTag)
		}
	}

	return -1, -1, "", false, "", ""
}

// countTagOccurrences counts how many times an opening tag appears in a string
func countTagOccurrences(s, tagPrefix string) int {
	count := 0
	idx := 0
	for {
		loc := strings.Index(s[idx:], tagPrefix)
		if loc == -1 {
			break
		}
		// Check if followed by space, > or /
		checkPos := idx + loc + len(tagPrefix)
		if checkPos < len(s) {
			c := s[checkPos]
			if c == ' ' || c == '>' || c == '/' {
				count++
			}
		}
		idx += loc + 1
	}
	return count
}

// processSingleFetchElement processes a single fetch element and returns the modified HTML and script
func processSingleFetchElement(fe FetchElement) (string, string, error) {
	// Extract suspense, fallback, and regular content
	suspenseContent, fallbackContent, regularContent := extractFetchChildren(fe.InnerContent)

	// Process for loops in the regular content
	processedContent, forLoops := processForElements(regularContent)

	// Build the modified element HTML by removing fetch-related attributes from opening tag
	// Find where the opening tag ends
	openTagEnd := strings.Index(fe.FullElement, ">")
	if openTagEnd == -1 {
		return "", "", fmt.Errorf("invalid element: missing >")
	}

	// Get just the opening tag
	openTag := fe.FullElement[:openTagEnd+1]

	// Remove fetch and as attributes from the opening tag
	modifiedOpenTag := reFetchAttr.ReplaceAllString(openTag, "")
	modifiedOpenTag = reAsAttr.ReplaceAllString(modifiedOpenTag, "")

	// Add the unique ID to the tag
	// Find position to insert ID (after tag name)
	spacePos := strings.IndexAny(modifiedOpenTag, " >")
	if spacePos == -1 {
		spacePos = len(modifiedOpenTag) - 1
	}

	if modifiedOpenTag[spacePos] == '>' {
		// No existing attributes, insert ID before >
		modifiedOpenTag = modifiedOpenTag[:spacePos] + fmt.Sprintf(" id=\"%s\"", fe.ID) + modifiedOpenTag[spacePos:]
	} else {
		// Has attributes, insert ID after tag name
		modifiedOpenTag = modifiedOpenTag[:spacePos] + fmt.Sprintf(" id=\"%s\"", fe.ID) + modifiedOpenTag[spacePos:]
	}

	// Build the complete element with closing tag but empty content
	// (content will be filled by JavaScript)
	modifiedElement := modifiedOpenTag + "</" + fe.TagName + ">"

	// Generate the JavaScript
	script := generateFetchScript(fe, suspenseContent, fallbackContent, processedContent, forLoops)

	return modifiedElement, script, nil
}

// extractFetchChildren extracts suspense, fallback, and regular content from fetch element children
func extractFetchChildren(content string) (suspense, fallback, regular string) {
	// Find and extract suspense element
	suspenseStart := findElementWithAttr(content, "suspense")
	if suspenseStart != -1 {
		startIdx, endIdx, _, _, _, inner := findElementAt(content, suspenseStart)
		if startIdx != -1 {
			suspense = inner
			content = content[:startIdx] + content[endIdx:]
		}
	}

	// Find and extract fallback element
	fallbackStart := findElementWithAttr(content, "fallback")
	if fallbackStart != -1 {
		startIdx, endIdx, _, _, _, inner := findElementAt(content, fallbackStart)
		if startIdx != -1 {
			fallback = inner
			content = content[:startIdx] + content[endIdx:]
		}
	}

	regular = strings.TrimSpace(content)
	return
}

// findElementWithAttr finds the start position of an element with the given attribute
func findElementWithAttr(html string, attrName string) int {
	pattern := regexp.MustCompile(`<\w+[^>]*\s+` + attrName + `(\s|>|/)`)
	loc := pattern.FindStringIndex(html)
	if loc == nil {
		return -1
	}
	return loc[0]
}

// forLoopCounter is used to generate unique IDs for for loops
var forLoopCounter int

// processForElements finds and processes elements with for attributes
func processForElements(content string) (string, []ForLoop) {
	var forLoops []ForLoop
	result := content

	// Find all elements with for attribute
	offset := 0
	for {
		loc := reForAttr.FindStringIndex(result[offset:])
		if loc == nil {
			break
		}

		// Find start of element
		forAttrPos := offset + loc[0]
		elementStart := -1
		for i := forAttrPos; i >= 0; i-- {
			if result[i] == '<' && i+1 < len(result) && result[i+1] != '/' {
				elementStart = i
				break
			}
		}
		if elementStart == -1 {
			offset = forAttrPos + 1
			continue
		}

		// Parse the element
		startIdx, endIdx, tagName, _, attrsStr, innerContent := findElementAt(result, elementStart)
		if startIdx == -1 {
			offset = forAttrPos + 1
			continue
		}

		// Extract for attribute value
		forMatch := reForAttr.FindStringSubmatch(attrsStr)
		if forMatch == nil {
			offset = endIdx
			continue
		}
		forValue := forMatch[1]

		// Parse "item in items" format
		forLoop, err := ParseForAttribute(forValue)
		if err != nil {
			offset = endIdx
			continue
		}

		// Recursively process nested for loops in inner content
		processedInner, nestedLoops := processForElements(innerContent)

		// Assign unique ID to this for loop
		forLoopCounter++
		templateID := fmt.Sprintf("gtml-for-%d", forLoopCounter)
		forLoop.TemplateID = templateID

		forLoops = append(forLoops, forLoop)
		// Append nested loops
		forLoops = append(forLoops, nestedLoops...)

		// Remove for attribute and add template markers
		newAttrs := reForAttr.ReplaceAllString(attrsStr, "")
		newElement := fmt.Sprintf("<%s%s data-gtml-for=\"%s\" data-gtml-item=\"%s\" data-gtml-source=\"%s\" style=\"display:none\">%s</%s>",
			tagName, newAttrs, templateID, forLoop.ItemName, forLoop.SourcePath, processedInner, tagName)

		result = result[:startIdx] + newElement + result[endIdx:]
		offset = startIdx + len(newElement)
	}

	return result, forLoops
}

// ParseForAttribute parses a for attribute value like "user in users" or "color in user.colors"
func ParseForAttribute(value string) (ForLoop, error) {
	parts := strings.Split(strings.TrimSpace(value), " in ")
	if len(parts) != 2 {
		return ForLoop{}, fmt.Errorf("invalid for attribute format: expected 'item in items', got '%s'", value)
	}

	itemName := strings.TrimSpace(parts[0])
	sourcePath := strings.TrimSpace(parts[1])

	// Get the base source name (first part before any dots)
	sourceName := sourcePath
	if idx := strings.Index(sourcePath, "."); idx != -1 {
		sourceName = sourcePath[:idx]
	}

	return ForLoop{
		ItemName:   itemName,
		SourceName: sourceName,
		SourcePath: sourcePath,
	}, nil
}

// generateFetchScript generates the JavaScript code for a fetch element
func generateFetchScript(fe FetchElement, suspenseContent, fallbackContent, regularContent string, forLoops []ForLoop) string {
	var script strings.Builder
	script.WriteString("\n<script>\n(function() {\n")

	// Get references to the container
	script.WriteString(fmt.Sprintf("  const container = document.getElementById('%s');\n", fe.ID))
	script.WriteString("  if (!container) return;\n\n")

	// Create suspense element if needed
	if suspenseContent != "" {
		script.WriteString("  // Create and show suspense element\n")
		script.WriteString("  const suspenseEl = document.createElement('div');\n")
		script.WriteString(fmt.Sprintf("  suspenseEl.innerHTML = `%s`;\n", escapeJSTemplate(suspenseContent)))
		script.WriteString("  suspenseEl.setAttribute('data-gtml-suspense', '');\n")
		script.WriteString("  container.appendChild(suspenseEl);\n\n")
	}

	// Create fallback element (hidden initially) if needed
	if fallbackContent != "" {
		script.WriteString("  // Create fallback element (hidden initially)\n")
		script.WriteString("  const fallbackEl = document.createElement('div');\n")
		script.WriteString(fmt.Sprintf("  fallbackEl.innerHTML = `%s`;\n", escapeJSTemplate(fallbackContent)))
		script.WriteString("  fallbackEl.setAttribute('data-gtml-fallback', '');\n")
		script.WriteString("  fallbackEl.style.display = 'none';\n")
		script.WriteString("  container.appendChild(fallbackEl);\n\n")
	}

	// Store the template content for iteration
	script.WriteString("  // Store template content\n")
	script.WriteString(fmt.Sprintf("  const templateContent = `%s`;\n\n", escapeJSTemplate(regularContent)))

	// Perform the fetch
	script.WriteString(fmt.Sprintf("  fetch('%s', { method: '%s' })\n", fe.URL, fe.Method))
	script.WriteString("    .then(response => {\n")
	script.WriteString("      if (!response.ok) throw new Error('Request failed');\n")
	script.WriteString("      return response.json();\n")
	script.WriteString("    })\n")
	script.WriteString(fmt.Sprintf("    .then(%s => {\n", fe.AsName))

	// Hide suspense
	if suspenseContent != "" {
		script.WriteString("      // Hide suspense\n")
		script.WriteString("      const suspense = container.querySelector('[data-gtml-suspense]');\n")
		script.WriteString("      if (suspense) suspense.remove();\n\n")
	}

	// Process for loops and render content
	if len(forLoops) > 0 {
		script.WriteString("      // Process iteration\n")
		script.WriteString("      const contentDiv = document.createElement('div');\n")
		script.WriteString("      contentDiv.innerHTML = templateContent;\n\n")

		// Add processForLoops function to handle nested iterations
		script.WriteString("      // Process all for loops recursively\n")
		script.WriteString("      function processForLoops(element, scope) {\n")
		script.WriteString("        const forElements = element.querySelectorAll('[data-gtml-for]');\n")
		script.WriteString("        forElements.forEach(template => {\n")
		script.WriteString("          // Skip if already processed (no longer has the attribute)\n")
		script.WriteString("          if (!template.hasAttribute('data-gtml-for')) return;\n")
		script.WriteString("          const itemName = template.getAttribute('data-gtml-item');\n")
		script.WriteString("          const sourcePath = template.getAttribute('data-gtml-source');\n")
		script.WriteString("          // Get source data from scope using path\n")
		script.WriteString("          const source = getValueByPath(scope, sourcePath);\n")
		script.WriteString("          if (!Array.isArray(source)) {\n")
		script.WriteString("            template.remove();\n")
		script.WriteString("            return;\n")
		script.WriteString("          }\n")
		script.WriteString("          const parent = template.parentNode;\n")
		script.WriteString("          source.forEach(item => {\n")
		script.WriteString("            const clone = template.cloneNode(true);\n")
		script.WriteString("            clone.removeAttribute('data-gtml-for');\n")
		script.WriteString("            clone.removeAttribute('data-gtml-item');\n")
		script.WriteString("            clone.removeAttribute('data-gtml-source');\n")
		script.WriteString("            clone.style.display = '';\n")
		script.WriteString("            // Create new scope with current item\n")
		script.WriteString("            const newScope = Object.assign({}, scope);\n")
		script.WriteString("            newScope[itemName] = item;\n")
		script.WriteString("            // Replace expressions in text nodes and attributes\n")
		script.WriteString("            clone.innerHTML = replaceExpressions(clone.innerHTML, newScope);\n")
		script.WriteString("            // Recursively process nested for loops\n")
		script.WriteString("            processForLoops(clone, newScope);\n")
		script.WriteString("            parent.insertBefore(clone, template);\n")
		script.WriteString("          });\n")
		script.WriteString("          template.remove();\n")
		script.WriteString("        });\n")
		script.WriteString("      }\n\n")

		// Initial scope with the fetched data
		script.WriteString(fmt.Sprintf("      const initialScope = { '%s': %s };\n", fe.AsName, fe.AsName))
		script.WriteString("      processForLoops(contentDiv, initialScope);\n\n")

		script.WriteString("      container.innerHTML = contentDiv.innerHTML;\n")
	} else {
		script.WriteString("      // Render content directly\n")
		script.WriteString("      container.innerHTML = templateContent;\n")
	}

	script.WriteString("    })\n")
	script.WriteString("    .catch(error => {\n")
	script.WriteString("      console.error('Fetch error:', error);\n")

	// Show fallback on error
	if suspenseContent != "" {
		script.WriteString("      // Hide suspense\n")
		script.WriteString("      const suspense = container.querySelector('[data-gtml-suspense]');\n")
		script.WriteString("      if (suspense) suspense.remove();\n")
	}
	if fallbackContent != "" {
		script.WriteString("      // Show fallback\n")
		script.WriteString("      const fallback = container.querySelector('[data-gtml-fallback]');\n")
		script.WriteString("      if (fallback) fallback.style.display = '';\n")
	}

	script.WriteString("    });\n\n")

	// Add helper function for getting value by path
	script.WriteString("  // Get value from object by dot-notation path\n")
	script.WriteString("  function getValueByPath(obj, path) {\n")
	script.WriteString("    const parts = path.split('.');\n")
	script.WriteString("    let value = obj[parts[0]];\n")
	script.WriteString("    for (let i = 1; i < parts.length && value !== undefined; i++) {\n")
	script.WriteString("      value = value[parts[i]];\n")
	script.WriteString("    }\n")
	script.WriteString("    return value;\n")
	script.WriteString("  }\n\n")

	// Add helper function for expression replacement
	script.WriteString("  // Helper function to replace expressions like {user.name} with actual values\n")
	script.WriteString("  function replaceExpressions(html, scope) {\n")
	script.WriteString("    return html.replace(/\\{([^}]+)\\}/g, (match, expr) => {\n")
	script.WriteString("      expr = expr.trim();\n")
	script.WriteString("      // Try to resolve the expression from scope\n")
	script.WriteString("      const parts = expr.split('.');\n")
	script.WriteString("      let value = scope[parts[0]];\n")
	script.WriteString("      for (let i = 1; i < parts.length && value !== undefined; i++) {\n")
	script.WriteString("        value = value[parts[i]];\n")
	script.WriteString("      }\n")
	script.WriteString("      return value !== undefined ? value : match;\n")
	script.WriteString("    });\n")
	script.WriteString("  }\n")

	script.WriteString("})();\n</script>\n")

	return script.String()
}

// escapeJSTemplate escapes a string for use in JavaScript template literals
func escapeJSTemplate(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "`", "\\`")
	s = strings.ReplaceAll(s, "${", "\\${")
	return s
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

		compiledHTML, err := CompileHTML(string(contentBytes), state, map[string]Value{}, true)
		if err != nil {
			return fmt.Errorf("error compiling %s: %v", path, err)
		}

		// Inject inline CSS into the head for reliable styling
		cssContent := state.CSSOutput.String()
		if cssContent != "" && strings.Contains(compiledHTML, "</head>") {
			inlineStyle := fmt.Sprintf("<style>\n%s</style>\n</head>", cssContent)
			compiledHTML = strings.Replace(compiledHTML, "</head>", inlineStyle, 1)
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

	// Copy static files first
	srcStatic := filepath.Join(basePath, opts.StaticDir)
	if _, err := os.Stat(srcStatic); err == nil {
		copyDir(srcStatic, staticDistDir)
	}

	// Write generated CSS last so it overwrites any placeholder from source static
	cssFile := filepath.Join(staticDistDir, "styles.css")
	if err := os.WriteFile(cssFile, []byte(state.CSSOutput.String()), 0644); err != nil {
		return err
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

const SignalLibrary = `// GTML Signal Library
class GtmlSignal {
  constructor(value) {
    this._value = value;
    this._subscribers = [];
  }

  get value() {
    return this._value;
  }

  set value(newValue) {
    const oldValue = this._value;
    this._value = newValue;
    this._subscribers.forEach(callback => callback(newValue, oldValue));
  }

  subscribe(callback) {
    this._subscribers.push(callback);
    return () => {
      this._subscribers = this._subscribers.filter(cb => cb !== callback);
    };
  }

  update(fn) {
    this.value = fn(this._value);
  }
}

function createSignal(initialValue) {
  return new GtmlSignal(initialValue);
}

const _gtmlSignalStore = new Map();

function getSignal(name) {
  if (!_gtmlSignalStore.has(name)) {
    _gtmlSignalStore.set(name, new GtmlSignal(null));
  }
  return _gtmlSignalStore.get(name);
}

function setSignal(name, value) {
  const signal = getSignal(name);
  signal.value = value;
}

function initSignal(name, value) {
  if (!_gtmlSignalStore.has(name)) {
    _gtmlSignalStore.set(name, new GtmlSignal(value));
  }
}

function _gtmlPropValue(propName) {
  const el = document.currentScript.previousElementSibling;
  if (el && el.hasAttribute('data-gtml-prop-' + propName)) {
    const value = el.getAttribute('data-gtml-prop-' + propName);
    if (value === 'null') return null;
    if (value === 'undefined') return undefined;
    if (value === 'true') return true;
    if (value === 'false') return false;
    if (!isNaN(value) && value !== '') return Number(value);
    return value;
  }
  return null;
}

function _gtmlRenderSignalValues() {
  document.querySelectorAll('[data-gtml-signal-value]').forEach(el => {
    const name = el.getAttribute('data-gtml-signal-value');
    const signal = getSignal(name);
    const unsubscribe = signal.subscribe((newVal) => {
      el.textContent = newVal;
    });
    el.textContent = signal.value;
  });
}
`

func ProcessInlineEvents(html string, props map[string]Value) (string, string, error) {
	matches := reInlineGtmlEvent.FindAllStringSubmatchIndex(html, -1)

	if len(matches) == 0 {
		return html, "", nil
	}

	signalRegistrations := make(map[string]bool)
	propSignals := make(map[string]bool)
	declaredSignals := make(map[string]bool)

	for _, match := range matches {
		eventName := html[match[2]:match[3]]
		eventStart := match[0]
		eventEnd := match[1]

		gtmlCode := strings.TrimSpace(html[match[4]:match[5]])

		declaredInScript := extractDeclaredSignals(gtmlCode)
		for _, sigName := range declaredInScript {
			declaredSignals[sigName] = true
			propSignals[sigName] = false
		}

		compiledScript, signals := CompileGtmlScript(gtmlCode)

		for sigName := range signals {
			if !signalRegistrations[sigName] {
				signalRegistrations[sigName] = true
			}
		}

		currentPropSignals := extractPropSignals(gtmlCode)
		for sigName := range currentPropSignals {
			if !propSignals[sigName] {
				propSignals[sigName] = true
			}
		}

		eventHandler := fmt.Sprintf(" %s=\"function() { %s }\"", eventName, compiledScript)
		// Escape braces to prevent expression regex from matching
		eventHandler = strings.ReplaceAll(eventHandler, "{", "&#123;")
		eventHandler = strings.ReplaceAll(eventHandler, "}", "&#125;")
		html = html[:eventStart] + eventHandler + html[eventEnd:]
	}

	signalInitCode := ""
	for sigName := range signalRegistrations {
		if propSignals[sigName] {
			if propVal, ok := props[sigName]; ok {
				var jsValue string
				switch propVal.Type {
				case PropTypeString:
					jsValue = fmt.Sprintf("'%s'", propVal.StrVal)
				case PropTypeInt:
					jsValue = fmt.Sprintf("%d", propVal.IntVal)
				case PropTypeBoolean:
					if propVal.BoolVal {
						jsValue = "true"
					} else {
						jsValue = "false"
					}
				default:
					jsValue = "null"
				}
				signalInitCode += fmt.Sprintf("  initSignal('%s', %s);\n", sigName, jsValue)
			} else {
				signalInitCode += fmt.Sprintf("  initSignal('%s', null);\n", sigName)
			}
		} else {
			signalInitCode += fmt.Sprintf("  initSignal('%s', null);\n", sigName)
		}
	}

	return html, signalInitCode, nil
}

func ProcessGtmlScripts(html string, props map[string]Value) (string, string, error) {
	matches := reGtmlScript.FindAllStringSubmatchIndex(html, -1)

	if len(matches) == 0 {
		return html, "", nil
	}

	signalRegistrations := make(map[string]bool)
	var compiledScripts strings.Builder
	var declaredSignals []string
	propSignals := make(map[string]bool)

	for _, match := range matches {
		scriptStart := match[2]
		scriptEnd := match[3]
		gtmlCode := html[scriptStart:scriptEnd]

		declaredSignals = extractDeclaredSignals(gtmlCode)
		compiledScript, signals := CompileGtmlScript(gtmlCode)

		for sigName := range signals {
			if !signalRegistrations[sigName] {
				signalRegistrations[sigName] = true
			}
		}

		for _, sigName := range declaredSignals {
			propSignals[sigName] = false
		}

		currentPropSignals := extractPropSignals(gtmlCode)
		for sigName := range currentPropSignals {
			if !propSignals[sigName] {
				propSignals[sigName] = true
			}
		}

		compiledScripts.WriteString(compiledScript)
	}

	signalInitCode := ""
	for sigName := range signalRegistrations {
		if propSignals[sigName] {
			// Use actual prop value if available, otherwise null
			if propVal, ok := props[sigName]; ok {
				// Format value appropriately for JavaScript
				var jsValue string
				switch propVal.Type {
				case PropTypeString:
					jsValue = fmt.Sprintf("'%s'", propVal.StrVal)
				case PropTypeInt:
					jsValue = fmt.Sprintf("%d", propVal.IntVal)
				case PropTypeBoolean:
					if propVal.BoolVal {
						jsValue = "true"
					} else {
						jsValue = "false"
					}
				default:
					jsValue = "null"
				}
				signalInitCode += fmt.Sprintf("  initSignal('%s', %s);\n", sigName, jsValue)
			} else {
				signalInitCode += fmt.Sprintf("  initSignal('%s', null);\n", sigName)
			}
		} else {
			signalInitCode += fmt.Sprintf("  initSignal('%s', null);\n", sigName)
		}
	}

	if signalInitCode != "" {
		signalInitCode = "\n" + signalInitCode
	}

	headerScript := fmt.Sprintf("\n<script>\n%s\n(function() {%s\n  _gtmlRenderSignalValues();\n})();\n</script>\n",
		signalInitCode, compiledScripts.String())

	result := reGtmlScript.ReplaceAllString(html, "")

	for sigName := range signalRegistrations {
		signalRe := regexp.MustCompile(`\{` + sigName + `\}`)
		result = signalRe.ReplaceAllString(result, fmt.Sprintf("<span data-gtml-signal-value='%s'></span>", sigName))
	}

	return result, headerScript, nil
}

func CompileGtmlScript(gtmlCode string) (string, map[string]bool) {
	signals := make(map[string]bool)

	code := gtmlCode

	code = convertElementSelectors(code, signals)
	code = convertEventBindings(code)
	code = convertSignalOperations(code, signals)
	code = convertSignalAccess(code, signals)

	return strings.TrimSpace(code), signals
}

// convertEventBindings converts .onclick(function() {...}) to .onclick = function() {...};
func convertEventBindings(code string) string {
	// Pattern to match event bindings like .onclick(function() or .onclick(function ()
	reEventBinding := regexp.MustCompile(`\.(on[a-z]+)\(function\s*\(\s*\)\s*\{`)
	code = reEventBinding.ReplaceAllString(code, ".$1 = function() {")

	// Convert closing }); to }; for event bindings (with semicolon)
	reClosingParenSemi := regexp.MustCompile(`\}\);`)
	code = reClosingParenSemi.ReplaceAllString(code, "};")

	// Convert closing }) to }; for event bindings (without semicolon)
	// This handles the closing parenthesis from the original .onclick(function(){}) syntax
	reClosingParen := regexp.MustCompile(`\}\)(\s*)$`)
	code = reClosingParen.ReplaceAllString(code, "};$1")

	// Handle }) followed by newlines and more code
	reClosingParenMidCode := regexp.MustCompile(`\}\)(\s*\n)`)
	code = reClosingParenMidCode.ReplaceAllString(code, "};$1")

	return code
}

func convertElementSelectors(code string, signals map[string]bool) string {
	reElementSel = regexp.MustCompile(`(^|[(\s;,])#([a-zA-Z_-][a-zA-Z0-9_-]*)(\*?)($|[)\s.,;])`)
	code = reElementSel.ReplaceAllStringFunc(code, func(match string) string {
		submatch := reElementSel.FindStringSubmatch(match)
		if len(submatch) < 5 {
			return match
		}
		prefix := submatch[1]
		id := submatch[2]
		isAll := submatch[3] == "*"
		suffix := submatch[4]
		selector := fmt.Sprintf("document.querySelector%s('#%s')", map[bool]string{true: "All", false: ""}[isAll], id)
		return prefix + selector + suffix
	})

	reClassSel = regexp.MustCompile(`(^|[(\s;,])\.([a-zA-Z_-][a-zA-Z0-9_-]*)(\*?)($|[)\s.,;])`)
	code = reClassSel.ReplaceAllStringFunc(code, func(match string) string {
		submatch := reClassSel.FindStringSubmatch(match)
		if len(submatch) < 5 {
			return match
		}
		prefix := submatch[1]
		class := submatch[2]
		isAll := submatch[3] == "*"
		suffix := submatch[4]
		selector := fmt.Sprintf("document.querySelector%s('.%s')", map[bool]string{true: "All", false: ""}[isAll], class)
		return prefix + selector + suffix
	})

	return code
}

func convertSignalOperations(code string, signals map[string]bool) string {
	code = reSignalSet.ReplaceAllStringFunc(code, func(match string) string {
		submatch := reSignalSet.FindStringSubmatch(match)
		if len(submatch) < 3 {
			return match
		}
		sigName := submatch[1]
		rightSide := strings.TrimSpace(submatch[2])
		signals[sigName] = true

		compiledRight := compileSignalExpression(rightSide, signals)

		return fmt.Sprintf("setSignal('%s', %s)", sigName, compiledRight)
	})

	return code
}

func compileSignalExpression(expr string, signals map[string]bool) string {
	processed := expr

	reCompoundOp := regexp.MustCompile(`\$([a-zA-Z_][a-zA-Z0-9_]*)\s*([+\-*/%]|==|!=|<=|>=|<|>)\s*(.+)`)
	processed = reCompoundOp.ReplaceAllStringFunc(processed, func(match string) string {
		submatch := reCompoundOp.FindStringSubmatch(match)
		if len(submatch) < 4 {
			return match
		}
		sigName := submatch[1]
		op := submatch[2]
		rightSide := strings.TrimSpace(submatch[3])
		signals[sigName] = true

		compiledRight := compileSignalExpression(rightSide, signals)

		return fmt.Sprintf("(getSignal('%s').value %s %s)", sigName, op, compiledRight)
	})

	processed = reSignalAccess.ReplaceAllStringFunc(processed, func(match string) string {
		submatch := reSignalAccess.FindStringSubmatch(match)
		if len(submatch) < 2 {
			return match
		}
		sigName := submatch[1]
		signals[sigName] = true
		return fmt.Sprintf("getSignal('%s').value", sigName)
	})

	return processed
}

func convertSignalAccess(code string, signals map[string]bool) string {
	code = reSignalAccess.ReplaceAllStringFunc(code, func(match string) string {
		submatch := reSignalAccess.FindStringSubmatch(match)
		if len(submatch) < 2 {
			return match
		}
		sigName := submatch[1]
		signals[sigName] = true
		return fmt.Sprintf("getSignal('%s').value", sigName)
	})
	return code
}

func extractSignalNames(script string) map[string]bool {
	signals := make(map[string]bool)
	reSignalUsage := regexp.MustCompile(`setSignal\('([a-zA-Z_][a-zA-Z0-9_]*)'`)
	matches := reSignalUsage.FindAllStringSubmatchIndex(script, -1)
	for _, match := range matches {
		if len(match) >= 4 {
			sigName := script[match[2]:match[3]]
			signals[sigName] = true
		}
	}
	return signals
}

func extractPropSignals(script string) map[string]bool {
	propSignals := make(map[string]bool)
	signals := extractSignalNames(script)
	reSignalAccess := regexp.MustCompile(`\$([a-zA-Z_][a-zA-Z0-9_]*)`)
	matches := reSignalAccess.FindAllStringSubmatch(script, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			sigName := match[1]
			if !signals[sigName] && sigName != "this" {
				propSignals[sigName] = true
			}
		}
	}
	return propSignals
}

func extractDeclaredSignals(script string) []string {
	var declared []string
	re := regexp.MustCompile(`\$([a-zA-Z_][a-zA-Z0-9_]*)\s*=`)
	matches := re.FindAllStringSubmatch(script, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			sigName := match[1]
			declared = append(declared, sigName)
		}
	}
	return declared
}
