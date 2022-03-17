# hashmap

[![GoDoc](https://img.shields.io/badge/api-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/tidwall/hashmap)

An efficient hashmap implementation in Go.

## Features

- Support for [Generics](#generics) (Go 1.18+).
- `Map` and `Set` types for unordered key-value maps and sets,
- [xxhash algorithm](http://www.xxhash.com)
- [Open addressing](https://en.wikipedia.org/wiki/Hash_table#Open_addressing)
- [Robin hood hashing](https://en.wikipedia.org/wiki/Hash_table#Robin_Hood_hashing)

For ordered key-value data, check out the [tidwall/btree](https://github.com/tidwall/btree) package.

## Using

To start using this package, install Go and run:

```sh
$ go get github.com/tidwall/hashmap
```

# Getting Started

## Installing

To start using `hashmap`, install Go and run `go get`:

```sh
$ go get github.com/tidwall/hashmap
```

This will retrieve the library.

## Usage

The `Map` type works similar to a standard Go map, and includes the methods:
`Set`, `Get`, `Delete`, `Len`, and `Scan`.

```go
var m hashmap.Map[string, string]
m.Set("Hello", "Dolly!")
val, _ := m.Get("Hello")
fmt.Printf("%v\n", val)
val, _ = m.Delete("Hello")
fmt.Printf("%v\n", val)
val, _ = m.Get("Hello")
fmt.Printf("%v\n", val)

// Output:
// Dolly!
// Dolly!
//
```

The `Set` type is like `Map` but only for keys.
It includes the methods: `Insert`, `Contains`, `Delete`, `Len`, and `Scan`.

```go
var m hashmap.Set[string]
m.Insert("Andy")
m.Insert("Kate")
m.Insert("Janet")

fmt.Printf("%v\n", m.Contains("Kate"))
fmt.Printf("%v\n", m.Contains("Bob"))
fmt.Printf("%v\n", m.Contains("Andy"))

// Output:
// true
// false
// true
```

## Performance

While this implementation was designed with performance in mind, it's not 
necessarily better, faster, smarter, or sexier that the built-in Go hashmap. 

See [BENCH.md](BENCH.md) for mor info.

## Contact

Josh Baker [@tidwall](http://twitter.com/tidwall)

## License

Source code is available under the MIT [License](/LICENSE).
