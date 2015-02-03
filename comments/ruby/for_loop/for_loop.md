Rubyists tend to prefer enumerable methods over `for` loops. This isn't
a rule, it's more of a strong cultural preference.

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

Check out [Enumerable](http://ruby-doc.org/core-2.1.2/Enumerable.html) to see what's available.
