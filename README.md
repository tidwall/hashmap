# hashmap

[![GoDoc](https://img.shields.io/badge/api-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/tidwall/hashmap)

An [efficient](BENCH.md) hashmap implementation in Go.

## Features

- Support for [Generics](#generics).
- `Map` and `Set` types for unordered key-value maps and sets,
- [xxh3 algorithm](https://github.com/zeebo/xxh3)
- [Open addressing](https://en.wikipedia.org/wiki/Hash_table#Open_addressing) with [Robin hood hashing](https://en.wikipedia.org/wiki/Hash_table#Robin_Hood_hashing)
- Automatically shinks memory on deletes (no memory leaks).
- Pretty darn good performance. 🚀 ([benchmarks](BENCH.md)).

For ordered key-value data, check out the [tidwall/btree](https://github.com/tidwall/btree) package.

## Getting Started

### Installing

To start using `hashmap`, install Go and run `go get`:

```sh
go get github.com/tidwall/hashmap
```

This will retrieve the library.

## Usage

The `Map` type works similar to a standard Go map, and includes the methods:
`Set`, `Get`, `Delete`, `Len`, `Scan`, `Keys`, `Values`, and `Copy`.

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
It includes the methods: `Insert`, `Contains`, `Delete`, `Len`, `Scan` and `Keys`.

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

See [BENCH.md](BENCH.md) for more info.

## Contact

Josh Baker [@tidwall](http://twitter.com/tidwall)

## License

Source code is available under the MIT [License](LICENSE).
