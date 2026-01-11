# Component Props

## The `props` Attribute
Components may have a `prop` attirbute in their outermost tag which indiactes the dynamic data available within the component. 

The props attribute is formatted like so:
```html
<h1 props='PROP1NAME PROP1TYPE, PROP2NAME, PROP2TYPE...'</h1>
```

Each `NAME TYPE` pair is seperated by a comma.

## How To Render A `props`
Props are evaluated within expressions found within a component. Expressions are found within double-curly braces like so `{}`. More on expressions can be found below.
