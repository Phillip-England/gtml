# Testing Prop Drilling

## Test Purpose
Validate that props can be passed from parent components to child components.

## Basic Prop Drilling

Take this component:
`./myapp/components/InnerComponent.html`
```html
<div props='text string'>
  <p>{text}</p>
</div>
```

And this layout:
`./myapp/components/OuterComponent.html`
```html
<div props='heading string'>
  <h1>{heading}</h1>
  <InnerComponent text={heading} />
</div>
```

Used like:
```html
<OuterComponent heading='some heading' />
```

Should produce:
```html
<div>
  <h1>some heading</h1>
  <div>
    <p>some heading</p>
  </div>
</div>
```

## Multi-Level Prop Drilling

Take three components:
`./myapp/components/LevelThree.html`
```html
<span props='value string'>{value}</span>
```

`./myapp/components/LevelTwo.html`
```html
<div props='passthrough string'>
  <LevelThree value={passthrough} />
</div>
```

`./myapp/components/LevelOne.html`
```html
<section props='data string'>
  <LevelTwo passthrough={data} />
</section>
```

Used like:
```html
<LevelOne data='deep value' />
```

Should produce:
```html
<section>
  <div>
    <span>deep value</span>
  </div>
</section>
```

## Drilling With Transformation

Take this component:
```html
<div props='base int'>
  <ChildComp value={base * 2} />
</div>
```

The child receives the transformed value (doubled).

## Drilling Non-Existent Prop
If a parent tries to drill a prop that doesn't exist:
```html
<div props='name string'>
  <ChildComp text={nonexistent} />
</div>
```

Should error indicating `nonexistent` is not defined.
