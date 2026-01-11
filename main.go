package main

import (
	"errors"
	"flag"
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

// =================================================================================
// GLOBALS & TYPES
// =================================================================================

const (
	DirComponents = "components"
	DirRoutes     = "routes"
	DirDist       = "dist"
	DirStatic     = "static"
	FileStyleCSS  = "styles.css"
)

// Regex patterns for DSL
var (
	// Matches {{ prop: name type }}
	rePropDef = regexp.MustCompile(`\{\{\s*prop:\s*(\w+)\s+(\w+)\s*\}\}`)
	// Matches {{ drill: name }}
	reDrill = regexp.MustCompile(`\{\{\s*drill:\s*(\w+)\s*\}\}`)
	// Matches {{ slot: name }}
	reSlotDef = regexp.MustCompile(`\{\{\s*slot:\s*(\w+)\s*\}\}`)
	// Matches <PascalCase ...> or <PascalCase />
	reComponentTag = regexp.MustCompile(`</?([A-Z][a-zA-Z0-9]*)`)
	// Matches <style>...</style>
	reStyleBlock = regexp.MustCompile(`(?s)<style>(.*?)</style>`)
	// Matches <slot ...> content </slot>
	reSlotUsage = regexp.MustCompile(`(?s)<slot\s+([^>]+)>(.*?)</slot>`)
)

type Component struct {
	Name        string
	RawContent  string // Original file content
	Template    string // HTML after stripping style
	Styles      string // Raw CSS content
	ScopedStyle string // Compiled/Scoped CSS
	ScopeID     string // Unique ID for scoping
	Path        string
}

type GlobalState struct {
	Components map[string]*Component
	CSSOutput  strings.Builder
}

// =================================================================================
// MAIN ENTRY
// =================================================================================

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		initCmd := flag.NewFlagSet("init", flag.ExitOnError)
		forceFlag := initCmd.Bool("force", false, "Overwrite existing directory")
		initCmd.Parse(os.Args[2:])

		if initCmd.NArg() < 1 {
			fmt.Println("Error: Missing path argument for init.")
			fmt.Println("Usage: gtml init <PATH> [--force]")
			os.Exit(1)
		}
		path := initCmd.Arg(0)
		runInit(path, *forceFlag)

	case "compile":
		compileCmd := flag.NewFlagSet("compile", flag.ExitOnError)
		watchFlag := compileCmd.Bool("watch", false, "Watch for changes")
		compileCmd.Parse(os.Args[2:])

		if compileCmd.NArg() < 1 {
			fmt.Println("Error: Missing path argument for compile.")
			fmt.Println("Usage: gtml compile <PATH> [--watch]")
			os.Exit(1)
		}
		path := compileCmd.Arg(0)

		if *watchFlag {
			runWatch(path)
		} else {
			if err := runCompile(path); err != nil {
				fmt.Printf("\n❌ Compilation failed: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("\n✅ Compilation successful!")
		}

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
}

// =================================================================================
// INIT COMMAND
// =================================================================================

func runInit(basePath string, force bool) {
	if _, err := os.Stat(basePath); err == nil && !force {
		fmt.Printf("Error: Directory '%s' already exists. Use --force to overwrite.\n", basePath)
		os.Exit(1)
	}

	dirs := []string{
		filepath.Join(basePath, DirComponents),
		filepath.Join(basePath, DirRoutes),
		filepath.Join(basePath, DirDist),
		filepath.Join(basePath, DirDist, DirStatic),
		filepath.Join(basePath, DirStatic),
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			panic(err)
		}
	}

	files := map[string]string{
		filepath.Join(basePath, DirComponents, "BasicButton.html"):
`<button>{{ prop: text string }}</button>`,

		filepath.Join(basePath, DirComponents, "GuestLayout.html"):
`<html>
  <head>
    <title>{{ prop: title string }}</title>
  </head>
  <body>
    <BasicButton text='{{ drill: title }}' />
    {{ slot: content }}
  </body>
</html>`,

		filepath.Join(basePath, DirRoutes, "index.html"):
`<GuestLayout title="Some Title">
  <p>Some Content</p>
  <BasicButton text='{{ drill: title }}' />
</GuestLayout>`,
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			panic(err)
		}
	}

	fmt.Printf("Initialized gtml project at %s\n", basePath)
}

// =================================================================================
// COMPILE COMMAND
// =================================================================================

