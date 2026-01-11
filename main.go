package main

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

// Directory and File Constants
const (
	DirComponents   = "components"
	DirRoutes       = "routes"
	DirDist         = "dist"
	DirStatic       = "static"
	FileStyleCSS    = "styles.css"
	DirPreinstalled = "spec/components/preinstalled_components"
)

// Prop Type Constants
const (
	PropTypeString  = "string"
	PropTypeInt     = "int"
	PropTypeBoolean = "boolean"
)

// Regular Expressions
var (
	// Matches <Component ...>
	reComponentTag = regexp.MustCompile(`</?([A-Z][a-zA-Z0-9]*)`)
	// Matches <style>...</style>
	reStyleBlock = regexp.MustCompile(`(?s)<style>(.*?)</style>`)
	// Matches props='...'
	rePropsAttr = regexp.MustCompile(`\s+props\s*=\s*['"]([^'"]+)['"]`)
	// Matches {expression}
	reExpression = regexp.MustCompile(`\{([^{}]+)\}`)
	// Matches <slot name='...' /> or <slot name='...'>default</slot> in Component Definition
	reSlotPlaceholder = regexp.MustCompile(`(?s)<slot\s+name=['"](\w+)['"]\s*/?>`)
	// Matches usage of slots in children: <slot name='x' tag='y'>...</slot>
	reSlotUsage = regexp.MustCompile(`(?s)<slot\s+([^>]+)>(.*?)</slot>`)
	// Matches {{ slot: name }} in Component Definition
	reSlotDef = regexp.MustCompile(`\{\{\s*slot:\s*(\w+)\s*\}\}`)
)

// Data Structures

type PropDef struct {
	Name string
	Type string // "string", "int", "boolean"
}

type Component struct {
	Name        string
	RawContent  string // Original file content
	Template    string // HTML after stripping style and props attr
	Styles      string // Raw CSS content
	ScopedStyle string // Compiled/Scoped CSS
	ScopeID     string // Unique ID for scoping
	Path        string
	PropDefs    map[string]PropDef // Prop definitions from props attribute
}

type GlobalState struct {
	Components map[string]*Component
	CSSOutput  strings.Builder
}

type Value struct {
	Type    string // "string", "int", "boolean"
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

// CLI Entry Point

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		force := false
		path := ""
		for _, arg := range os.Args[2:] {
			if arg == "--force" {
				force = true
			} else if !strings.HasPrefix(arg, "-") && path == "" {
				path = arg
			}
		}
		if path == "" {
			fmt.Println("Error: Missing path argument for init.")
			fmt.Println("Usage: gtml init <PATH> [--force]")
			os.Exit(1)
		}
		runInit(path, force)

	case "compile":
		watch := false
		path := ""
		for _, arg := range os.Args[2:] {
			if arg == "--watch" {
				watch = true
			} else if !strings.HasPrefix(arg, "-") && path == "" {
				path = arg
			}
		}
		if path == "" {
			fmt.Println("Error: Missing path argument for compile.")
			fmt.Println("Usage: gtml compile <PATH> [--watch]")
			os.Exit(1)
		}

		if watch {
			runWatch(path)
		} else {
			if err := runCompile(path); err != nil {
				fmt.Printf("\n❌ Compilation failed: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("\n✅ Compilation successful!")
		}

	case "test":
		path := ""
		for _, arg := range os.Args[2:] {
			if !strings.HasPrefix(arg, "-") && path == "" {
				path = arg
			}
		}
		if path == "" {
			path = "."
		}
		runTests(path)

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("gtml - A Static Site Generator")
	fmt.Println("Usage:")
	fmt.Println("  gtml init <PATH> [--force]")
	fmt.Println("  gtml compile <PATH> [--watch]")
	fmt.Println("  gtml test [PATH]")
}

// Init Logic

func runInit(basePath string, force bool) {
	if _, err := os.Stat(basePath); err == nil && !force {
		fmt.Printf("Error: Directory '%s' already exists. Use --force to overwrite.\n", basePath)
		os.Exit(1)
	}

	if force {
		os.RemoveAll(basePath)
	}

	dirs := []string{
		filepath.Join(basePath, DirComponents),
		filepath.Join(basePath, DirRoutes),
		filepath.Join(basePath, DirStatic),
		filepath.Join(basePath, DirDist),
		filepath.Join(basePath, DirDist, DirStatic),
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			panic(fmt.Sprintf("Failed to create directory %s: %v", d, err))
		}
	}

	// Updated to match spec: uses <slot name='...' /> syntax
	files := map[string]string{
		filepath.Join(basePath, DirComponents, "BasicButton.html"): `<button props='text string'>{text}</button>`,
		filepath.Join(basePath, DirComponents, "GuestLayout.html"): `<html props='title string'>
  <head>
    <title>{title}</title>
  </head>
  <body>
    <BasicButton text={title} />
    <slot name='content' />
  </body>
</html>`,
		filepath.Join(basePath, DirRoutes, "index.html"): `<GuestLayout title="Some Title">
  <slot name='content' tag='div'>
    <p>Some Content</p>
  </slot>
</GuestLayout>`,
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			panic(fmt.Sprintf("Failed to write file %s: %v", path, err))
		}
	}

	// Copy preinstalled components
	if err := copyPreinstalledComponents(basePath); err != nil {
		fmt.Printf("Warning: Failed to copy preinstalled components: %v\n", err)
	}

	// Attempt initial compilation to ensure generated code works
	if err := runCompile(basePath); err != nil {
		fmt.Printf("Warning: Initial compilation failed: %v\n", err)
	} else {
		fmt.Printf("Initialized gtml project at %s\n", basePath)
	}
}

