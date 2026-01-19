gtml scripts provide an easy way for us to gain access to elements. We may access elements directly by their css selectors in place like so:

```html
<div props='name string'>
  <button id='btn'>button</button>
  <p class='name'>bob</p>
  <p class='name'>bob</p>
  <p class='name'>bob</p>
  <p class='name'>bob</p>
</div>

<script type='gtml'>
  console.log(#btn) // compiles down into document.querySelector("#btn")
  console.log(.name) // compiles down into doucment.querySelector('.name')
  console.log(.name*) // compiles down into doucment.querySelectorAll('.name')
</script>
```

We make use of the `*` at the end of the css selector to dictate if we want to use querySelectorAll of just querySelector.
