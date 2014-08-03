In Ruby `and` has similar behavior to `&&`, but they mean very different things.

`&&`, `||`, and `!` are logical operators. These are meant to combine multiple boolean (`true`/`false`) statements into one single `true` or `false`.

For example:

```
hungry? && thirsty? || cranky?
```

The more English-y operators, `and`, `or`, and `not` are control flow operators.

These get strung together when you need to decide whether or not to perform an operation based on whether or not a previous operation failed or succeeded.

For example:

```
# only drink coffee if you're sleep deprived
sleep or drink_coffee
```

or

```
# you only get to celebrate if the deploy succeeded
deploy and celebrate
```
