Instead of ignoring the second value in your loop with an underscore, you
 can simplify and omit it altogether.

This works when ranging over both slices and maps.

For example, if you are comparing two slices of int that you expect to be the same, you can say:

```
for i := range expected {
	if expected[i] != actual[i] {
		// something went wrong
	}
}
```

Likewise, when you're ranging over a map and only need the key, you can omit the value.

```
children := map[string]int{
	`alice`:   9,
	`bob`:     12,
	`charlie`: 7,
}
for name := range children {
	// do something with the name
}
```
