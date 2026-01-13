# Client Side Fetching

## The `fetch` Attribute
An element with the `fetch` attribute will make a client-side fetch request AFTER the page has loaded. The value found within the `fetch` attribute must be string formatted like so:
```html
<div fetch='GET localhost:8000/api/users'></div>
```

## `fetch` Attribute Formatting
See how the string is formatted `METHOD URL`? That formatting is required. I considered using the `method` attribute, but opted out because I did not want to create potential conflicts with html `<form>` elements which make use of the `method` attribute.

## Data Expectations
gtml expects data received from fetch requests to be json. If the received data is not json, then we will resort to displaying the `<fallback>` element instead. More on `fallbacks` in a bit.

## Naming Incoming Data via `as`
An element with the `fetch` attribute may also have an `as` attribute. The value found within the `as` attribute will be the "name" of the incoming body data found within the fetch request initiated by the element. This is required if we intent to iterate over the incoming data and display it as html.

## JSON Data Iteration
If json data is found within the body of the incoming fetch request, then that data may be iterated through. However, this may only be done if we provide the json data a name via the `as` attribute. Here is a common example:

`IterComp`
```html
<div fetch='GET localhost:8080/api/users' as='users'>
  <ul>
    <li for='user in users'>
      <p>{user.name}</p>
      <ul>
        <li for='color in user.colors'><p>{color.name}</p></li>
      </ul>
    </li>
  </ul>
</div>
```

## Iteration Compilation
The above code will compile in the following manner. The compilers goal is to initiate client-side behavior, so at the end of the day, javascript is our compiler's output target for this bit of compilation work.

Here is the breakdown, we are referring to `IterComp` above.

First, gtml sees we are making a GET request to `localhost:8080/api/users`. so it drafts up the javascript code to make such a request. It also notes the element in which this request is derived from, because this element will need to be targeted via javascript in our output.

Second, gtml sees we named the incoming JSON data `users` via `as="users"`. The word `users` now becomes associated with the scope created while compiling this `fetch` element.

Third, gtml sees we have a `<li>` with a `for` attribute which points to our users via `for="user in users"`. Lets break this down very clearly because a lot happens in this moment. The element with the `for` attribute was found within a `fetch` element which signals it might be a potential iterative block. So, we check the value in the `for` element. We check to see if the formatting is good. `ITEM in ITEMS` where `ITEMS` must be the name of the JSON data captured by the `fetch` element. If `ITEMS` is in fact an actual piece of data we named via `as`, then we have identified an iterative block. Once we identify an iterative block, we must get the name of the iteration item which is the `ITEM`. So, in `for="user in users"`, `users` is the `ITEMS` and `user` is the `ITEM`.

Fourth, now that we know each item in our `users` is called `user`, we can duplicate the `for` element for each user. So, imagine we pulled in 20 users via our `fetch` request, well this means the `for` element will be duplicated 20 times. One time for each user.

Fifth, `for` blocks must be evaluated. Since each `for` block has an `ITEM` associated with it, we may use that item to access the JSON data's properties. For example, in `IterComp` we access the user's name like so:

```html
<!--snip-->
<li for='user in users'>
  <p>{user.name}</p>
<!--...-->
```

`{user.name}` will be replaced with the actual user's name found within the json data.

BUT, what about nested iteration? What if we pull in some json data which itself has nested data which may need to be iterated through? Well, an `ITEM` may itself be utilized in a nested `for` element like so:

```html
<!--snip-->
<li for='user in users'>
  <p>{user.name}</p>
  <ul>
    <li for='color in user.colors'><p>{color.name}</p></li>
  </ul>
<!--...-->
```

Do you see how we say `for='color in user.colors'`? Well, since we already have access to `user`, we can then point to its nested data within the `for` element in which it is contained.

This provides us an easy way to handle iteration and nested iteration by depending on these naming schemes which are scoped and compiled into native javascript.

## Use <script> Tags
For now, make sure that we do not compiled to an external javascript file. For the time being, just insert raw script tags to manage things and make sure the script tags are isolated to the component the are associated with. I do not want script tags conflicting with each other. In the same way that styles are clearly isolated on a per-component basis, so may the script tags be for iteration.

## Suspense
Okay, fetch requests take time, and we do not want to just render a black screen during the flight of the requests initiated by our `fetch` elements. That is where the `suspense` attribute comes in handy. For example:

```html
<div fetch='GET localhost:8080/api/users' as='users'>
  <div suspense>
    <p>loading</p>
  </div>
  <ul>
  <li for='user in users'>{user.name}</li>
  </ul>
</div>
```

In the above example, the element with the `suspense` attribute will be displayed during the flight of the request, and once the request is resolved, it will be removed from the users view and perception.


## fallback
Okay, sometimes fetch requests fail, so we need a way to say, "show this html in the event of request failure". That is where the `fallback` attribute comes in hand. It is sort of the opposite of the `suspense` attribute.

Take this example:
```html
<div fetch='GET localhost:8080/api/users' as='users'>
  <div suspense>
    <p>loading</p>
  </div>
  <div fallback>
    <p>error loading users</p>
  </div>
  <ul>
  <li for='user in users'>{user.name}</li>
  </ul>
</div>
```

In the above example, the `fallback` element will only be displayed if the request to `localhost:8080/api/users` does not result in a successful request.
