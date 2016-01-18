`go vet` is complaining about your solution.

From the documentation:

> Vet examines Go source code and reports suspicious constructs, such as Printf
> calls whose arguments do not align with the format string. Vet uses heuristics
> that do not guarantee all reports are genuine problems, but it can find errors
> not caught by the compilers.

To run it against the package in your current working directory, use:

```
go vet
```

The tool ships with Go, so you don't need to install anything extra in order to run it.

For more details about the various ways you can use it check out the
documentation with `go doc vet`.
