# Documentation Site Overview

## Built with gtml
The gtml documentation website is built entirely using gtml itself. This serves as both a real-world test of the framework and a demonstration of its capabilities.

## Location
The documentation site source is located at `./docs` in the project root. This is a standard gtml project with the following structure:

```
./docs
./docs/components    # Reusable components for the docs site
./docs/routes        # Page routes
./docs/dist          # Compiled output
./docs/static        # Static assets
```

## Built with Preinstalled Components
The documentation website uses the preinstalled gtml components to showcase how they look and function. The site itself demonstrates practical usage of these components.

## Links Within the Site
Links in the documentation site use clean URLs without the `.html` extension:

```html
<!-- Correct -->
<a href="/docs">Documentation</a>
<a href="/components/buttons/button-primary">ButtonPrimary</a>

<!-- Not needed -->
<a href="/docs.html">Documentation</a>
```

## Compiling the Documentation Site
To compile the documentation site, run:

```bash
gtml compile ./docs
```

This outputs the compiled static HTML files to `./docs/dist`.
