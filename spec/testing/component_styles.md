# Testing Component Styles

## Test Purpose
Validate that component styles are scoped correctly and compiled to `styles.css`.

## Basic Component Style

Take this component:
`./myapp/components/StyledButton.html`
```html
<style>
  button {
    background: blue;
    color: white;
  }
</style>

<button props='text string'>{text}</button>
```

The styles should be scoped so they only apply to this component.

## Style Scoping Prevents Conflicts

Take two components:

`./myapp/components/RedParagraph.html`
```html
<style>
  p {
    color: red;
  }
</style>

<p props='text string'>{text}</p>
```

`./myapp/components/BlueParagraph.html`
```html
<style>
  p {
    color: blue;
  }
</style>

<p props='text string'>{text}</p>
```

When both are used in a route:
```html
<RedParagraph text='Red text' />
<BlueParagraph text='Blue text' />
```

The styles should not conflict. Each paragraph should have its correct color.

## Styles Output Location

All component styles should be compiled into:
```
./myapp/dist/static/styles.css
```

## Component Without Styles

A component without a `<style>` section should work normally:
```html
<button props='text string'>{text}</button>
```

## Style With Complex Selectors

```html
<style>
  .container > p:first-child {
    font-weight: bold;
  }
  .container p + p {
    margin-top: 1rem;
  }
</style>

<div class='container' props='title string'>
  <p>{title}</p>
  <p>More content</p>
</div>
```

Complex selectors should be scoped correctly.

## Empty Style Block

```html
<style></style>

<div props='text string'>{text}</div>
```

Should not cause errors.
