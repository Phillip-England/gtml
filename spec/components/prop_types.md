# Component Props Types

## Strings
Strings a valid type and may be declared within a components prop like so:
`./myapp/components/AnotherComp.html`
```html
<div props='name string'>
  <p>{name}</p>
</div>
```

And then used like so:
```html
<AnotherComp name='Bob'>
```

Strings are the only type which can be passed as a raw string into the prop attribute. All other types must pass their value as an expression like so:
```html
<AnotherComp name={'Bob'}>
```

Now, we can do either with strings, I just want to make this distintion, this type of raw insertion is only legal with raw string types, an error will be produced if attempted with other types.

## Ints
You may declare an int within a component like so:
`./myapp/components/IntBoy.html`
```html
<div props='favoriteNumber int'>
  <p>My favorite number is {favoriteNumber}</p>
</div>
```

You may use the component like so:
```html
<IntBoy favoriteNumber={2}>
```

Which produces:
```html
<div props='favoriteNumber int'>
  <p>My favorite number is 2</p>
</div>
```

Take note, you must pass `int` values in as expressions which means they can be manipulated during compilation like so:

```html
<IntBoy favoriteNumber={2+2}>
```

However, exprssions are strict and will not allow mixing of different types with the following expression producing an error:
```html
<IntBoy favoriteNumber={2+"2"}>
```

## Booleans
Booleans are another primitive type we may use in our component's props. Here is an example:

`./myapp/components/IsAdmin.html`
```html
<div props='isAdmin boolean'>
  { isAdmin ? (
    <p>I am the admin</p>
  ) : (
    </p>I am not the admin</p>
  ) }
</div>
```

And that component may be used like so:
```html
<IsAdmin isAdmin={true}/>
```

Which will produce the following output:
```html
<div props='isAdmin boolean'>
  <p>I m the admin</p>
</div>
```

Take note, booleans are evaluate within an expression, which means we can do something like this:
```html
<IsAdmin isAdmin={2 > 1}/>
```

That would resolve to `true` as well.
