# The Spec
"The spec" is all of the documentation outlined in `./spec`. It describes, in natural language, how our program ought to work and be constructed. Anything you need to know about the project should be available somewhere in the spec.

## The Spec as an Intermediate Representation
The spec serves as a sort of IR. When we makes changes to this program, we do so via the spec. We do not read the code nor do we write the code. Rather, we construct the spec in such a way that the spec itself may be used to generate a fully-operative program via a LLM.

## Prone to Logical Errors
Because the spec is written in natural language, it is prone to logically errors in the same way a program is itself. The spec is itself a sort of program. In the same way a program expresses what the computer ought to do, so does the spec.