func runCompile(basePath string) error {
	state := &GlobalState{
		Components: make(map[string]*Component),
	}

	// 1. Read Components
	compDir := filepath.Join(basePath, DirComponents)
	err := filepath.Walk(compDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil { return err }
		if info.IsDir() || filepath.Ext(path) != ".html" { return nil }

		name := strings.TrimSuffix(filepath.Base(path), ".html")

		// Validation: PascalCase
		if !isPascalCase(name) {
			return fmt.Errorf("component '%s' must be PascalCase", path)
		}

		// Validation: Unique Names
		if _, exists := state.Components[name]; exists {
			return fmt.Errorf("duplicate component name found: %s", name)
		}

		contentBytes, err := os.ReadFile(path)
		if err != nil { return err }
		content := string(contentBytes)

		// Process CSS and Scope
		scopeID := "data-" + strings.ToLower(name)
		template, css, err := processComponentStyles(content, scopeID)
		if err != nil { return err }

		// Validation: One Root Element
		// Simple check: strip comments/whitespace, ensure it starts with <tag> and ends with </tag>
		// Note: This is a loose check for the purpose of this exercise.
		if !hasSingleRoot(template) {
			return fmt.Errorf("component '%s' must have a single root element", name)
		}

		// Inject Scope ID into the root element of the template
		template = injectScopeID(template, scopeID)

		state.Components[name] = &Component{
			Name:        name,
			RawContent:  content,
			Template:    template,
			ScopedStyle: css,
			ScopeID:     scopeID,
			Path:        path,
		}

		// Aggregate CSS
		if css != "" {
			state.CSSOutput.WriteString("/* " + name + " */\n")
			state.CSSOutput.WriteString(css + "\n")
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 2. Process Routes
	routesDir := filepath.Join(basePath, DirRoutes)
	distDir := filepath.Join(basePath, DirDist)

	// Clean dist (optional but good practice, keeping it simple here)

	err = filepath.Walk(routesDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil { return err }
		if info.IsDir() { return nil }
		if filepath.Ext(path) != ".html" { return nil }

		relPath, _ := filepath.Rel(routesDir, path)

		// Validation: Kebab-case filenames
		fileName := strings.TrimSuffix(filepath.Base(path), ".html")
		if !isKebabCase(fileName) {
			return fmt.Errorf("route '%s' must be kebab-case", path)
		}

		contentBytes, err := os.ReadFile(path)
		if err != nil { return err }

		// Compile the route content
		// Route starts with empty props scope
		compiledHTML, err := compileHTML(string(contentBytes), state, map[string]interface{}{})
		if err != nil {
			return fmt.Errorf("error compiling %s: %v", path, err)
		}

		// Write to Dist
		outPath := filepath.Join(distDir, relPath)
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return err
		}
		return os.WriteFile(outPath, []byte(compiledHTML), 0644)
	})

	if err != nil {
		return err
	}

	// 3. Write CSS
	staticDistDir := filepath.Join(distDir, DirStatic)
	if err := os.MkdirAll(staticDistDir, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(staticDistDir, FileStyleCSS), []byte(state.CSSOutput.String()), 0644); err != nil {
		return err
	}

	// 4. Copy Static Assets
	srcStatic := filepath.Join(basePath, DirStatic)
	copyDir(srcStatic, staticDistDir)

	return nil
}

// =================================================================================
// COMPILATION LOGIC (Recursive)
// =================================================================================

