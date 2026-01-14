package main_test

import (
	"strings"
	"testing"

	"github.com/phillip-england/gtml/pkg/gtml"
)

// TestBasicFetchElement tests that a basic fetch element is correctly processed
func TestBasicFetchElement(t *testing.T) {
	html := `<div fetch='GET localhost:8080/api/users' as='users'></div>`

	result, err := gtml.ProcessFetchElements(html)
	if err != nil {
		t.Fatalf("ProcessFetchElements failed: %v", err)
	}

	// Should contain the modified div with an ID
	if !strings.Contains(result, "id=\"gtml-fetch-") {
		t.Error("Expected fetch element to have gtml-fetch ID")
	}

	// Should contain a script tag
	if !strings.Contains(result, "<script>") {
		t.Error("Expected script tag to be generated")
	}

	// Should contain fetch call to the URL
	if !strings.Contains(result, "fetch('localhost:8080/api/users'") {
		t.Error("Expected fetch call with correct URL")
	}

	// Should contain GET method
	if !strings.Contains(result, "method: 'GET'") {
		t.Error("Expected GET method in fetch options")
	}

	// Should use the 'as' name for the data
	if !strings.Contains(result, ".then(users =>") {
		t.Error("Expected data to be named 'users' from 'as' attribute")
	}
}

// TestFetchWithPOSTMethod tests that POST method is correctly handled
func TestFetchWithPOSTMethod(t *testing.T) {
	html := `<div fetch='POST localhost:8080/api/create' as='response'></div>`

	result, err := gtml.ProcessFetchElements(html)
	if err != nil {
		t.Fatalf("ProcessFetchElements failed: %v", err)
	}

	if !strings.Contains(result, "method: 'POST'") {
		t.Error("Expected POST method in fetch options")
	}
}

// TestFetchWithSuspense tests that suspense elements are correctly processed
func TestFetchWithSuspense(t *testing.T) {
	html := `<div fetch='GET localhost:8080/api/data' as='data'>
  <div suspense>
    <p>Loading...</p>
  </div>
  <p>Content</p>
</div>`

	result, err := gtml.ProcessFetchElements(html)
	if err != nil {
		t.Fatalf("ProcessFetchElements failed: %v", err)
	}

	// Should create suspense element
	if !strings.Contains(result, "data-gtml-suspense") {
		t.Error("Expected suspense element creation")
	}

	// Should contain loading text
	if !strings.Contains(result, "Loading...") {
		t.Error("Expected suspense content to be preserved")
	}

	// Should remove suspense on success
	if !strings.Contains(result, "suspense.remove()") {
		t.Error("Expected suspense removal on success")
	}
}

// TestFetchWithFallback tests that fallback elements are correctly processed
func TestFetchWithFallback(t *testing.T) {
	html := `<div fetch='GET localhost:8080/api/data' as='data'>
  <div fallback>
    <p>Error loading data</p>
  </div>
  <p>Content</p>
</div>`

	result, err := gtml.ProcessFetchElements(html)
	if err != nil {
		t.Fatalf("ProcessFetchElements failed: %v", err)
	}

	// Should create fallback element (hidden initially)
	if !strings.Contains(result, "data-gtml-fallback") {
		t.Error("Expected fallback element creation")
	}

	// Should contain error text
	if !strings.Contains(result, "Error loading data") {
		t.Error("Expected fallback content to be preserved")
	}

	// Should show fallback on error
	if !strings.Contains(result, "fallback.style.display = ''") {
		t.Error("Expected fallback to be shown on error")
	}
}

// TestFetchWithSuspenseAndFallback tests both suspense and fallback together
func TestFetchWithSuspenseAndFallback(t *testing.T) {
	html := `<div fetch='GET localhost:8080/api/users' as='users'>
  <div suspense>
    <p>Loading users...</p>
  </div>
  <div fallback>
    <p>Failed to load users</p>
  </div>
  <ul>
    <li>User list will go here</li>
  </ul>
</div>`

	result, err := gtml.ProcessFetchElements(html)
	if err != nil {
		t.Fatalf("ProcessFetchElements failed: %v", err)
	}

	// Should have both elements
	if !strings.Contains(result, "data-gtml-suspense") {
		t.Error("Expected suspense element")
	}
	if !strings.Contains(result, "data-gtml-fallback") {
		t.Error("Expected fallback element")
	}
}

