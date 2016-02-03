Go programmers value consistency very highly.

If you define methods on a type, then the receiver name should be
the same in all the method definitions.

For example, imagine that you have a type `Bacterium`, and define a method `Grow()` for that type that
names the receiver `b`:

```
func (b Bacterium) Grow() {
	// ....
}
```

If you define a second method, `Split()` for the `Bacterium` type, then don't use
`bact` or `bacterium` or `originalBacterium` for the receiver name. Use `b`, since
that is what you used before.

```
func (b Bacterium) Split() (Bacterium, Bacterium) {
	// ...
}
```

