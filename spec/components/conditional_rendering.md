# Conditional Rendering with Ternary Operators

## How Conditionals Work
In `gtml`, conditional rendering is done using ternary operators within expressions. The syntax is:
```
{ condition ? (truthy content) : (falsy content) }
```

This allows you to conditionally render HTML based on prop values.

## Basic Ternary Example
For example, let's take the following component:
`./myapp/components/Conditional.html`
```html
<div props='x int'>
  { x == 2 ? (
    <p>I am only shown if x equals 2</p>
  ) : (
    <p>I am shown if x does not equal 2</p>
  ) }
</div>
```

Then, we can use this component like so:
```html
<Conditional x={2} />
<Conditional x={0} />
```

Which should output:
```html
<div>
  <p>I am only shown if x equals 2</p>
</div>
<div>
  <p>I am shown if x does not equal 2</p>
</div>
```

## Not Equal Comparison
We can use `!=` to check for inequality:
`./myapp/components/NegativeCond.html`
```html
<div props='x int'>
  { x != 2 ? (
    <p>I am shown if x does not equal 2</p>
  ) : (
    <p>I am shown if x is equal to 2</p>
  ) }
</div>
```

## String Comparisons
We can also compare strings:
`./myapp/components/StringCond.html`
```html
<div props='color string'>
  <h1>Favorite Color</h1>
  { color == "blue" ? (
    <p>blue</p>
  ) : (
    <p>not blue</p>
  ) }
</div>
```

## Nested Ternary for Multiple Conditions
For multiple conditions (similar to else-if), you can nest ternary operators:
`./myapp/components/CanYouDrink.html`
```html
<div props='age int'>
  { age < 21 ? (
    <p>You cannot drink</p>
  ) : (
    { age > 21 && age < 25 ? (
      <p>You can drink but you shouldn't, your brain is not fully formed.</p>
    ) : (
      <p>You can drink</p>
    ) }
  ) }
</div>
```

The above component makes use of `&&` which allows us to chain conditions. It is our way of saying, "both conditions on both sides must resolve to true."

We can also use `||` to say, "one or more of the conditions on both sides must resolve to true." In this way, `gtml` offers conditional primitives for `and` and `or`.

The above component may be used like so:
```html
<CanYouDrink age={20} />
<CanYouDrink age={23} />
<CanYouDrink age={29} />
```

which should resolve as:
```html
<div>
  <p>You cannot drink</p>
</div>
<div>
  <p>You can drink but you shouldn't, your brain is not fully formed.</p>
</div>
<div>
  <p>You can drink</p>
</div>
```

## Using Props in Both Conditions and Content
Props used in ternary conditions can also be displayed elsewhere in the component:

`./myapp/components/ShowValue.html`
```html
<div props='x int'>
  <p>The value of x is: {x}</p>
  { x == 5 ? (
    <p>x is five!</p>
  ) : (
    <p>x is not five</p>
  ) }
</div>
```

When used like this:
```html
<ShowValue x={5} />
```

It would output:
```html
<div>
  <p>The value of x is: 5</p>
  <p>x is five!</p>
</div>
```

This works because `x` is declared as a prop in the `props` attribute. The ternary checks the value of `x`, and `x` is also accessible throughout the entire component via the expression syntax `{x}`.

## Comparison Operators
The following comparison operators are available in ternary conditions:
- `==` (equals)
- `!=` (not equals)
- `<` (less than)
- `>` (greater than)
- `<=` (less than or equal)
- `>=` (greater than or equal)

## Logical Operators
The following logical operators are available:
- `&&` (and) - both conditions must be true
- `||` (or) - at least one condition must be true
