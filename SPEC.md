
=====================================
# What is `gtml`
=====================================

## `gtml` is a Static Site Generator
`gtml` is an opinionated static site generator which compiles components and routes into static html assets.

## `gtml` is a Templating System
`gtml` is centered around html generation. The goal is to make reusable bits of html which you can reuse throughout your entire webpage.

## `gtml` is a Compiler
`gtml` takes in our routes, resolves the components found within, and compiles everything into static html.

=====================================
# Project Structure
=====================================

## Basic Project Layout
`gtml` expects a very specific project structure. It looks like this:
```bash
./myapp
./myapp/components
./myapp/routes
./myapp/dist
./myapp/static
```

These directories form the foundation of a `gtml` application and are required. Without these directories, `gtml` will fail to compile.

=====================================
# Command Line Usage
=====================================

## `gtml init <PATH>`
You may initalize a `gtml` project by calling `gtml init <PATH>`. For example, I may say `gtml init myapp`. This would generate the `./myapp` directory if it does not exist. However, if `./myapp` does exist, `gtml` will fail. To force initalization you may pass the `--force` flag like so: `gtml init myapp --force`.

Upon initalization using `gtml init myapp`, the following directory structure will be generated:
```bash
./myapp
./myapp/components
./myapp/components/BasicButton.html
./myapp/components/GuestLayout.html
./myapp/routes
./myapp/routes/index.html
./myapp//static
./myapp/dist
./myapp/dist/index.html
./myapp/dist/static
```

The html files within the initalized project will look as follows:

`./myapp/components/BasicButton.html`
```html
<button>{{ prop: text string }}</button>
```

`./myapp/components/GuestLayout.html`
```html
<html>
  <head>
    <title>{{ prop: title string }}</title>
  </head>
  <body>
  <BasicButton text='{{ drill: title }}' />
    {{ slot: content }}
  </body>
</html>
```

`./myapp/routes/index.html`
```html
<GuestLayout title="Some Title">
  <p>Some Content</p>
  <BasicButton text='{{ drill: title }}' />
</GuestLayout>
```

`./myapp/dist/index.html`
```html
<html>
  <head>
    <title>Some Title</title>
  </head>
  <body>
    <button>Some Title</button>
    <p>Some Content</p>
    <button>Some Title</button>
  </body>
</html>
```

## `gtml compile ./somedir`
The `compile` command will attempt to compile the routes found at `./somedir/routes`. This command will check to ensure our project structure is correct and that everything checks out. Upon failure, this command will let you know exactly why things failed. If things are successful, you should have your static `html` in `./somedir/dist`

If you pass the `--watch` flag, changes to the any file within the `./somedir` directory will trigger recompilation.

=====================================
# The Components Directory
=====================================

## Subdirectories Are Allowed
The `./myapp/components` directory may have subdirectories like so: `./myapp/components/buttons`. That is permitted.

## Unique Names Required
Each `html` file within the `./myapp/components` directory must be unique, even if two files with the same name are in different directories. For example, lets imagine we have `./myapp/components/SomeComponent.html` and then `./myapp/components/inner/SomeComponent.html`, that would result in an error because `SomeComponent.html` may only be found once within the entire components directories, including all subdirectories.

## PascalCase Names Required
All components in the `./myapp/components` directory must have PascalCase names. For example, if we find `./myapp/components/somecomponent.html`, that would result in an error. Instead, you would do: `./myapp/components/SomeComponent.html`.

## One Component Per File
Each `html` file within the `./myapp/components` directory must contain a single element. For example, this is an example of an invalid component:

```html
<p>One</p>
<p>Two</p>
```

Instead, you would do:

```html
<div>
  <p>One</p>
  <p>Two</p>
</div>
```

=====================================
# The Routes Directory
=====================================

## Lowercase `kebab-case` Names only 
The `./myapp/routes` directory must only contain `html` files with `kebab-casing` only. For example, `./myapp/routes/SomeRoute.html` would result in an error. Instead you would do: `./myapp/routes/some-route.html`.

