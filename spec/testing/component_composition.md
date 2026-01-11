# Testing Component Composition

## Test Purpose
Validate that components can be nested within other components correctly.

## Basic Nesting

`./myapp/components/Card.html`
```html
<div class='card' props='title string'>
  <h2>{title}</h2>
  {{ slot: body }}
</div>
```

`./myapp/components/Button.html`
```html
<button props='label string'>{label}</button>
```

Used together:
```html
<Card title='My Card'>
  <slot name='body' tag='div'>
    <p>Card content</p>
    <Button label='Click me' />
  </slot>
</Card>
```

Should produce:
```html
<div class='card'>
  <h2>My Card</h2>
  <div>
    <p>Card content</p>
    <button>Click me</button>
  </div>
</div>
```

## Deeply Nested Components

Test components nested 3+ levels deep to ensure recursive compilation works.

## Component Using Itself (Recursion Guard)

A component should NOT be able to use itself:
```html
<div props='text string'>
  <SelfReference text={text} />
</div>
```

This should error or have defined behavior to prevent infinite recursion.

## Circular Component References

Component A uses Component B, and Component B uses Component A:

`./myapp/components/CompA.html`
```html
<div>
  <CompB />
</div>
```

`./myapp/components/CompB.html`
```html
<div>
  <CompA />
</div>
```

This should error to prevent infinite recursion.