// TestFetchWithForIteration tests basic for loop iteration
func TestFetchWithForIteration(t *testing.T) {
	html := `<div fetch='GET localhost:8080/api/users' as='users'>
  <ul>
    <li for='user in users'>{user.name}</li>
  </ul>
</div>`

	result, err := gtml.ProcessFetchElements(html)
	if err != nil {
		t.Fatalf("ProcessFetchElements failed: %v", err)
	}

	// Should contain data-gtml-for attribute
	if !strings.Contains(result, "data-gtml-for") {
		t.Error("Expected data-gtml-for attribute on template element")
	}

	// Should contain data-gtml-item attribute
	if !strings.Contains(result, "data-gtml-item=\"user\"") {
		t.Error("Expected data-gtml-item attribute with 'user'")
	}

	// Should contain data-gtml-source attribute
	if !strings.Contains(result, "data-gtml-source=\"users\"") {
		t.Error("Expected data-gtml-source attribute with 'users'")
	}

	// Should process the iteration in JavaScript using processForLoops
	if !strings.Contains(result, "processForLoops") {
		t.Error("Expected processForLoops function for iteration")
	}

	// Should have replaceExpressions helper
	if !strings.Contains(result, "replaceExpressions") {
		t.Error("Expected replaceExpressions helper function")
	}
}

// TestFetchWithNestedIteration tests nested for loops
func TestFetchWithNestedIteration(t *testing.T) {
	html := `<div fetch='GET localhost:8080/api/users' as='users'>
  <ul>
    <li for='user in users'>
      <p>{user.name}</p>
      <ul>
        <li for='color in user.colors'>{color.name}</li>
      </ul>
    </li>
  </ul>
</div>`

	result, err := gtml.ProcessFetchElements(html)
	if err != nil {
		t.Fatalf("ProcessFetchElements failed: %v", err)
	}

	// Should have nested iteration data attributes
	if !strings.Contains(result, "data-gtml-source=\"users\"") {
		t.Error("Expected outer iteration source")
	}
	if !strings.Contains(result, "data-gtml-source=\"user.colors\"") {
		t.Error("Expected nested iteration source with parent path")
	}
}

// TestParseForAttribute tests the for attribute parsing
func TestParseForAttribute(t *testing.T) {
	tests := []struct {
		input        string
		itemName     string
		sourcePath   string
		shouldError  bool
	}{
		{"user in users", "user", "users", false},
		{"item in items", "item", "items", false},
		{"color in user.colors", "color", "user.colors", false},
		{"child in parent.children", "child", "parent.children", false},
		{"invalid format", "", "", true},
		{"missing_in_keyword", "", "", true},
	}

	for _, tt := range tests {
		forLoop, err := gtml.ParseForAttribute(tt.input)

		if tt.shouldError {
			if err == nil {
				t.Errorf("ParseForAttribute(%q) should have errored", tt.input)
			}
			continue
		}

		if err != nil {
			t.Errorf("ParseForAttribute(%q) unexpected error: %v", tt.input, err)
			continue
		}

		if forLoop.ItemName != tt.itemName {
			t.Errorf("ParseForAttribute(%q) ItemName = %q, expected %q", tt.input, forLoop.ItemName, tt.itemName)
		}

		if forLoop.SourcePath != tt.sourcePath {
			t.Errorf("ParseForAttribute(%q) SourcePath = %q, expected %q", tt.input, forLoop.SourcePath, tt.sourcePath)
		}
	}
}

// TestMultipleFetchElements tests multiple fetch elements on the same page
func TestMultipleFetchElements(t *testing.T) {
	html := `<div>
  <div fetch='GET localhost:8080/api/users' as='users'>
    <ul><li for='user in users'>{user.name}</li></ul>
  </div>
  <div fetch='GET localhost:8080/api/posts' as='posts'>
    <ul><li for='post in posts'>{post.title}</li></ul>
  </div>
</div>`

	result, err := gtml.ProcessFetchElements(html)
	if err != nil {
		t.Fatalf("ProcessFetchElements failed: %v", err)
	}

	// Should have two different fetch IDs
	count := strings.Count(result, "id=\"gtml-fetch-")
	if count != 2 {
		t.Errorf("Expected 2 fetch IDs, got %d", count)
	}

	// Should have two script tags
	scriptCount := strings.Count(result, "<script>")
	if scriptCount != 2 {
		t.Errorf("Expected 2 script tags, got %d", scriptCount)
	}

	// Should have both URLs
	if !strings.Contains(result, "localhost:8080/api/users") {
		t.Error("Expected users API URL")
	}
	if !strings.Contains(result, "localhost:8080/api/posts") {
		t.Error("Expected posts API URL")
	}
}

// TestFetchElementPreservesOtherAttributes tests that other attributes are preserved
func TestFetchElementPreservesOtherAttributes(t *testing.T) {
	html := `<div class="container" data-test="value" fetch='GET /api/data' as='data'></div>`

	result, err := gtml.ProcessFetchElements(html)
	if err != nil {
		t.Fatalf("ProcessFetchElements failed: %v", err)
	}

	// Should preserve class attribute
	if !strings.Contains(result, "class=\"container\"") {
		t.Error("Expected class attribute to be preserved")
	}

	// Should preserve data-test attribute
	if !strings.Contains(result, "data-test=\"value\"") {
		t.Error("Expected data-test attribute to be preserved")
	}
}

