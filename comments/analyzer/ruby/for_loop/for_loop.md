Rubyists tend to prefer enumerable methods over `for` loops. `for` loops affect variables
outside of the iteration. Here's a good example:

```ruby
person = 'alice'
for person in ["alice", "bob", "charlie"]
  # your code
end
person #=> 'charlie'
```

It's worth being aware of the Ruby enumerable methods, since there are some
very powerful and expressive ones.

The most basic one is `Enumerable#each`, which can be used with enumerable
objects like `Array` and `Hash`:

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

`Enumerable` methods do not affect surounding variables.

```ruby
person = 'alice'
["alice", "bob", "charlie"].each do |person|
  # your code
end
person #=> 'alice'
```

Check out [Enumerable](http://ruby-doc.org/core-2.1.2/Enumerable.html) to see what's available.
