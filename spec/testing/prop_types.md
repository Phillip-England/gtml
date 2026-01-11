# Testing Prop Types

## Test Purpose
Validate that all prop types work correctly and type checking is enforced.

## String Props

### Basic String Usage
```html
<p props='message string'>{message}</p>
```

Used as:
```html
<MyComp message='Hello World' />
```

Produces:
```html
<p>Hello World</p>
```

### String With Quotes
```html
<MyComp message="It's a nice day" />
```

Should handle apostrophes and special characters correctly.

### String With Expression Syntax
```html
<MyComp message={'Hello'} />
```

Should work the same as raw string.

### Empty String
```html
<MyComp message='' />
```

Should produce empty content where the prop is used.

## Int Props

### Basic Int Usage
```html
<span props='count int'>Count: {count}</span>
```

Used as:
```html
<Counter count={42} />
```

Produces:
```html
<span>Count: 42</span>
```

### Int With Arithmetic
```html
<Counter count={40 + 2} />
```

Should produce:
```html
<span>Count: 42</span>
```

### Negative Int
```html
<Counter count={-5} />
```

Should produce:
```html
<span>Count: -5</span>
```

### Zero
```html
<Counter count={0} />
```

Should produce:
```html
<span>Count: 0</span>
```

### Int Type Error - String Passed
Passing a raw string to an int prop should error:
```html
<Counter count='42' />
```

This should fail because `'42'` is a string, not `{42}` which is an int expression.

## Boolean Props

### True Value
```html
<div props='visible boolean'>
  { visible ? (
    <p>Visible</p>
  ) : (
    <p>Hidden</p>
  ) }
</div>
```

Used as:
```html
<VisibilityToggle visible={true} />
```

Produces:
```html
<div>
  <p>Visible</p>
</div>
```

### False Value
```html
<VisibilityToggle visible={false} />
```

Produces:
```html
<div>
  <p>Hidden</p>
</div>
```

### Boolean From Comparison
```html
<VisibilityToggle visible={1 > 0} />
```

Should evaluate to `true`.

### Boolean Type Error - String Passed
Passing a raw string to a boolean prop should error:
```html
<VisibilityToggle visible='true' />
```