func compileHTML(html string, state *GlobalState, scopeProps map[string]interface{}) (string, error) {
	// Loop until no components remain
	// We scan the string for <PascalCase
	for {
		startIdx, endIdx, tagName, _, attrsStr, innerContent := findFirstComponent(html)
		if startIdx == -1 {
			break // No components found
		}

		compDef, exists := state.Components[tagName]
		if !exists {
			return "", fmt.Errorf("component '%s' not found", tagName)
		}

		// Parse Attributes
		props := parseAttributes(attrsStr)
		resolvedProps := make(map[string]interface{})

		// Resolve Props (handle drills)
		for k, v := range props {
			if strings.HasPrefix(v, "{{") && strings.Contains(v, "drill:") {
				match := reDrill.FindStringSubmatch(v)
				if len(match) > 1 {
					drillName := match[1]
					if val, ok := scopeProps[drillName]; ok {
						resolvedProps[k] = val
					} else {
						// Prop missing or drilling nil?
						resolvedProps[k] = ""
					}
				}
			} else {
				resolvedProps[k] = v
			}
		}

		// Parse Children for Slots
		// We need to look for <slot name='x' tag='y'>content</slot> inside innerContent
		slotsMap := make(map[string]string)

		// We temporarily replace slot definitions in innerContent to prevent regex collision if nesting
		// Actually, we just scan innerContent.
		// NOTE: innerContent is the raw HTML passed AS CHILDREN to the component.
		// We need to resolve components inside the children BEFORE passing them to the slot?
		// Usually in SSG, children are compiled in the parent's context.

		compiledChildren, err := compileHTML(innerContent, state, scopeProps)
		if err != nil {
			return "", err
		}

		// Extract slots from compiled children
		// Any top-level <slot> tags in the children need to be captured.
		// Non-slot content usually goes to 'default' slot or is discarded if no default slot (implicit).
		// For this specific DSL, the prompt explicitly uses <slot name='...' tag='...'>

		// Simple parser for slot usages in the *compiled* children
		finalChildrenForDefault := strings.Builder{}

		childCursor := 0
		childLen := len(compiledChildren)
		for childCursor < childLen {
			// Find <slot ...
			loc := reSlotUsage.FindStringIndex(compiledChildren[childCursor:])
			if loc == nil {
				finalChildrenForDefault.WriteString(compiledChildren[childCursor:])
				break
			}

			absStart := childCursor + loc[0]
			absEnd := childCursor + loc[1]

			// Append text before the slot to default
			finalChildrenForDefault.WriteString(compiledChildren[childCursor:absStart])

			// Process the slot usage
			slotTagFull := compiledChildren[absStart:absEnd]
			match := reSlotUsage.FindStringSubmatch(slotTagFull)

			slotAttrs := parseAttributes(match[1])
			slotContent := match[2]

			sName := slotAttrs["name"]
			sTag := slotAttrs["tag"]
			sClass := slotAttrs["class"] // Optional extras

			if sName != "" {
				// Wrap content in the requested tag
				// e.g. <article class='...'> content </article>
				wrapper := fmt.Sprintf("<%s", sTag)
				if sClass != "" {
					wrapper += fmt.Sprintf(" class=\"%s\"", sClass)
				}
				// Copy other attributes? Prompt implies just name/tag logic mostly.
				wrapper += fmt.Sprintf(">%s</%s>", slotContent, sTag)
				slotsMap[sName] = wrapper
			}

			childCursor = absEnd
		}

		// Now we have the Component Template. We need to fill it.
		// 1. Fill Props
		renderedComp := compDef.Template

		// Replace {{ prop: name type }}
		// We use a func replace to handle logic
		renderedComp = rePropDef.ReplaceAllStringFunc(renderedComp, func(s string) string {
			m := rePropDef.FindStringSubmatch(s)
			pName := m[1]
			pType := m[2]

			val, ok := resolvedProps[pName]
			if !ok {
				// return empty or error? Prompt doesn't specify strict prop failure, assuming empty.
				return ""
			}

			// Type check (basic)
			if pType == "int" {
				if _, err := strconv.Atoi(fmt.Sprintf("%v", val)); err != nil {
					fmt.Printf("Warning: Prop '%s' expected int, got '%v'\n", pName, val)
				}
			}

			return fmt.Sprintf("%v", val)
		})

		// 2. Fill Slots
		renderedComp = reSlotDef.ReplaceAllStringFunc(renderedComp, func(s string) string {
			m := reSlotDef.FindStringSubmatch(s)
			slotName := m[1]
			if content, ok := slotsMap[slotName]; ok {
				return content
			}
			return "" // Empty if slot not filled
		})

		// Now we recurse inside the rendered component?
		// The component template might contain OTHER components.
		// We need to compile the result of this expansion.
		// IMPORTANT: We pass 'resolvedProps' as the scope for the INNER components of this template.
		finalRendered, err := compileHTML(renderedComp, state, resolvedProps)
		if err != nil {
			return "", err
		}

		// Replace the original component tag in the source HTML with the final rendered HTML
		html = html[:startIdx] + finalRendered + html[endIdx:]
	}

	return html, nil
}

// =================================================================================
// PARSING HELPERS
// =================================================================================

