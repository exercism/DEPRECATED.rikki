Rubyists tend to prefer enumerable methods over `for` loops.

The most basic choice is `Enumerable#each`, which can be used enumerable objects like `Array` and `Hash`:

```ruby
["alice", "bob", "charlie"].each do |name|
  puts name.upcase
end
```

```ruby
{"alice" => 51, "bob" => 18, "charlie" => 63}.each do |name, age|
  puts "#{name} is #{age} years old."
end
```

There are many useful enumerable methods beyond `each` that you might want to use, depending on what you're trying to do, but everything is based on `each` so it's a good place to start.

Check out [Enumerable](http://ruby-doc.org/core-2.1.2/Enumerable.html) to see what's available.
