gtml should offer a way to add gtml scripts directly within events within our dom.

For example:

```html
<div props='name string'>
  <p>My name is {name}!</p>
  <button onclick={() => {
    $name = "John"
  }}>Click Me!</button>
</div>
```

If we find any event like that ^ we will treat it as gtml syntax as well. This allows users to add events directly to elements as well as adding them into gtml scripts.

these events are found within an expression like `onclick={}` we open up an expression to indicate that this is an event which needs compiled as a gtml script.