## Repeated Names Are Fine
Unlike the `./myapp/components` directory, the `./myapp/routes` directory is allowed to have repeated names, as long as they occur in different subdirectories.

## Subdirectories Are Allowed
The `./myapp/routes` directory may have subdirectories like so: `./myapp/routes/blog`. That is permitted.

## Compiling To `./myapp/dist`
The `./myapp/routes` directory is compiled into static `html` and is then copied over to the `./myapp/dist` directory with the exact same naming. For example, when we compile `./myapp/routes/index.html`, the components within `./myapp/index.html` will be resolved from the `./myapp/components` directory and then we copy the fully rendered `html` into the `./myapp/dist` directory.

=====================================
# Compilation Workflow
=====================================

## Reading The `./myapp/components` Directory
First, we read the `./myapp/components` directory and create a map of all of our components. This gives us the names of all of our components and what they look like. We must also have acces to the components styles and other parameters down the line so make this data type something we can access to get information about our components anywhere during the compilation process.

## Reading The `./myapp/routes` Directory
Second, we read the `./myapp/routes` directory and we scan each file for the existence of components. When we find a component, we compile the component with the provided parameters. This must be done in a recursive manner because components themselves may contain other components within themselves. Once we compile the component down all the way into pure html, we replace the component in the route with the fully compiled component. This process is then repeated for all the components within the route until no components are left, resulting in pure html left.

## Copying The Routes to The `./myapp/dist` Directory
Third, the `html` which is derived from the routes is copied over into the `./myapp/dist` directory. When everything is said and done, the `./myapp/dist` directory should contain all of the fully compiled, static `html`. This `html` can then be served.

## Copying The `./myapp/static` Directory into `./myapp/dist/static`
Fourth, we copy over our static assets into the `./myapp/dist` directory. This will ensure the fully static `html` has access to the static assets for the site.

## Generating a Final `./myapp/dist/static/styles.css`
Each component may have a `<style></style>` section at the top. These styles are generated in such a way that they do not conflict between different components. For example, one component may say `<style>p {background:red;}</style>` and another one may say `<style>p {background:blue;}</style>` and when everything is said and done, these styles should not conflict with each other. These styles are all gathered and compiled into a single `css` file at `./myapp/dist/static/styles.css`. This enables users to style their components within a single file without conflicting with the styles of other components.

=====================================
# Components
=====================================

## Where Are Components Found?
Components exist within the `./myapp/components` directory and contain a single `html` component.

## A Basic Component
Here is an example of a basic component:
```html
<style>
  h1 {
    color:red;
  }
</style>

<div>
  <h1>{{ prop: heading string }}</h1>
  <p>{{ prop: subheading string  }}</p>
<div>
```

Take note of the `{{ prop: ... }}` syntax. This is called a `prop` and more information about how they work can below.


=====================================
# Using `props`
=====================================

## Found Within Components
`props` are found within components and are indicated by their use of double brackets in conjunction with the `prop` keywork like so: `{{ prop: }}`. 

## Name and Type Required
All `props` require a `name` and a `type`.

## Available Types
`string` and `int` are the only available types.

## Full example
Here is a full example of a component which makes use of a `prop` with the name being `text` and the type being `string`:

```html
<button>{{ prop: text string }}</button>
```

=====================================
# Using `drill`
=====================================

## Passing Props to Children Components
Imagine the following situation: we have a component which already accepts a `prop` but also has children components which accept `props` themselves.

Here is `./myapp/components/AnotherComponentAgain.html`
```html
<div>
  <p>{{ prop: text string }}</p>
</div>
```

Here is `./myapp/components/AnotherComponent.html`
```html
<div>
  <p>{{ prop: text string }}</p>
  <AnotherComponentAgain text="{{ drill: text }}" />
</div>
```

And now, we use that component in `./myapp/components/SomeComponent.html`
```html
<div>
  {{ prop: title string }}
  <AnotherComponent text="{{ drill: title }}" />
</div>
```

See how we make use of `{{ drill: ... }}` where `...` is the name of the `prop` you would like to `drill`.

