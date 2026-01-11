# Prop Drilling

## Passing Props to Children Components
Imagine the following situation: we have a component which already accepts a `prop` but also has children components which accept `props` themselves.

Here is `./myapp/components/AnotherComponentAgain.html`
```html
<div props='text string'>
  <p>{text}</p>
</div>
```

Here is `./myapp/components/AnotherComponent.html`
```html
<div props='text string'>
  <p>{text}</p>
  <AnotherComponentAgain text={text} />
</div>
```

And now, we use that component in `./myapp/components/SomeComponent.html`
```html
<div props='title string'>
  <AnotherComponent text={title} />
</div>
```

See how we pass the parent prop `title` to the child component's `text` prop using expression syntax: `text={title}`.

This is how prop drilling works in `gtml` - you simply pass the parent's prop value using an expression to the child component's prop.

In this way, prop drilling enables us to pass values down the component tree.

Also take note, we can drill props down multiple layers.
