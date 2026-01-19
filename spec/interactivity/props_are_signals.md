When we gtml detects props within a gtml script, it will convert them into signals in place.

Here is what I mean:

```html
<div props='name string'>
  <p>My name is {name}!</p>
  <button id='btn'>Click Me!</button>
</div>

<script type='gtml'>
  #btn.onclick(() => {
    $name = "John"
  })
</script>
```

Do you see how the above component has the "name" prop? Well, because we find the `$name` within our gtml script, we know that the `$name` variable is actually a signal. This tells our compiler to go ahead and produce a signal for this value.

Then, we can apply operations to the signal directly like in th example above we simply set the value of the signal to "John"

This will compile down into the actual javascript necessary to make these changes, and we will make use of our signal class to do so.
