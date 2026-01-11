# The Compilation Process

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