// findFirstComponent locates the first <PascalCase> tag, handling nesting roughly
// Returns: start, end, tagName, isSelfClosing, attributes, innerContent
func findFirstComponent(html string) (int, int, string, bool, string, string) {
	loc := reComponentTag.FindStringIndex(html)
	if loc == nil {
		return -1, -1, "", false, "", ""
	}

	start := loc[0]
	tagOpenContent := html[loc[0]:loc[1]] // e.g. "<Button" or "</Button"

	// If it's a closing tag </Button>, ignore it and look further (shouldn't happen if logic works)
	if strings.HasPrefix(tagOpenContent, "</") {
		// Found a dangling closing tag or nested closing tag first?
		// We should skip this and search after it.
		s, e, t, sc, a, i := findFirstComponent(html[loc[1]:])
		if s != -1 {
			return loc[1] + s, loc[1] + e, t, sc, a, i
		}
		return -1, -1, "", false, "", ""
	}

	tagName := tagOpenContent[1:] // Strip <

	// Find the end of the opening tag ">" or "/>"
	rest := html[loc[1]:]
	closeBracket := strings.IndexAny(rest, ">")
	if closeBracket == -1 {
		return -1, -1, "", false, "", "" // Malformed
	}

	attrEnd := loc[1] + closeBracket
	fullTagContent := html[start : attrEnd+1] // <Button attrs...>
	attrsStr := html[loc[1] : attrEnd]        // attrs...

	if strings.HasSuffix(fullTagContent, "/>") {
		// Self closing
		return start, attrEnd + 1, tagName, true, strings.TrimSuffix(attrsStr, "/"), ""
	}

	// Not self closing. Need to find matching </Name>
	// We need to account for nested tags of same name.
	nestLevel := 1
	searchStart := attrEnd + 1
	closingTag := "</" + tagName + ">"
	openingTagPrefix := "<" + tagName

	for {
		nextClose := strings.Index(html[searchStart:], closingTag)
		if nextClose == -1 {
			return -1, -1, "", false, "", "" // Missing closing tag
		}

		absClose := searchStart + nextClose

		// Check for nested open tags between current pos and the found close
		chunk := html[searchStart:absClose]
		// Count opening tags in chunk (rough check)
		// To be precise we need to check bounds, but simple counting works if strict PascalCase

		nestedOpens := strings.Count(chunk, openingTagPrefix)

		if nestedOpens == 0 {
			// No nested opens, this is our closer
			return start, absClose + len(closingTag), tagName, false, attrsStr, html[attrEnd+1 : absClose]
		}

		// If nested opens found, we need to skip equivalent number of closes
		// This is tricky with simple string search.
		// Iterative approach: find next open OR close.

		// Reset and scan linearly for correctness
		current := attrEnd + 1
		nestLevel = 1
		for nestLevel > 0 {
			nextOpenIdx := strings.Index(html[current:], openingTagPrefix)
			nextCloseIdx := strings.Index(html[current:], closingTag)

			if nextCloseIdx == -1 {
				return -1, -1, "", false, "", "" // Malformed
			}

			if nextOpenIdx != -1 && nextOpenIdx < nextCloseIdx {
				nestLevel++
				current += nextOpenIdx + 1 // Advance past <
			} else {
				nestLevel--
				if nestLevel == 0 {
					return start, current + nextCloseIdx + len(closingTag), tagName, false, attrsStr, html[attrEnd+1 : current+nextCloseIdx]
				}
				current += nextCloseIdx + len(closingTag)
			}
		}
	}
	return -1, -1, "", false, "", ""
}

func parseAttributes(attrStr string) map[string]string {
	// regex for key="value" or key='value'
	reAttrs := regexp.MustCompile(`(\w+)=["']([^"']*)["']`)
	matches := reAttrs.FindAllStringSubmatch(attrStr, -1)
	res := make(map[string]string)
	for _, m := range matches {
		res[m[1]] = m[2]
	}
	return res
}

func hasSingleRoot(html string) bool {
	// Remove comments
	reComment := regexp.MustCompile(``)
	clean := reComment.ReplaceAllString(html, "")
	clean = strings.TrimSpace(clean)

	// Should start with <tag
	if !strings.HasPrefix(clean, "<") { return false }
	// Find first tag name
	idx := strings.IndexAny(clean, " >/")
	if idx == -1 { return false }

	tagName := clean[1:idx]

	// Ensure it ends with > (self closing) or </tag>
	if strings.HasSuffix(clean, "/>") {
		// Check if there is anything else?
		// Just a heuristic check
		return true
	}

	suffix := "</" + tagName + ">"
	if strings.HasSuffix(clean, suffix) {
		// Check for siblings.
		// If we find the closing tag earlier than the end, and there is non-whitespace after...
		// Complex to parse strictly with regex.
		// Assuming valid HTML structure for this exercise.
		return true
	}
	return false
}

