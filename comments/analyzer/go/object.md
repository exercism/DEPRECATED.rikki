Go doesn't have a concept of _objects_.

Given a type `Cupcake`

```
type Cupcake struct {
	Type string
	Grams int
}
```

when you create one with

```
c := Cupcake{
	Type: `blueberry`,
	Grams: 275,
}
```

a Go programmer would say that you're creating a _value of type Cupcake_.
