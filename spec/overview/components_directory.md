# The Components Directory

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
