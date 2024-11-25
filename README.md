# gtml
Html in Go doesn't *have* to suck. Introducing... 🥁 gtml.

## Hello, World

Turn this:
```html
<div _component='HelloWorld'>
    <h1>Hello, World!</h1>
</div>
```

Into this:
```go
func HelloWorld() string {
    var builder strings.Builder
    builder.WriteString(`<div _component='GreetingCard'><h1>Hello, World!</h1>`)
    return builder.String()
}
```

## Attribute Listing
gtml uses html attributes to determine the structure of our components. Here is a list of the available attributes:

### _component
A _component element is the root-level element for a gtml component. A _component element must have no parents and it must be named using PascalCasing. A _component element may not have any _component elements within it.

Here is a basic _component element making use of a prop (more on props later):
```html
<div _component='Greeting'>
    <h1>Hello, {{ name }}</h1>
</div>
```

gtml will scan all the `.html` files in a given directory, generating a Go function for each _component element. The above component will resolve to:
```go
func Greeting(name string) string {
    var builder strings.Builder
    builder.WriteString(`<div _component='Greeting'><h1>Hello, `)
    builder.WriteString(name)
    builder.WriteString(`!</h1>`)
    return builder.String()
}
```

### _for
A _for element allows us to iterate over a slice of type [T]. The value of the `_for` attribute must follow this schema:

```html
<div _for="ITEM of ITEMS []TYPE">...</div>
```

Here is a for element. Take note of how we access the underlying type's data using `{{}}` along with `item.Property` syntax:
```html
<div _component="GuestList">
    <p>{{ name }}</p>
    <ul _for='guest of guests []Guest'>
        <li>{{ guest.Name }}</li>
    </ul>
</div>
```

The output:
```go
func GuestList(name string, guests []Guest) string {
    var builder strings.Builder
    guestFor := gtml.For(guests, func(i int, guest Guest) string {
        var guestBuilder strings.Builder
        guestBuilder.WriteString(`<ul _for="guest of guests []Guest"><li>`)
        guestBuilder.WriteString(guest.Name)
        guestBuilder.WriteString(`</li></ul>`)
        return guestBuilder.String()
    })
    builder.WriteString(`<div _component="GuestList"><p>`)
    builder.WriteString(name)
    builder.WriteString(`</p>`)
    builder.WriteString(guestFor)
    builder.WriteString(`</div>`)
    return builder.String()
}
```

You can also choose to iterate over a `[]string` slice. Take note of how we access the strings value using `{{ .strname }}` syntax. Doing this will prevent `{{ color }}` from appearing in the output functions parameter definition.
```html
<div _component="FavoriteColors">
    <ul _for="color of colors []string">
        <li>{{ .color }}</li>
    </ul>
</div>
```

The output:
```go
func FavoriteColors(colors []string) string {
    var builder strings.Builder
    colorFor := gtml.For(colors, func(i int, color string) string {
        var colorBuilder strings.Builder
        colorBuilder.WriteString(`<ul _for="color of colors []string"><li>`)
        colorBuilder.WriteString(color)
        colorBuilder.WriteString(`</li></ul>`)
        return colorBuilder.String()
    })
    builder.WriteString(`<div _component="FavoriteColors">`)
    builder.WriteString(colorFor)
    builder.WriteString(`</div>`)
    return builder.String()
}
```

## _if & _else
_if and _else Elements are opposites. _if will allow a block of html to be rendered if a boolean is true. _else will allow a block to be rendered if a boolean is false.

Here an example:



# Features for v0.1.0
- For Loops ✅
- Conditionals ✅
- Slots ✅
- Internal Placeholders ✅
- Root Placeholders ✅
- Automatic Attribute Organizing ✅
- Basic Props ✅
- Command Line Tool For Generation
- A --watch Command
- Type Generation
- Solid README.md
- Managing Imports and Package Names in Output File
- Tests For Many Multilayered Components
- Attributes can use props ✅
- Varname Randomization (found instances where multiple uses of the same name will be needed)

# Rules Noted
- all HTML attribute names must be written in kebab casing while attribute values may be camel case
- when declaring a prop using {{ propName }} syntax, you must use camel casing to define the name
- use @ to pipe props into an child Elements.