// TestNoFetchElements tests that HTML without fetch elements passes through unchanged
func TestNoFetchElements(t *testing.T) {
	html := `<div class="container">
  <p>No fetch here</p>
</div>`

	result, err := gtml.ProcessFetchElements(html)
	if err != nil {
		t.Fatalf("ProcessFetchElements failed: %v", err)
	}

	// Should be unchanged
	if normalizeHTML(result) != normalizeHTML(html) {
		t.Error("Expected HTML without fetch elements to pass through unchanged")
	}
}

// TestFetchInComponent tests fetch elements work within components
func TestFetchInComponent(t *testing.T) {
	state := createTestState(map[string]string{
		"UserList": `<div props='apiUrl string'>
  <div fetch='GET {apiUrl}' as='users'>
    <ul><li for='user in users'>{user.name}</li></ul>
  </div>
</div>`,
	})

	html := `<UserList apiUrl='localhost:8080/api/users' />`
	result, err := gtml.CompileHTML(html, state, map[string]gtml.Value{})
	if err != nil {
		t.Fatalf("CompileHTML failed: %v", err)
	}

	// Should have compiled the fetch element
	if !strings.Contains(result, "gtml-fetch-") {
		t.Error("Expected fetch element to be compiled")
	}
}

// TestFetchScriptIsolation tests that fetch scripts are isolated using IIFE
func TestFetchScriptIsolation(t *testing.T) {
	html := `<div fetch='GET /api/data' as='data'></div>`

	result, err := gtml.ProcessFetchElements(html)
	if err != nil {
		t.Fatalf("ProcessFetchElements failed: %v", err)
	}

	// Should use IIFE for isolation
	if !strings.Contains(result, "(function() {") {
		t.Error("Expected IIFE for script isolation")
	}
	if !strings.Contains(result, "})();") {
		t.Error("Expected IIFE closing for script isolation")
	}
}

// TestFetchExpressionReplacement tests the expression replacement helper
func TestFetchExpressionReplacement(t *testing.T) {
	html := `<div fetch='GET /api/users' as='users'>
  <li for='user in users'><span>{user.name}</span><span>{user.email}</span></li>
</div>`

	result, err := gtml.ProcessFetchElements(html)
	if err != nil {
		t.Fatalf("ProcessFetchElements failed: %v", err)
	}

	// Should have replaceExpressions function
	if !strings.Contains(result, "function replaceExpressions(html, scope)") {
		t.Error("Expected replaceExpressions helper function")
	}

	// The template should still contain expression syntax for client-side replacement
	if !strings.Contains(result, "{user.name}") {
		t.Error("Expected expressions to be preserved for client-side replacement")
	}
}

// TestFetchTemplateHidden tests that for templates are initially hidden
func TestFetchTemplateHidden(t *testing.T) {
	html := `<div fetch='GET /api/items' as='items'>
  <li for='item in items'>{item.name}</li>
</div>`

	result, err := gtml.ProcessFetchElements(html)
	if err != nil {
		t.Fatalf("ProcessFetchElements failed: %v", err)
	}

	// Template elements should be hidden initially
	if !strings.Contains(result, "style=\"display:none\"") {
		t.Error("Expected template elements to be hidden initially")
	}
}

// TestFetchErrorHandling tests proper error handling in generated script
func TestFetchErrorHandling(t *testing.T) {
	html := `<div fetch='GET /api/data' as='data'>
  <div fallback><p>Error</p></div>
</div>`

	result, err := gtml.ProcessFetchElements(html)
	if err != nil {
		t.Fatalf("ProcessFetchElements failed: %v", err)
	}

	// Should have catch block
	if !strings.Contains(result, ".catch(error =>") {
		t.Error("Expected catch block for error handling")
	}

	// Should log error
	if !strings.Contains(result, "console.error('Fetch error:'") {
		t.Error("Expected error logging")
	}
}

// TestFetchResponseValidation tests that non-ok responses are handled
func TestFetchResponseValidation(t *testing.T) {
	html := `<div fetch='GET /api/data' as='data'></div>`

	result, err := gtml.ProcessFetchElements(html)
	if err != nil {
		t.Fatalf("ProcessFetchElements failed: %v", err)
	}

	// Should check response.ok
	if !strings.Contains(result, "if (!response.ok) throw new Error") {
		t.Error("Expected response.ok validation")
	}
}
