# Testing Props Usage

## Test Purpose
Validate that props are correctly used and referenced within components.

## Basic Prop Usage

Take this component:
`./myapp/components/ThatButton.html`
```html
<button props='title string'>{title}</button>
```

If I use it like this:
```html
<ThatButton title='some title' />
```

I should get:
```html
<button>some title</button>
```

## Prop Used Multiple Times

Take this component:
```html
<div props='name string'>
  <h1>{name}</h1>
  <p>Welcome, {name}!</p>
  <footer>Goodbye, {name}</footer>
</div>
```

Used like:
```html
<WelcomeCard name='Alice' />
```

Should produce:
```html
<div>
  <h1>Alice</h1>
  <p>Welcome, Alice!</p>
  <footer>Goodbye, Alice</footer>
</div>
```

## Prop In Nested Elements

Take this component:
```html
<div props='title string'>
  <header>
    <nav>
      <h1>{title}</h1>
    </nav>
  </header>
</div>
```

Used like:
```html
<NestedTitle title='My Site' />
```

Should produce:
```html
<div>
  <header>
    <nav>
      <h1>My Site</h1>
    </nav>
  </header>
</div>
```

## Missing Required Prop
If a component expects a prop but it's not provided:
```html
<ThatButton />
```

Should error with a message indicating the missing prop `title`.

## Extra Props Passed
If extra props are passed that aren't declared:
```html
<ThatButton title='Click' extra='ignored' />
```

Should either ignore the extra prop or error. Test for consistent behavior.
