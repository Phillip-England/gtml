=====================================
# Testing Error Messages
=====================================

## Test Purpose
Validate that error messages are clear, accurate, and helpful.

## Error Should Include File Path

When a component has an error, the message should include the full file path:
```
Error in ./myapp/components/BrokenComponent.html: ...
```

## Error Should Include Line Number

When possible, errors should include line numbers:
```
Error in ./myapp/components/BrokenComponent.html:15: Unexpected token...
```

## Error Should Be Descriptive

Error messages should explain what went wrong and how to fix it:
```
Error: Component 'Button' expects prop 'text' of type 'string', but it was not provided.
```

## Error Should Not Crash Silently

All errors should produce visible error messages, not silent failures or empty output.
