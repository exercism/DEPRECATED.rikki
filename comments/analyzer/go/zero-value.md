In Go if you don't initialize a variable when you declare it, the value of that variable will be its _zero value_.

The zero values are:

- `false` for `bool`
- `0` for `int`
- `0.0` for `float`
- `""` for strings
- `nil` for everything else (pointers, functions, interfaces, slices, channels, and maps).

Because of this, the convention is to not explicitly initialize a variable if you want the zero value.

So, instead of:

```
var foo int = 0
```

Gophers will write:

```
var foo int
```