func copyPreinstalledComponents(basePath string) error {
	srcDir := DirPreinstalled
	dstDir := filepath.Join(basePath, DirComponents)

	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		return fmt.Errorf("preinstalled components directory not found: %s", srcDir)
	}

	return filepath.Walk(srcDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".html" {
			return nil
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dstDir, relPath)

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
			return err
		}

		return os.WriteFile(dstPath, content, 0644)
	})
}

// Compilation Logic

func runCompile(basePath string) error {
	state := &GlobalState{
		Components: make(map[string]*Component),
	}

	// 1. Parse Components
	compDir := filepath.Join(basePath, DirComponents)
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
		if !isPascalCase(name) {
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
		template, css, err := processComponentStyles(content, scopeID)
		if err != nil {
			return err
		}

		propDefs, template, err := parsePropsAttribute(template)
		if err != nil {
			return fmt.Errorf("error parsing props in %s: %v", path, err)
		}

		if !hasSingleRoot(template) {
			return fmt.Errorf("component '%s' must have a single root element", name)
		}

		template = injectScopeID(template, scopeID)

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

	// 2. Parse and Compile Routes
	routesDir := filepath.Join(basePath, DirRoutes)
	if _, err := os.Stat(routesDir); os.IsNotExist(err) {
		return fmt.Errorf("missing required directory: %s", routesDir)
	}

	distDir := filepath.Join(basePath, DirDist)

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
		// Check for valid kebab-case, allow index
		if fileName != "index" && !isKebabCase(fileName) {
			return fmt.Errorf("route '%s' must be kebab-case", path)
		}

		contentBytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		compiledHTML, err := compileHTML(string(contentBytes), state, map[string]Value{})
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

	// 3. Write CSS and Static Assets
	staticDistDir := filepath.Join(distDir, DirStatic)
	if err := os.MkdirAll(staticDistDir, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(staticDistDir, FileStyleCSS), []byte(state.CSSOutput.String()), 0644); err != nil {
		return err
	}

	srcStatic := filepath.Join(basePath, DirStatic)
	if _, err := os.Stat(srcStatic); err == nil {
		copyDir(srcStatic, staticDistDir)
	}

	return nil
}

// HTML Compilation Core

func compileHTML(html string, state *GlobalState, scopeProps map[string]Value) (string, error) {
	// 1. Evaluate Ternaries (Pre-processing for control flow)
	var err error
	html, err = evaluateTernaries(html, scopeProps)
	if err != nil {
		return "", err
	}

	// 2. Resolve Components
	for {
		startIdx, endIdx, tagName, _, attrsStr, innerContent := findFirstComponent(html)
		if startIdx == -1 {
			break // No components found
		}

		compDef, exists := state.Components[tagName]
		if !exists {
			return "", fmt.Errorf("component '%s' not found", tagName)
		}

		// Parse Props provided to this component instance
		props, err := parseComponentAttributes(attrsStr, scopeProps, compDef.PropDefs)
		if err != nil {
			return "", fmt.Errorf("error parsing attributes for %s: %v", tagName, err)
		}

		// Compile children (recursion) before extracting slots
		// We pass the *current* scopeProps to children because they exist in the current scope
		compiledChildren, err := compileHTML(innerContent, state, scopeProps)
		if err != nil {
			return "", err
		}

		// Map Slots provided in the child content
		slotsMap := extractSlots(compiledChildren)

		// Prepare component template
		renderedComp := compDef.Template

		// Convert {{ slot: name }} syntax to <slot name='name' /> for slot injection
		renderedComp = reSlotDef.ReplaceAllStringFunc(renderedComp, func(match string) string {
			subMatch := reSlotDef.FindStringSubmatch(match)
			if len(subMatch) < 2 {
				return match
			}
			slotName := subMatch[1]
			return "<slot name='" + slotName + "' />"
		})

		// Evaluate expressions within the component template using its OWN props
		renderedComp, err = evaluateExpressions(renderedComp, props)
		if err != nil {
			return "", fmt.Errorf("error evaluating expressions in %s: %v", tagName, err)
		}

		// Inject Slots into the Component Template
		renderedComp = reSlotPlaceholder.ReplaceAllStringFunc(renderedComp, func(match string) string {
			subMatch := reSlotPlaceholder.FindStringSubmatch(match)
			if len(subMatch) < 2 {
				return ""
			}
			slotName := subMatch[1]
			if content, ok := slotsMap[slotName]; ok {
				return content
			}
			// Spec doesn't explicitly detail default slot fallback in main.go context,
			// but if the user provides empty slot in template, we replace with empty.
			// To support default content inside <slot>Default</slot>, regex parsing would need to be deeper.
			// Assuming <slot name='x' /> usage for now based on prompt.
			return ""
		})

		// Recursively compile the result (in case the component template contained other components)
		// Note: We use `props` here because we are now inside the component context
		finalRendered, err := compileHTML(renderedComp, state, props)
		if err != nil {
			return "", err
		}

		// Replace the component tag in the original HTML with the rendered output
		html = html[:startIdx] + finalRendered + html[endIdx:]
	}

	// 3. Evaluate Expressions in the current scope (for pure HTML elements)
	html, err = evaluateExpressions(html, scopeProps)
	if err != nil {
		return "", err
	}

	return html, nil
}

// Slot Extraction Logic
func extractSlots(content string) map[string]string {
	slots := make(map[string]string)
	// Find all <slot name='...' tag='...'>...</slot> blocks in the content
	matches := reSlotUsage.FindAllStringSubmatchIndex(content, -1)

	// Process backwards to handle replacements or simple extraction
	// Since we are extracting to a map, order matters less for the map,
	// but strictly speaking, we just need to pull the data.
	for _, loc := range matches {
		// loc[0]: start, loc[1]: end
		// loc[2]: start of attrs, loc[3]: end of attrs
		// loc[4]: start of inner, loc[5]: end of inner

		attrsStr := content[loc[2]:loc[3]]
		innerContent := content[loc[4]:loc[5]]

		attrs := parseAttributes(attrsStr)
		name := attrs["name"]
		tag := attrs["tag"]

		if name != "" && tag != "" {
			// Construct the wrapper tag as per spec
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

// Expression Evaluation

func evaluateExpressions(html string, props map[string]Value) (string, error) {
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

		// Ignore double braces if they are part of a ternary or object logic not yet processed
		// or if we are inside a props="..." attribute (handled elsewhere)
		if isInsideComponentTag(result, fullStart) {
			offset = fullEnd
			continue
		}

		// Simple heuristics to skip things that look like ternaries (handled by evaluateTernaries)
		// This prevents { cond ? ( ... ) : ( ... ) } from being mangled by simple expression eval
		if strings.Contains(expr, "?") && strings.Contains(expr, "(") {
			offset = fullEnd
			continue
		}

		value, err := evaluateExpression(expr, props)
		if err != nil {
			return "", err
		}

		result = result[:fullStart] + value.String() + result[fullEnd:]
		// Adjust offset
		offset = fullStart + len(value.String())
	}
	return result, nil
}

func isInsideComponentTag(html string, pos int) bool {
	// Look backwards. If we find '<' followed by UpperCase before we find '>', we are inside.
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

// Ternary Logic

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

		condValue, err := evaluateExpression(condition, props)
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

		// Recursively evaluate ternaries inside the result
		replacement, err = evaluateTernaries(replacement, props)
		if err != nil {
			return "", err
		}

		result = result[:ternaryStart] + replacement + result[ternaryEnd:]
	}
	return result, nil
}

// Identifies the start of { ... ? ( ... ) : ( ... ) }
func findTernaryStart(s string) int {
	for i := 0; i < len(s); i++ {
		// Look for {
		if s[i] == '{' {
			// Avoid {{ logic if needed, but simple check:
			// Look ahead for ? and (
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
						return i // Found potential ternary start
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

// Parses: { condition ? (truthy) : (falsy) }
func parseTernary(s string, pos int) (int, string, string, string, error) {
	if s[pos] != '{' {
		return 0, "", "", "", fmt.Errorf("expected '{' at position %d", pos)
	}

	// 1. Find condition end (?)
	// We need to respect nested parens inside condition
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
			// Abort if we hit closing brace before question mark
			return 0, "", "", "", fmt.Errorf("malformed ternary, missing '?'")
		}
	}

	if conditionEnd == -1 {
		return 0, "", "", "", fmt.Errorf("missing '?' in ternary expression")
	}

	condition := strings.TrimSpace(s[pos+1 : conditionEnd])

	// 2. Find Truthy Block (...)
	// Scan for first '('
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

	// 3. Find Colon
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

	// 4. Find Falsy Block (...)
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

	// 5. Find Closing Brace }
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

// Expression Parser

func evaluateExpression(expr string, props map[string]Value) (Value, error) {
	expr = strings.TrimSpace(expr)

	// Boolean literals
	if expr == "true" {
		return Value{Type: PropTypeBoolean, BoolVal: true}, nil
	}
	if expr == "false" {
		return Value{Type: PropTypeBoolean, BoolVal: false}, nil
	}

	// String literals (single or double quoted)
	if (strings.HasPrefix(expr, "'") && strings.HasSuffix(expr, "'")) ||
		(strings.HasPrefix(expr, "\"") && strings.HasSuffix(expr, "\"")) {
		if len(expr) < 2 {
			return Value{Type: PropTypeString, StrVal: ""}, nil
		}
		return Value{Type: PropTypeString, StrVal: expr[1 : len(expr)-1]}, nil
	}

	// Integer literals
	if i, err := strconv.Atoi(expr); err == nil {
		return Value{Type: PropTypeInt, IntVal: i}, nil
	}
	// Negative Ints
	if strings.HasPrefix(expr, "-") {
		rest := strings.TrimSpace(expr[1:])
		if i, err := strconv.Atoi(rest); err == nil {
			return Value{Type: PropTypeInt, IntVal: -i}, nil
		}
	}

	// Logic Operators (||, &&)
	if idx := findOperator(expr, "||"); idx != -1 {
		left, err := evaluateExpression(expr[:idx], props)
		if err != nil {
			return Value{}, err
		}
		right, err := evaluateExpression(expr[idx+2:], props)
		if err != nil {
			return Value{}, err
		}
		if left.Type != PropTypeBoolean || right.Type != PropTypeBoolean {
			return Value{}, fmt.Errorf("|| operator requires boolean operands")
		}
		return Value{Type: PropTypeBoolean, BoolVal: left.BoolVal || right.BoolVal}, nil
	}
	if idx := findOperator(expr, "&&"); idx != -1 {
		left, err := evaluateExpression(expr[:idx], props)
		if err != nil {
			return Value{}, err
		}
		right, err := evaluateExpression(expr[idx+2:], props)
		if err != nil {
			return Value{}, err
		}
		if left.Type != PropTypeBoolean || right.Type != PropTypeBoolean {
			return Value{}, fmt.Errorf("&& operator requires boolean operands")
		}
		return Value{Type: PropTypeBoolean, BoolVal: left.BoolVal && right.BoolVal}, nil
	}

	// Comparison Operators
	// Order matters: match <=, >= before <, >
	for _, op := range []string{"==", "!=", "<=", ">=", "<", ">"} {
		if idx := findOperator(expr, op); idx != -1 {
			left, err := evaluateExpression(expr[:idx], props)
			if err != nil {
				return Value{}, err
			}
			right, err := evaluateExpression(expr[idx+len(op):], props)
			if err != nil {
				return Value{}, err
			}
			return compareValues(left, right, op)
		}
	}

	// Arithmetic (RTL for proper associativity in simple recursive parser)
	if idx := findOperatorRTL(expr, "+", "-"); idx != -1 {
		op := string(expr[idx])
		left, err := evaluateExpression(expr[:idx], props)
		if err != nil {
			return Value{}, err
		}
		right, err := evaluateExpression(expr[idx+1:], props)
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
		left, err := evaluateExpression(expr[:idx], props)
		if err != nil {
			return Value{}, err
		}
		right, err := evaluateExpression(expr[idx+1:], props)
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

	// Parentheses
	if strings.HasPrefix(expr, "(") && strings.HasSuffix(expr, ")") {
		return evaluateExpression(expr[1:len(expr)-1], props)
	}

	// Variable Lookup
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
		// Handle strings
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
						// Ensure it's not a unary operator (negative number check simplified)
						// If - is at start, it's unary. If preceded by operator, it's unary.
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
		// String comparison lexically
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

// Attribute Parsing

func parseComponentAttributes(attrStr string, scopeProps map[string]Value, propDefs map[string]PropDef) (map[string]Value, error) {
	result := make(map[string]Value)
	i := 0
	for i < len(attrStr) {
		// Skip whitespace
		for i < len(attrStr) && unicode.IsSpace(rune(attrStr[i])) {
			i++
		}
		if i >= len(attrStr) {
			break
		}

		// Name
		nameStart := i
		for i < len(attrStr) && attrStr[i] != '=' && !unicode.IsSpace(rune(attrStr[i])) {
			i++
		}
		if nameStart == i {
			break
		}
		name := attrStr[nameStart:i]

		// Skip to equals
		for i < len(attrStr) && unicode.IsSpace(rune(attrStr[i])) {
			i++
		}
		if i >= len(attrStr) || attrStr[i] != '=' {
			// Boolean prop shorthand? Not fully spec'd, assuming value required for now or boolean true
			// For this spec, let's assume attr="val" format is strictly required based on examples
			continue
		}
		i++ // skip =

		// Skip whitespace
		for i < len(attrStr) && unicode.IsSpace(rune(attrStr[i])) {
			i++
		}
		if i >= len(attrStr) {
			break
		}

		var value Value
		var err error

		// Determine if Expression {} or String ""
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
			value, err = evaluateExpression(expr, scopeProps)
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
			i++ // skip closing quote

			// Strict Type Check for raw strings
			if def, ok := propDefs[name]; ok {
				if def.Type != PropTypeString {
					return nil, fmt.Errorf("prop '%s' expects type '%s', but got raw string value. Use {expression} syntax for non-string types", name, def.Type)
				}
			}
			value = Value{Type: PropTypeString, StrVal: strVal}
		} else {
			continue
		}

		// Type Validation
		if def, ok := propDefs[name]; ok {
			if value.Type != def.Type {
				return nil, fmt.Errorf("prop '%s' expects type '%s', but got '%s'", name, def.Type, value.Type)
			}
		}
		result[name] = value
	}
	return result, nil
}

func parsePropsAttribute(template string) (map[string]PropDef, string, error) {
	propDefs := make(map[string]PropDef)
	match := rePropsAttr.FindStringSubmatch(template)
	if match == nil {
		return propDefs, template, nil // No props defined
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

func parseAttributes(attrStr string) map[string]string {
	reAttrs := regexp.MustCompile(`(\w+)=["']([^"']*)["']`)
	matches := reAttrs.FindAllStringSubmatch(attrStr, -1)
	res := make(map[string]string)
	for _, m := range matches {
		res[m[1]] = m[2]
	}
	return res
}

// Parsing Utilities

func findFirstComponent(html string) (int, int, string, bool, string, string) {
	loc := reComponentTag.FindStringIndex(html)
	if loc == nil {
		return -1, -1, "", false, "", ""
	}

	start := loc[0]
	tagOpenContent := html[loc[0]:loc[1]]

	// If we hit closing tag first </Comp>, ignore it, continue search
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

	// Self closing <Comp />
	if strings.HasSuffix(fullTagContent, "/>") {
		return start, attrEnd + 1, tagName, true, strings.TrimSuffix(attrsStr, "/"), ""
	}

	// Nested <Comp>...</Comp>
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

		// Check for nested opens in between
		chunk := html[searchStart:absClose]
		// Simple counting isn't perfect but works for well-formed HTML
		nestedOpens := 0
		openIdx := 0
		for {
			idx := strings.Index(chunk[openIdx:], openingTagPrefix)
			if idx == -1 {
				break
			}
			// Verify it's a tag
			charAfter := chunk[openIdx+idx+len(openingTagPrefix)]
			if charAfter == ' ' || charAfter == '>' || charAfter == '/' {
				nestedOpens++
			}
			openIdx += idx + 1
		}

		if nestedOpens == 0 {
			// Found matching close
			return start, absClose + len(closingTag), tagName, false, attrsStr, html[attrEnd+1 : absClose]
		}

		// Handle nesting logic more robustly
		// Actually, we can just jump ahead if we found nesting, but standard parsing is hard with Regex/Index.
		// Let's iterate block by block.
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

func hasSingleRoot(html string) bool {
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

// CSS Scoping

func processComponentStyles(raw string, scopeID string) (string, string, error) {
	loc := reStyleBlock.FindStringSubmatchIndex(raw)
	if loc == nil {
		return raw, "", nil
	}

	cssContent := raw[loc[2]:loc[3]]
	htmlContent := raw[:loc[0]] + raw[loc[1]:] // Remove style block from HTML
	htmlContent = strings.TrimSpace(htmlContent)

	var scopedCSS strings.Builder
	// Basic CSS Parser to inject scope
	// Note: This is a fragile regex parser. A real CSS parser is recommended for production.
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
			// Apply scope: p -> p[data-scope]
			newSels = append(newSels, fmt.Sprintf("%s[%s]", s, scopeID))
		}
		scopedCSS.WriteString(strings.Join(newSels, ", ") + " {" + body + "}\n")
	}

	return htmlContent, scopedCSS.String(), nil
}

func injectScopeID(html string, scopeID string) string {
	clean := strings.TrimSpace(html)
	// Find the first tag
	firstSpace := strings.IndexAny(clean, " />")
	if firstSpace == -1 {
		return html // Should not happen for valid html
	}
	// Inject attribute
	return clean[:firstSpace] + " " + scopeID + "=\"\"" + clean[firstSpace:]
}

// Watcher

func runWatch(basePath string) {
	fmt.Printf("Watching %s for changes...\n", basePath)
	if err := runCompile(basePath); err != nil {
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
				if strings.Contains(path, DirDist) {
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
			if err := runCompile(basePath); err != nil {
				fmt.Printf("❌ Compile Error: %v\n", err)
			} else {
				fmt.Println("✅ Built successfully.")
			}
			lastMod = time.Now()
		}
	}
}

// Helpers

func isPascalCase(s string) bool {
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

func isKebabCase(s string) bool {
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

func runTests(basePath string) {
	fmt.Printf("Running tests in %s...\n", basePath)

	failed := 0
	passed := 0

	err := filepath.Walk(basePath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if path == filepath.Join(basePath, "dist") || path == filepath.Join(basePath, "static") {
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Ext(path) != ".html" {
			return nil
		}

		if !strings.HasSuffix(path, "-test.html") {
			return nil
		}

		relPath, err := filepath.Rel(basePath, path)
		if err != nil {
			return err
		}

		fmt.Printf("  Testing %s... ", relPath)

		if err := runCompile(filepath.Dir(path)); err != nil {
			fmt.Printf("FAILED\n    Error: %v\n", err)
			failed++
		} else {
			fmt.Printf("PASSED\n")
			passed++
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error running tests: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n--- Test Results ---\n")
	fmt.Printf("Passed: %d\n", passed)
	fmt.Printf("Failed: %d\n", failed)

	if failed > 0 {
		os.Exit(1)
	}

	fmt.Println("All tests passed!")
}
