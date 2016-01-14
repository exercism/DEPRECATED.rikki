The Go style guide specifies camel case for all names:

> [T]he convention in Go is to use MixedCaps or mixedCaps rather than underscores to write multiword names.- [Effective Go](https://golang.org/doc/effective_go.html#mixed-caps)

This also includes constants. Where other languages often use `ALL_CAPS`, Go uses `mixedCaps` for unexported constants and `MixedCaps` for exported ones--same as any other name.
