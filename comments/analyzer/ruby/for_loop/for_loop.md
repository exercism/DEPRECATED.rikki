Rubyists tend to prefer methods over `for` loops. `for` loops
affect variables outside of the iteration. Here's a good example:

```ruby
person = 'alice'
for person in ["alice", "bob", "charlie"]
  # your code
end
person #=> 'charlie'
```

It's worth being aware of the methods available for Ruby, since there are
some very powerful and expressive ones.

The most basic one is `each`, which can be used with objects like
`Array` and `Hash`:

```ruby
["alice", "bob", "charlie"].each do |name|
  puts name.upcase
end
```

```ruby
{"alice" => 51, "bob" => 18, "charlie" => 63}.each do |person, age|
  puts "#{person} is #{age} years old."
end
```

These blocks do not affect surounding variables.

```ruby
person = 'alice'
["alice", "bob", "charlie"].each do |person|
  # your code
end
person #=> 'alice'
```

Enumerable doesn't provide the `each` method, but it does require that
it exists.

Check out [Enumerable](http://ruby-doc.org/core/Enumerable.html) to
learn more about it.
