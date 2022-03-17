## Performance

A very rough comparison of this implementation and the built-in Go map.

The following benchmarks were run on my 2019 Macbook Pro (2.4 GHz 8-Core Intel Core i9) using Go version 1.18. The key types are either strings or ints and the values are always ints.

In all cases the maps start from zero capacity, like:

```go
m := make(map[string]int)      // go stdlib
var m hashmap.Map[string, int] // this package
```

```
MAPBENCH=100000 go test
```

### 100,000 random string keys

```
## STRING KEYS

-- tidwall --
set          100,000 ops     19ms      5,405,235/sec
get          100,000 ops     12ms      8,599,131/sec
reset        100,000 ops     12ms      8,463,694/sec
scan              20 ops     14ms          1,428/sec
delete       100,000 ops     11ms      8,749,727/sec
memory     4,194,288 bytes                  41/entry

-- stdlib --
set          100,000 ops     28ms      3,602,606/sec
get          100,000 ops     17ms      5,842,126/sec
reset        100,000 ops     13ms      7,489,159/sec
scan              20 ops     32ms            627/sec
delete       100,000 ops     13ms      7,545,664/sec
memory     3,968,784 bytes                  39/entry
```

### 100,000 random int keys

```
## INT KEYS

-- tidwall --
set          100,000 ops     10ms      9,624,083/sec
get          100,000 ops      5ms     21,056,856/sec
reset        100,000 ops      5ms     21,281,182/sec
scan              20 ops     10ms          1,917/sec
delete       100,000 ops      5ms     18,342,582/sec
memory     3,143,352 bytes                  31/entry

-- stdlib --
set          100,000 ops     10ms     10,354,476/sec
get          100,000 ops      4ms     25,693,552/sec
reset        100,000 ops      4ms     24,752,983/sec
scan              20 ops     21ms            967/sec
delete       100,000 ops      5ms     19,239,275/sec
memory     2,772,744 bytes                  27/entry
```

### 1,000,000 random string keys

```
## STRING KEYS

-- tidwall --
set        1,000,000 ops    299ms      3,342,247/sec
get        1,000,000 ops    141ms      7,099,822/sec
reset      1,000,000 ops    155ms      6,444,007/sec
scan              20 ops    244ms             82/sec
delete     1,000,000 ops    178ms      5,622,242/sec
memory    67,108,848 bytes                  67/entry

-- stdlib --
set        1,000,000 ops    426ms      2,348,509/sec
get        1,000,000 ops    142ms      7,019,978/sec
reset      1,000,000 ops    182ms      5,496,549/sec
scan              20 ops    297ms             67/sec
delete     1,000,000 ops    187ms      5,337,449/sec
memory    57,931,472 bytes                  57/entry
```

### 1,000,000 random int keys

```
## INT KEYS

-- tidwall --
set        1,000,000 ops    146ms      6,838,679/sec
get        1,000,000 ops     72ms     13,798,363/sec
reset      1,000,000 ops     76ms     13,236,277/sec
scan              20 ops    243ms             82/sec
delete     1,000,000 ops    112ms      8,893,494/sec
memory    50,329,280 bytes                  50/entry

-- stdlib --
set        1,000,000 ops    171ms      5,850,975/sec
get        1,000,000 ops     71ms     14,096,964/sec
reset      1,000,000 ops     75ms     13,279,320/sec
scan              20 ops    285ms             70/sec
delete     1,000,000 ops     90ms     11,131,406/sec
memory    40,146,760 bytes                  40/entry
```

### 10,000,000 random string keys (int values)

```
## STRING KEYS

-- tidwall --
set       10,000,000 ops   3185ms      3,139,265/sec
get       10,000,000 ops   1624ms      6,158,609/sec
reset     10,000,000 ops   1816ms      5,505,645/sec
scan              20 ops   2037ms              9/sec
delete    10,000,000 ops   2137ms      4,678,937/sec
memory   536,870,904 bytes                  53/entry

-- stdlib --
set       10,000,000 ops   4623ms      2,163,071/sec
get       10,000,000 ops   1764ms      5,670,292/sec
reset     10,000,000 ops   2389ms      4,185,314/sec
scan              20 ops   2975ms              6/sec
delete    10,000,000 ops   2280ms      4,385,089/sec
memory   463,468,352 bytes                  46/entry

```

### 10,000,000 random int keys (int values)

```
## INT KEYS

-- tidwall --
set       10,000,000 ops   1743ms      5,736,861/sec
get       10,000,000 ops    878ms     11,394,400/sec
reset     10,000,000 ops   1003ms      9,965,713/sec
scan              20 ops   1973ms             10/sec
delete    10,000,000 ops   1126ms      8,883,210/sec
memory   402,650,808 bytes                  40/entry

-- stdlib --
set       10,000,000 ops   1989ms      5,027,700/sec
get       10,000,000 ops   1021ms      9,796,267/sec
reset     10,000,000 ops   1033ms      9,685,070/sec
scan              20 ops   2847ms              7/sec
delete    10,000,000 ops   1182ms      8,456,779/sec
memory   321,976,040 bytes                  32/entry
```
