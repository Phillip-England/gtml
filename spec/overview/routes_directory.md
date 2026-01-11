# The Routes Directory

## Lowercase `kebab-case` Names only 
The `./myapp/routes` directory must only contain `html` files with `kebab-casing` only. For example, `./myapp/routes/SomeRoute.html` would result in an error. Instead you would do: `./myapp/routes/some-route.html`.

## Repeated Names Are Fine
Unlike the `./myapp/components` directory, the `./myapp/routes` directory is allowed to have repeated names, as long as they occur in different subdirectories.

## Subdirectories Are Allowed
The `./myapp/routes` directory may have subdirectories like so: `./myapp/routes/blog`. That is permitted.

## Compiling To `./myapp/dist`
The `./myapp/routes` directory is compiled into static `html` and is then copied over to the `./myapp/dist` directory with the exact same naming. For example, when we compile `./myapp/routes/index.html`, the components within `./myapp/index.html` will be resolved from the `./myapp/components` directory and then we copy the fully rendered `html` into the `./myapp/dist` directory.
