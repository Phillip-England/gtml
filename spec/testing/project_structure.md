# Testing Project Structure

## Test Purpose
Validate that `gtml` correctly enforces the required project structure and fails appropriately when the structure is invalid.

## Valid Project Structure
A valid project must have all required directories:
```bash
./myapp
./myapp/components
./myapp/routes
./myapp/dist
./myapp/static
```

Test that compilation succeeds when all directories exist.

## Missing Components Directory
If `./myapp/components` is missing, compilation should fail with a clear error message indicating the missing directory.

## Missing Routes Directory
If `./myapp/routes` is missing, compilation should fail with a clear error message indicating the missing directory.

## Missing Dist Directory
If `./myapp/dist` is missing, it should either be auto-created or compilation should fail with a clear error. Test for consistent behavior.

## Missing Static Directory
If `./myapp/static` is missing, it should either be auto-created or compilation should fail with a clear error. Test for consistent behavior.

## Empty Components Directory
An empty `./myapp/components` directory should be valid (routes may not use any components).

## Empty Routes Directory
An empty `./myapp/routes` directory should either be valid (no output) or produce a warning. Test for consistent behavior.
