# Testing Slots

## Test Purpose
Validate that slots work correctly for content injection.

## Basic Slot Usage

Take this layout component:
`./myapp/components/PageLayout.html`
```html
<html>
  <body>
    <header>Site Header</header>
    {{ slot: content }}
    <footer>Site Footer</footer>
  </body>
</html>
```

If I use it like this:
```html
<PageLayout>
  <slot name='content' tag='main'>
    <p>Hello World</p>
  </slot>
</PageLayout>
```

I should get:
```html
<html>
  <body>
    <header>Site Header</header>
    <main>
      <p>Hello World</p>
    </main>
    <footer>Site Footer</footer>
  </body>
</html>
```

## Multiple Slots

Take this layout:
`./myapp/components/TwoColumnLayout.html`
```html
<div class='container'>
  {{ slot: sidebar }}
  {{ slot: main }}
</div>
```

Used like:
```html
<TwoColumnLayout>
  <slot name='sidebar' tag='aside'>
    <nav>Navigation</nav>
  </slot>
  <slot name='main' tag='section'>
    <p>Main content here</p>
  </slot>
</TwoColumnLayout>
```

Should produce:
```html
<div class='container'>
  <aside>
    <nav>Navigation</nav>
  </aside>
  <section>
    <p>Main content here</p>
  </section>
</div>
```

## Slot With Attributes

Slots should preserve additional attributes:
```html
<slot name='content' tag='div' class='my-class' id='my-id'>
  <p>Content</p>
</slot>
```

Should resolve to:
```html
<div class='my-class' id='my-id'>
  <p>Content</p>
</div>
```

Note: The `name` and `tag` attributes are consumed and removed from output.

## Slot With Props From Parent

A slot can contain components that use props from the parent:
```html
<PageLayout>
  <slot name='content' tag='main'>
    <UserCard name={currentUser} />
  </slot>
</PageLayout>
```

## Slot Error Cases

### Missing Slot Name
A slot without a name should error:
```html
<slot tag='div'>
  <p>Content</p>
</slot>
```

### Missing Tag Attribute
A slot without a tag attribute should error:
```html
<slot name='content'>
  <p>Content</p>
</slot>
```

### Slot Not Found in Parent
If a slot references a name that doesn't exist in the parent component:

Parent component:
```html
<div>
  {{ slot: sidebar }}
</div>
```

Usage (should error because 'content' doesn't exist in parent):
```html
<ParentComponent>
  <slot name='content' tag='div'>
    <p>Content</p>
  </slot>
</ParentComponent>
```

### Unfilled Slot
If a parent component has a slot placeholder but the child doesn't provide content:
```html
<PageLayout>
  <!-- no slot provided -->
</PageLayout>
```

Should either render empty or error. Test for consistent behavior.

### Duplicate Slot Names
If a child provides multiple slots with the same name:
```html
<PageLayout>
  <slot name='content' tag='div'>First</slot>
  <slot name='content' tag='div'>Second</slot>
</PageLayout>
```

Should error.
