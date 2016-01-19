In Go the idiom is to ditch the else if there is a return in the `if` block. This minimizes indentation.

[Effective Go](https://golang.org/doc/effective_go.html#if) says:

> In the Go libraries, you'll find that when an if statement doesn't flow into the next statement—that is, the body ends in break, continue, goto, or return—the unnecessary else is omitted.
