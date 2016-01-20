Avoid naming your durations with unit-specific suffixes like `Secs`,
`Seconds`, `Ms`, `Millis` or `Milliseconds`.

The `time.Duration` type represents a span of time which can be expressed
in many different units, e.g. seconds, milliseconds, microseconds,
and nanoseconds).

When you pass a duration to `fmt.Print`, it displays the value in a way that
makes sense for that particular duration:

```
fmt.Println(1e6 * time.Millisecond)
// 16m40s

fmt.Println(31415 * time.Second)
// 8h43m35s
```
