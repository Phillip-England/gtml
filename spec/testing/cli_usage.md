# Testing CLI Commands

## Test Purpose
Validate that CLI commands work correctly and produce expected results.

## Testing `gtml init <PATH>`

### Basic Init
Running `gtml init myapp` should create:
```bash
./myapp
./myapp/components
./myapp/components/BasicButton.html
./myapp/components/GuestLayout.html
./myapp/routes
./myapp/routes/index.html
./myapp/static
./myapp/dist
./myapp/dist/index.html
./myapp/dist/static
```

### Init Fails If Directory Exists
Running `gtml init myapp` when `./myapp` already exists should fail with an error.

### Init With Force Flag
Running `gtml init myapp --force` when `./myapp` already exists should succeed and reinitialize the project.

### Init Creates Valid Default Files
After `gtml init myapp`, the generated files should be valid and compilable. Running `gtml compile myapp` immediately after init should succeed.

## Testing `gtml compile <PATH>`

### Basic Compile
Running `gtml compile myapp` on a valid project should:
1. Read all components from `./myapp/components`
2. Process all routes from `./myapp/routes`
3. Output compiled HTML to `./myapp/dist`
4. Copy static assets to `./myapp/dist/static`
5. Generate `./myapp/dist/static/styles.css`

### Compile Fails On Invalid Project
Running `gtml compile myapp` when `./myapp` is missing required directories should fail with a clear error.

### Compile Reports Component Errors
If a component has invalid syntax, compilation should fail with an error that identifies the file and line number.

### Watch Mode
Running `gtml compile myapp --watch` should:
1. Perform initial compilation
2. Watch for file changes in `./myapp`
3. Recompile when files change
