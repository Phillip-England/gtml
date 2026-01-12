# Using `slot`

## Inserting Content In Certain Places
Imagine we want to make a layout for a blog where each page has the same navbar, but we have a different post for each page? A `slot` could be used to make this trivial.

Look at this component:
`./myapp/components/GuestLayout.html`
```html
<html props='title string'>
  <head>
    <title>{title}</title>
  </head>
  <body>
    <nav>
      <ul>
      <a href='/'>Home</a>
      </ul>
    </nav>
    <slot name='post' />
  </body>
</html>
```

Then, in another component we can say:
`./myapp/components/BlogPost.html`
```html
<GuestLayout title="Some Page">
  <slot name='post' tag='article'>
    <h1>Some Post</h1>
    <p>Some Content</p>
  </slot>
</GuestLayout>
```

## How Slot Replacement Works
When we process a component like `./myapp/components/BlogPost.html`, we see the `<slot>` element. It is named post, so now we check to see if the parent component has a slot with a matching name. If it does, then we check for the tag attribute on the `slot`. In our final outputted html, the slot tag will be replaced with a tag with the name found in the tag attribute. For exmaple, this slot:

```html
<slot name='main-section' tag='article' class='myclass'>
  <h1>Some Article</h1>
</slot>
```

would resolve to:
```html
<article name='main-section' tag='article' class='myclass'>
  <h1>Some Article</h1>
</article>
```
