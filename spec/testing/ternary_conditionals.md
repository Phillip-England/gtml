# Testing Ternary Conditionals

## Test Purpose
Validate that conditional rendering with ternary operators `{ condition ? (truthy) : (falsy) }` works correctly.

## Basic Ternary

Take this component:
`./myapp/components/ShowIfActive.html`
```html
<div props='active boolean'>
  { active == true ? (
    <p>Active</p>
  ) : (
    <p>Inactive</p>
  ) }
</div>
```

Used like:
```html
<ShowIfActive active={true} />
```

Should produce:
```html
<div>
  <p>Active</p>
</div>
```

Used like:
```html
<ShowIfActive active={false} />
```

Should produce:
```html
<div>
  <p>Inactive</p>
</div>
```

## Numeric Comparisons

### Greater Than Or Equal
Take this component:
`./myapp/components/AgeCheck.html`
```html
<div props='age int'>
  { age >= 18 ? (
    <p>Adult</p>
  ) : (
    <p>Minor</p>
  ) }
</div>
```

Used like:
```html
<AgeCheck age={21} />
```

Should produce:
```html
<div>
  <p>Adult</p>
</div>
```

Used like:
```html
<AgeCheck age={15} />
```

Should produce:
```html
<div>
  <p>Minor</p>
</div>
```

### All Comparison Operators
Test each operator:
- `==` (equals)
- `!=` (not equals)
- `<` (less than)
- `>` (greater than)
- `<=` (less than or equal)
- `>=` (greater than or equal)

## String Comparisons

Take this component:
`./myapp/components/ColorPicker.html`
```html
<div props='color string'>
  { color == "red" ? (
    <span style="color:red;">Red</span>
  ) : (
    { color == "blue" ? (
      <span style="color:blue;">Blue</span>
    ) : (
      <span>Unknown color</span>
    ) }
  ) }
</div>
```

Used like:
```html
<ColorPicker color='red' />
```

Should produce:
```html
<div>
  <span style="color:red;">Red</span>
</div>
```

Used like:
```html
<ColorPicker color='green' />
```

Should produce:
```html
<div>
  <span>Unknown color</span>
</div>
```

## Nested Ternary (Multiple Conditions)

Take this component:
`./myapp/components/Grade.html`
```html
<div props='score int'>
  { score >= 90 ? (
    <p>A</p>
  ) : (
    { score >= 80 ? (
      <p>B</p>
    ) : (
      { score >= 70 ? (
        <p>C</p>
      ) : (
        { score >= 60 ? (
          <p>D</p>
        ) : (
          <p>F</p>
        ) }
      ) }
    ) }
  ) }
</div>
```

Used like:
```html
<Grade score={95} />
<Grade score={85} />
<Grade score={72} />
<Grade score={65} />
<Grade score={50} />
```

Should produce:
```html
<div>
  <p>A</p>
</div>
<div>
  <p>B</p>
</div>
<div>
  <p>C</p>
</div>
<div>
  <p>D</p>
</div>
<div>
  <p>F</p>
</div>
```

## Logical AND (`&&`)

Take this component:
`./myapp/components/AccessCheck.html`
```html
<div props='role string, active boolean'>
  { role == "admin" && active == true ? (
    <p>Full Access</p>
  ) : (
    <p>Limited Access</p>
  ) }
</div>
```

Used like:
```html
<AccessCheck role='admin' active={true} />
```

Should produce:
```html
<div>
  <p>Full Access</p>
</div>
```

Used like (one condition false):
```html
<AccessCheck role='admin' active={false} />
```

Should produce:
```html
<div>
  <p>Limited Access</p>
</div>
```

## Logical OR (`||`)

Take this component:
`./myapp/components/PriorityUser.html`
```html
<div props='role string'>
  { role == "admin" || role == "moderator" ? (
    <p>Priority User</p>
  ) : (
    <p>Standard User</p>
  ) }
</div>
```

Used like:
```html
<PriorityUser role='admin' />
<PriorityUser role='moderator' />
<PriorityUser role='guest' />
```

Should produce:
```html
<div>
  <p>Priority User</p>
</div>
<div>
  <p>Priority User</p>
</div>
<div>
  <p>Standard User</p>
</div>
```

## Not Equal (`!=`)

Take this component:
`./myapp/components/NotBanned.html`
```html
<div props='status string'>
  { status != "banned" ? (
    <p>Welcome back!</p>
  ) : (
    <p>Account suspended</p>
  ) }
</div>
```

Used like:
```html
<NotBanned status='active' />
```

Should produce:
```html
<div>
  <p>Welcome back!</p>
</div>
```

Used like:
```html
<NotBanned status='banned' />
```

Should produce:
```html
<div>
  <p>Account suspended</p>
</div>
```

## Ternary With Prop Expressions

Variables used in ternary conditions should also be accessible as regular props:

`./myapp/components/ScoreCard.html`
```html
<div props='score int'>
  <h2>Your score: {score}</h2>
  { score >= 50 ? (
    <p>You passed!</p>
  ) : (
    <p>You failed.</p>
  ) }
  <p>Final score: {score}</p>
</div>
```

Used like:
```html
<ScoreCard score={75} />
```

Should produce:
```html
<div>
  <h2>Your score: 75</h2>
  <p>You passed!</p>
  <p>Final score: 75</p>
</div>
```

## Deeply Nested Ternary

```html
<div props='level int, premium boolean'>
  { level >= 10 ? (
    { premium == true ? (
      <p>Premium high-level user</p>
    ) : (
      <p>High-level user</p>
    ) }
  ) : (
    <p>Regular user</p>
  ) }
</div>
```

Used like:
```html
<UserType level={15} premium={true} />
```

Should produce:
```html
<div>
  <p>Premium high-level user</p>
</div>
```

## Ternary Error Cases

### Malformed Ternary Syntax - Missing Colon
```html
<div props='x int'>
  { x == 1 ? (
    <p>One</p>
  ) }
</div>
```

Should error about incomplete ternary (missing `:` and falsy branch).

### Malformed Condition Syntax
```html
<div props='x int'>
  { x = 1 ? (
    <p>One</p>
  ) : (
    <p>Not one</p>
  ) }
</div>
```

Should error about invalid condition syntax (single `=` instead of `==`).

### Empty Condition
```html
<div>
  { ? (
    <p>Empty condition</p>
  ) : (
    <p>Other</p>
  ) }
</div>
```

Should error about empty condition.

### Undefined Variable In Condition
```html
<div>
  { undefinedVar == "test" ? (
    <p>Matched</p>
  ) : (
    <p>Not matched</p>
  ) }
</div>
```

Should error about undefined variable.

### Unbalanced Parentheses
```html
<div props='x int'>
  { x == 1 ? (
    <p>One</p>
  : (
    <p>Not one</p>
  ) }
</div>
```

Should error about unbalanced parentheses.
