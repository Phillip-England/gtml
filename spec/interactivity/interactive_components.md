Components are interactive. This is the true magic behind gtml. It gives us a way to abstract DOM manipulations into our html with the use of attributes and signals.

Here is the magic, every prop is a signal which is available to be modified in place.

For example:

```html
<div props='name string'>
  <p>My name is {name}!</p>
  <button id='btn'>Click Me!</button>
</div>

<script type='gtml'>
  #btn.onclick(() => {
    name = "John"
  })
</script>
```

That is just a general overview, I will break things down in others files and explain more how it works in depth.
