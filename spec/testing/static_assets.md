# Testing Static Assets

## Test Purpose
Validate that static assets are handled correctly during compilation.

## Static Directory Copied

Files in `./myapp/static` should be copied to `./myapp/dist/static`:
- `./myapp/static/image.png` → `./myapp/dist/static/image.png`
- `./myapp/static/fonts/custom.woff` → `./myapp/dist/static/fonts/custom.woff`

## Static Assets Not Processed

Static files should be copied as-is, not processed or modified.

## Referencing Static Assets

Routes and components should be able to reference static assets:
```html
<img src='/static/logo.png' alt='Logo' />
<link rel='stylesheet' href='/static/custom.css' />
```
