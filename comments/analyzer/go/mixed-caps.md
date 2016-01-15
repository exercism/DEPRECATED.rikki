The Go style guide specifies camel case for all names:

> [T]he convention in Go is to use MixedCaps or mixedCaps rather than underscores to write multiword names.- [Effective Go](https://golang.org/doc/effective_go.html#mixed-caps)

_Names_ here refers to anything that you would name: types, functions, methods, variables (whether they're package level or local), and constants.

Note that _constants_ is included in the above list. A lot of languages use `ALL_CAPS` for these, but Go uses `mixedCaps` for unexported constants and `MixedCaps` for exported ones--same as any other name.
