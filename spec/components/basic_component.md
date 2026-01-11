# Basic Component Example

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

<div props='heading string, subheading string'>
  <h1>{heading}</h1>
  <p>{subheading}</p>
<div>
```
