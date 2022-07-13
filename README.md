# Chainable Go Errors

This package provides two useful error types:
- Const: to represent constant type errors
- Error: to represent error chains or lists of errors

## Constant errors

The Const error is simply a string that supports the Error interface. This
allows us to write constant-time errors. These types of errors are useful
when we have very shallow error handling. For example a single struct and
its methods can probably make do internally with just constant errors.
Constant errors fall apart when errors are surfaced from deep in the call
stack from code we do not control. For these cases typically the error
originated from a 3rd partly library or even the Go standard library. These
errors are often custom in implementation so it becomes hard to distinguish
them.

# Chainable errors

The Error type is a pair or errors. This allows us to build an arbitrarily
long chain of errors that is not possible with fmt.Errorf's '%w' wrapping
which only allows a single wrapped error. This solves the problem identified
above where we want to wrap a 3rd party error with our own (typically
constant error) so it is easier to identify at other layers in our
application without losing the type information from the source error.

The Error type should be used with care because it enables lazily wrapping
errors to avoid proper fine-grained error handling. For example if we wrap
the source error several times with generic container errors, then the final
error handling code will include several possible matching conditions for
the errors.Is function. So it is strongly encouraged to use this
functionality sparingly at only at the point where the error originated.

This does not mean you should *not* have multiple chained errors however.
For example if you have a design in mind with classes of errors, then this
is a possible use-case that can be a good tradeoff.

