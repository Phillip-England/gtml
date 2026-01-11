# Testing Routes Directory Rules

## Test Purpose
Validate that `gtml` enforces all route directory rules correctly.

## Kebab-Case Names Required

### Valid Kebab-Case Names
These should be accepted:
- `index.html`
- `about.html`
- `contact-us.html`
- `blog-post-1.html`

### Invalid Non-Kebab-Case Names
These should produce errors:
- `Index.html` (starts with uppercase)
- `contactUs.html` (camelCase)
- `contact_us.html` (snake_case)
- `ContactUs.html` (PascalCase)

## Repeated Names In Different Subdirectories

### Same Name Different Directories Allowed
These should both be valid:
- `./myapp/routes/blog/index.html`
- `./myapp/routes/about/index.html`

## Output Mirrors Directory Structure
The compiled output should mirror the routes directory:
- `./myapp/routes/index.html` → `./myapp/dist/index.html`
- `./myapp/routes/blog/post.html` → `./myapp/dist/blog/post.html`
- `./myapp/routes/about/team.html` → `./myapp/dist/about/team.html`
