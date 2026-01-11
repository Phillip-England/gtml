# Testing Full Compilation

## Test Purpose
End-to-end tests that validate the complete compilation workflow.

## Simple Route Compilation

`./myapp/routes/index.html`
```html
<html>
  <head>
    <title>Home</title>
  </head>
  <body>
    <h1>Welcome</h1>
  </body>
</html>
```

Should compile to `./myapp/dist/index.html` with the same content.

## Route With Components

`./myapp/components/Header.html`
```html
<header props='title string'>
  <h1>{title}</h1>
</header>
```

`./myapp/routes/index.html`
```html
<html>
  <body>
    <Header title='My Site' />
    <p>Content</p>
  </body>
</html>
```

Should compile to:
```html
<html>
  <body>
    <header>
      <h1>My Site</h1>
    </header>
    <p>Content</p>
  </body>
</html>
```

## Route With Nested Directory

`./myapp/routes/blog/post.html` should compile to `./myapp/dist/blog/post.html`.

## Multiple Routes

Compiling a project with multiple routes should generate all output files:
- `./myapp/routes/index.html` → `./myapp/dist/index.html`
- `./myapp/routes/about.html` → `./myapp/dist/about.html`
- `./myapp/routes/contact.html` → `./myapp/dist/contact.html`

## Full Integration Test

Create a complete project with:
- Multiple components (some with styles)
- Multiple routes
- Use of props, slots, and conditionals
- Static assets

Verify:
1. All routes compile correctly
2. All components resolve
3. Styles are combined into `styles.css`
4. Static assets are copied
5. No errors or warnings
