# Expressions Within Components

## What Does `{}` Mean?
`{}` indicates we are evaluating an expression within a component. This `{}` may be found within an elements text area as well as within its attributes. Below we will discuss both types of scenarios, but one thing needs to be made very clear up front. `{}` is an indicator that our compiler needs to evaluate some sort of expression within a component. Ultimately, after compilation, all `{}` signs will be replaced with the value they resolve to.

## The Use of `{}` For Wrapping Props In Text Areas
In our components, we may use `{}` within the text area of a component. For example, let's take: `<p>{age}</p>`. When the compiler locates the word `age` within the `{}` it then seeks to replace the whole `{age}` construct with the actual value of age passed into the component.

For example, take this component `ThatComponent`:
```html
<div props="name string, age int">
  <p>My name is {name}</p>
  <p>My age is {age}</p>
</div>
```

If we were to use that component like so:
```html
<ThatComponet name='Bob' age={2}>
```

Then we would get the following output:
```html
<div>
  <p>My name is Bob</p>
  <p>My age is 2</p>
</div>
```

## Raw Values Versus Evaluated Values
Take note here:
```html
<ThatComponet name='Bob' age={2}>
```

See how `name='Bob'` is used instead of `name={"Bob"}`? Well, either are acceptable, but if raw values which do not need to be evaluated are passed in, `name='Bob'` is preffered.

However, we could have just as easiy said, `name={"B"+"o"+"b"}` and the same result would be produced. It would just require evaluation to produce such result, which is more work and should thus be avoid if at all possible.

## Strings Versus All Other Types
Take note here:
```html
<ThatComponet name='Bob' age={2}>
```

See how we say `age={2}` and not `age="2"`? Well, `int` values must ALWAYS be evaluated. This is because any value which is not a string must be evaluated. The only time passing a value into a component as a raw string is acceptible is if the value itself is actually a raw string. Plain and simple.

## Exmaples of Using `{}` Expression in Attributes
We may come across situations where expressions are being utilized within a components props as indicated by the `{}` sign:
```html
<ThatComponet name='Bob' age={2+2}>
```

In which case the expression would evalutate to 4, resulting in the following output:
```html
<div>
  <p>My name is Bob</p>
  <p>My age is 4</p>
</div>
```

## You Can Evaulate Strings in `{}` as Well
What if we wanted to evalaute an expression using a string?
```html
<ThatComponent name={"Bo"+"b"} age={2} />
```

In the above example we say `name={"Bo"+"b"}` which will evaluate to "Bob". Strings are also able to be evaluated within `{}`.

## More on `{}` Within Attributes
I mentioned briefly in the previous section that `{}` may also be used within an html elements attribute to say, "this is an expression which needs evaluated." However, I would like to explore futher exactly what this means for clarification.

Take the following component which takes in an age:
`./myapp/components/AgeComp.html`
```html
<div props='age int'>
  <p>I am {age} years old</p>
</div>
```

Then take this component which uses the above component within itself. Take note this component itself has a prop `extraYears` which is an `int`:
```html
<div props='extraYears int'> 
  <AgeComp age={extraYears + 2} />
</div>
```

See how we make use of `{}` on the right-handed side of the `age` attribute? Then within, we access the `extraYears` prop which is accessible to components deeper within the tree. Then, we add `2` to the value of `extraYears`.

In this way, users may perform operations on props within the template itself using expressions.

These expressions must be able to be elaborate like so:
```html
<div props='extraYears int'> 
  <AgeComp age={extraYears + 2 / (2 * 10) % 4} />
</div>
```

And we would expect this expression to be evaluated in standard PEMDAS order. 

This gives users a flexible way to pipe values to child components while being able to perform operations on values as well.

## Quirks When Evaluating expressions
Evaluating Expressions within components may get messy, so we need some structure to how things ought to play out to ensure we keep things clean and tidy.

Here are some examples of things that I mean.

What happens if we get an expression like so: `{someStr+2}` where a user tries to add an int to a string?

Well, we need to be very strict about these types of interactions and we need to implement multiple tests to ensure these sorts of edge cases are accounted for.
