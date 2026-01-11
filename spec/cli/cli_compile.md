# Command Line Compile

## `gtml compile ./somedir`
The `compile` command will attempt to compile the routes found at `./somedir/routes`. This command will check to ensure our project structure is correct and that everything checks out. Upon failure, this command will let you know exactly why things failed. If things are successful, you should have your static `html` in `./somedir/dist`

If you pass the `--watch` flag, changes to the any file within the `./somedir` directory will trigger recompilation.
