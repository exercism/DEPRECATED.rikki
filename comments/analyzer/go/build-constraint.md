The `+build !example` comment is a _build constraint_.

It's kind of complicated, and the short version is: you don't need
it, so you can safely delete it.

The longer version is that the Go track on Exercism has an automated
test suite that verifies all the exercises on a Continuous Integration
(CI) server. In other words, there are tests to test the tests.

Each exercise checks the test suite against a reference solution. The
problem is when there's also a file with stub code. They tend to have
duplicate definitions or code that doesn't compile, which makes the
build fail.

To avoid that, the CI server runs the tests with a tag, ignoring the
stub file.

```
go test --tags example ./...
```