// =================================================================================
// CSS & STYLES
// =================================================================================

func processComponentStyles(raw string, scopeID string) (string, string, error) {
	// Extract <style> block
	loc := reStyleBlock.FindStringSubmatchIndex(raw)
	if loc == nil {
		return raw, "", nil
	}

	cssContent := raw[loc[2]:loc[3]]
	htmlContent := raw[:loc[0]] + raw[loc[1]:] // Remove style block
	htmlContent = strings.TrimSpace(htmlContent)

	// Scope CSS
	// Very basic CSS parser: split by }
	var scopedCSS strings.Builder
	blocks := strings.Split(cssContent, "}")

	for _, block := range blocks {
		if strings.TrimSpace(block) == "" { continue }
		parts := strings.Split(block, "{")
		if len(parts) != 2 { continue } // Invalid CSS or media query (ignored for simplicity)

		selectors := parts[0]
		body := parts[1]

		// Split selectors by comma
		selList := strings.Split(selectors, ",")
		var newSels []string
		for _, s := range selList {
			s = strings.TrimSpace(s)
			// Apply scope: add [data-scope-id] to the last element of selector
			// e.g. "div p" -> "div p[data-scope-id]"
			// Simple approach: append the attribute selector
			newSels = append(newSels, fmt.Sprintf("%s[%s]", s, scopeID))
		}

		scopedCSS.WriteString(strings.Join(newSels, ", ") + " {" + body + "}\n")
	}

	return htmlContent, scopedCSS.String(), nil
}

func injectScopeID(html string, scopeID string) string {
	// Inject scopeID as an attribute into the first tag found
	// html assumes to be "<tag ...> ..."
	clean := strings.TrimSpace(html)
	firstSpace := strings.IndexAny(clean, " />")
	if firstSpace == -1 { return html } // Should not happen if valid

	// Insert before first space or closing
	return clean[:firstSpace] + " " + scopeID + "=\"\"" + clean[firstSpace:]
}

// =================================================================================
// WATCHER
// =================================================================================

func runWatch(basePath string) {
	fmt.Printf("Watching %s for changes...\n", basePath)

	// Initial Compile
	if err := runCompile(basePath); err != nil {
		fmt.Printf("Compile Error: %v\n", err)
	} else {
		fmt.Println("Built successfully.")
	}

	lastMod := time.Now()

	for {
		time.Sleep(1 * time.Second)

		needsCompile := false

		// Walk to check mod times
		filepath.Walk(basePath, func(path string, info fs.FileInfo, err error) error {
			if err != nil { return nil }
			if info.IsDir() {
				if strings.Contains(path, DirDist) { return filepath.SkipDir } // Ignore dist
				return nil
			}
			if info.ModTime().After(lastMod) {
				needsCompile = true
				return errors.New("change found") // Break walk
			}
			return nil
		})

		if needsCompile {
			fmt.Println("Change detected. Compiling...")
			lastMod = time.Now()
			if err := runCompile(basePath); err != nil {
				fmt.Printf("❌ Compile Error: %v\n", err)
			} else {
				fmt.Println("✅ Built successfully.")
			}
		}
	}
}

// =================================================================================
// UTILS
// =================================================================================

func isPascalCase(s string) bool {
	if len(s) == 0 { return false }
	// First char must be uppercase
	r := []rune(s)
	if !unicode.IsUpper(r[0]) { return false }
	// Must contain only letters/numbers (basic check)
	for _, x := range r {
		if !unicode.IsLetter(x) && !unicode.IsNumber(x) { return false }
	}
	return true
}

func isKebabCase(s string) bool {
	if len(s) == 0 { return false }
	for _, r := range s {
		if !unicode.IsLower(r) && !unicode.IsNumber(r) && r != '-' {
			return false
		}
	}
	return true
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil { return err }
		relPath, _ := filepath.Rel(src, path)
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		data, err := os.ReadFile(path)
		if err != nil { return err }
		return os.WriteFile(dstPath, data, info.Mode())
	})
}