In the above component where we made use `AnotherComponent` we said, `text="{{drill: title }}"`. That is our way of saying, "Whatever 'title' is equal to, well we would like to also use that for the value passed to 'text'"

In this way, `drill` enables us to pass values down the component tree.

Also take note, we can `drill` down multiple layers.

=====================================
# Using `slot`
=====================================

## Inserting Content In Certain Places
Imagine we want to make a layout for a blog where each page has the same navbar, but we have a different post for each page? A `slot` could be used to make this trivial.

Look at this component:
`./myapp/components/GuestLayout.html`
```html
<html>
  <head>
    <title>{{ prop: title string }}</title>
  </head>
  <body>
    <nav>
      <ul>
      <a href='/'>Home</a>
      </ul>
    </nav>
    {{ slot: post }}
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

=====================================
# Testing Props
=====================================

## Test Purpose
We should create numerous tests to ensure `props` are piped into components correctly. We also need to make sure that things fail as expected too.

## Example Props Test

Take this component:
```html
<button>{{ prop: title string }}</button>
```

If I use it like this:
```html
<ThatButton title='some title'>
```

I should get:
```html
<button>some title</button>
```

I need multiple tests like that. I also need tests which check for edge cases and failure like this:

```html
<div>{{ prop: title }}</div>
```

See how that doesn't include a type? Well, that should error and we need to test for these things.

Also, what happens if  we have two props with the same name but different types like this:

```html
<div>{{ prop: title string }}</div>
<div>{{ prop: title int }}</div>
```

Well, that should error out as well and is not allowed when using gtml.

=====================================
# Testing Prop Drilling
=====================================

## We should test for prop drilling to ensure it works as expected. For example, take this component named `SomeComponent`:

```html
<div>
  <p>{{ prop: text string }}</p>
</div>
```

And then this layout, `SomeLayout`:

```html
<div>
  <h1>{{ prop: heading string }}</h1>
  <SomeComponent text="{{ drill: heading }}" />
</div>
```

And then I compiled:
```html
<SomeLayout heading='some heading'/>
```

Well, that should produce:
```html
<div>
  <h1>some heading</h1>
  <div>
    <p>some heading</p>
  </div>
</div>
```

We need multiple tests to ensure prop drilling is working as expected. However, what if we try to drill a `string` into an `int` or vice-versa, well in that case we should expect an error. Check for that and multiple other types of errors.

=====================================
# Testing Slots
=====================================

## Test Purpose
We should create numerous tests to ensure `slots` are inserted into components correctly. We also need to make sure that things fail as expected too.

## Example Slot Test

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
    <main name='content' tag='main'>
      <p>Hello World</p>
    </main>
    <footer>Site Footer</footer>
  </body>
</html>
```

## Testing Multiple Slots

A component may have multiple slots. Take this layout:
`./myapp/components/TwoColumnLayout.html`
```html
<div class='container'>
  {{ slot: sidebar }}
  {{ slot: main }}
</div>
```

If I use it like this:
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

I should get:
```html
<div class='container'>
  <aside name='sidebar' tag='aside'>
    <nav>Navigation</nav>
  </aside>
  <section name='main' tag='section'>
    <p>Main content here</p>
  </section>
</div>
```

## Testing Slot Attributes

Slots should preserve additional attributes like `class`, `id`, etc. Take this:
```html
<slot name='content' tag='div' class='my-class' id='my-id'>
  <p>Content</p>
</slot>
```

Should resolve to:
```html
<div name='content' tag='div' class='my-class' id='my-id'>
  <p>Content</p>
</div>
```

## Testing Slot Error Cases

We need tests for error conditions:

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
If a slot references a name that doesn't exist in the parent component, it should error:

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

### Unfilled Required Slot
If a parent component has a slot placeholder but the child doesn't provide content for it, we should decide if this errors or renders empty. Test for consistent behavior.

### Duplicate Slot Names
If a child provides multiple slots with the same name, it should error:
```html
<PageLayout>
  <slot name='content' tag='div'>First</slot>
  <slot name='content' tag='div'>Second</slot>
</PageLayout>
```

=====================================
