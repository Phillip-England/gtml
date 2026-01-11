# Testing Props Attribute Syntax

## Test Purpose
Validate that the `props` attribute is parsed correctly and props are accessible within components.

## Basic Single Prop

Take this component:
`./myapp/components/SimpleButton.html`
```html
<button props='text string'>{text}</button>
```

Used like:
```html
<SimpleButton text='Click me' />
```

Should produce:
```html
<button>Click me</button>
```

## Multiple Props

Take this component:
`./myapp/components/UserGreeting.html`
```html
<div props='name string, age int'>
  <p>Hello, {name}!</p>
  <p>You are {age} years old.</p>
</div>
```

Used like:
```html
<UserGreeting name='Alice' age={25} />
```

Should produce:
```html
<div>
  <p>Hello, Alice!</p>
  <p>You are 25 years old.</p>
</div>
```

## Props Attribute Removed From Output

The `props` attribute should be removed from the compiled output. Take:
```html
<div props='title string'>
  <h1>{title}</h1>
</div>
```

Used like:
```html
<MyComponent title='Hello' />
```

Should produce:
```html
<div>
  <h1>Hello</h1>
</div>
```

Note: `props='title string'` is NOT in the output.

## Props With Extra Whitespace

The parser should handle extra whitespace in props:
```html
<div props='  name   string  ,   age   int  '>
  {name} is {age}
</div>
```

Should parse correctly as two props: `name` (string) and `age` (int).

## Props Error Cases

### Invalid Prop Type
A prop with an unrecognized type should error:
```html
<div props='data unknown'>
  {data}
</div>
```

### Missing Prop Type
A prop without a type should error:
```html
<div props='name'>
  {name}
</div>
```

### Duplicate Prop Names
Duplicate prop names should error:
```html
<div props='name string, name int'>
  {name}
</div>
```

### Props On Non-Root Element
If `props` appears on a nested element, it should be treated as a regular attribute (not parsed as component props, unless of course it is passed to a component):
```html
<div>
  <span props='some value'>{text}</span>
</div>
```
