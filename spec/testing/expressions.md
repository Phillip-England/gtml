# Testing Expressions

## Test Purpose
Validate that expressions within `{}` are evaluated correctly in all contexts.

## Expressions In Text Content

### Simple Prop Reference
```html
<p props='name string'>Hello, {name}!</p>
```

Used as:
```html
<Greeting name='World' />
```

Produces:
```html
<p>Hello, World!</p>
```

### Multiple Expressions In Same Element
```html
<p props='first string, last string'>{first} {last}</p>
```

Used as:
```html
<FullName first='John' last='Doe' />
```

Produces:
```html
<p>John Doe</p>
```

## Expressions In Attributes

### Prop In Attribute Value
```html
<a props='url string, label string' href={url}>{label}</a>
```

Used as:
```html
<Link url='https://example.com' label='Click here' />
```

Produces:
```html
<a href="https://example.com">Click here</a>
```

### Prop In Class Attribute
```html
<div props='className string' class={className}>Content</div>
```

Used as:
```html
<StyledDiv className='my-class' />
```

Produces:
```html
<div class="my-class">Content</div>
```

## Arithmetic Expressions

### Addition
```html
<span props='a int, b int'>{a + b}</span>
```

Used as:
```html
<Math a={5} b={3} />
```

Produces:
```html
<span>8</span>
```

### Subtraction
```html
<Math a={10} b={4} />
```

Produces:
```html
<span>6</span>
```

### Multiplication
```html
<span props='a int, b int'>{a * b}</span>
```

Used as:
```html
<Multiply a={6} b={7} />
```

Produces:
```html
<span>42</span>
```

### Division
```html
<span props='a int, b int'>{a / b}</span>
```

Used as:
```html
<Divide a={20} b={4} />
```

Produces:
```html
<span>5</span>
```

### Modulo
```html
<span props='a int, b int'>{a % b}</span>
```

Used as:
```html
<Modulo a={17} b={5} />
```

Produces:
```html
<span>2</span>
```

### Complex Expression With Parentheses
```html
<span props='x int'>{x + 2 / (2 * 10) % 4}</span>
```

Used as:
```html
<Complex x={10} />
```

Should evaluate following PEMDAS order.

### Division By Zero
```html
<Divide a={10} b={0} />
```

Should either error or produce a defined behavior (Infinity, error message, etc.). Test for consistent behavior.

## String Concatenation

### Basic Concatenation
```html
<p props='first string, last string'>{first + " " + last}</p>
```

Used as:
```html
<FullName first='John' last='Doe' />
```

Produces:
```html
<p>John Doe</p>
```

### Multiple Concatenations
```html
<p props='a string, b string, c string'>{a + b + c}</p>
```

Used as:
```html
<Concat a='Hello' b=' ' c='World' />
```

Produces:
```html
<p>Hello World</p>
```

## Ternary Expressions

### Basic Ternary
```html
<span props='isActive boolean'>{ isActive ? "Active" : "Inactive" }</span>
```

Used as:
```html
<Status isActive={true} />
```

Produces:
```html
<span>Active</span>
```

### Ternary With HTML
```html
<div props='showDetails boolean'>
  { showDetails ? (
    <p>Here are the details...</p>
  ) : (
    <p>Click to see more</p>
  ) }
</div>
```

Used as:
```html
<Details showDetails={true} />
```

Produces:
```html
<div>
  <p>Here are the details...</p>
</div>
```

### Nested Ternary
```html
<span props='value int'>
  { value > 0 ? "Positive" : (value < 0 ? "Negative" : "Zero") }
</span>
```

Used as:
```html
<Sign value={5} />
<Sign value={-3} />
<Sign value={0} />
```

Produces:
```html
<span>Positive</span>
<span>Negative</span>
<span>Zero</span>
```

## Expression Type Errors

### Adding Int To String
```html
<span props='num int, text string'>{num + text}</span>
```

Should error due to type mismatch.

### Comparing Different Types
```html
<span props='num int, text string'>{ num == text ? "equal" : "not equal" }</span>
```

Should error or have defined comparison behavior.

### Undefined Variable In Expression
```html
<span>{undefinedVar}</span>
```

Should error since `undefinedVar` is not declared in props.
