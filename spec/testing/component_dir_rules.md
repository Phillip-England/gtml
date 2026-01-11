# Testing Components Directory Rules

## Test Purpose
Validate that `gtml` enforces all component directory rules correctly.

## PascalCase Names Required

### Valid PascalCase Names
These should be accepted:
- `BasicButton.html`
- `UserCard.html`
- `HTTPClient.html`
- `MyComponent.html`

### Invalid Non-PascalCase Names
These should produce errors:
- `basicbutton.html` (all lowercase)
- `basic-button.html` (kebab-case)
- `basic_button.html` (snake_case)
- `basicButton.html` (camelCase, does not start with uppercase)

## Unique Names Across Subdirectories

### Duplicate Names Should Error
If we have both:
- `./myapp/components/Button.html`
- `./myapp/components/buttons/Button.html`

Compilation should fail with an error indicating duplicate component names.

### Same Name Different Extensions
If we have both:
- `./myapp/components/Button.html`
- `./myapp/components/Button.css`

Only `.html` files are components, so this should be allowed.

## One Component Per File

### Valid Single Root Element
```html
<div props='text string'>
  <p>{text}</p>
  <span>More content</span>
</div>
```

### Invalid Multiple Root Elements
This should error:
```html
<p>One</p>
<p>Two</p>
```

### Text Before Element
This should error:
```html
Some text
<div>Content</div>
```

### Comment Before Element
A comment before the root element should be allowed:
```html
<!-- This is a comment -->
<div>Content</div>
```

## Subdirectories Are Allowed
Components can be organized in subdirectories:
- `./myapp/components/buttons/PrimaryButton.html`
- `./myapp/components/layouts/MainLayout.html`
- `./myapp/components/forms/TextInput.html`

All should be accessible by their filename (without path).
