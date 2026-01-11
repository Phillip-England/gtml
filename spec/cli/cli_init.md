# Command Line Init
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
<button props='text string'>{text}</button>
```

`./myapp/components/GuestLayout.html`
```html
<html props='title string'>
  <head>
    <title>{title}</title>
  </head>
  <body>
  <BasicButton text={title} />
    <slot name='content'/>
  </body>
</html>
```

`./myapp/routes/index.html`
```html
<GuestLayout title="Some Title">
  <slot name='content'>
    <p>Some Content</p>
  </slot>
  <BasicButton text={title} />
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
