# Performance

While this implementation was designed with performance in mind, it's not
necessarily better, faster, smarter, or sexier than the built-in Go hashmap.

Here's a very rough comparison.

The following benchmarks were run on a Linux Desktop (3.8 GHz 8-Core AMD Ryzen 7 5800X) using Go version 1.19. The key types are either strings or ints and the values are always ints.

In all cases the maps start from zero capacity, like:

```go
m := make(map[string]int)      // go stdlib
var m hashmap.Map[string, int] // this package
```

```shell
MAPBENCH=100000 go test
```

## 100,000 random string keys

```shell
## STRING KEYS

-- tidwall --
set          100,000 ops     15ms      6,755,229/sec 
get          100,000 ops      9ms     11,480,505/sec 
reset        100,000 ops      6ms     15,393,057/sec 
scan              20 ops     12ms          1,643/sec 
delete       100,000 ops      8ms     12,024,028/sec 
memory     4,194,288 bytes                  41/entry 

-- stdlib --
set          100,000 ops     22ms      4,618,215/sec 
get          100,000 ops      9ms     10,560,800/sec 
reset        100,000 ops      6ms     15,763,082/sec 
scan              20 ops     16ms          1,275/sec 
delete       100,000 ops      9ms     11,645,027/sec 
memory     3,960,672 bytes                  39/entry 
```

## 100,000 random int keys

```shell
## INT KEYS

-- tidwall --
set          100,000 ops     13ms      7,655,039/sec 
get          100,000 ops      8ms     12,273,792/sec 
reset        100,000 ops      8ms     12,403,374/sec 
scan              20 ops     18ms          1,090/sec 
delete       100,000 ops     10ms     10,240,310/sec 
memory     3,143,448 bytes                  31/entry 

-- stdlib --
set          100,000 ops     10ms     10,356,513/sec 
get          100,000 ops      5ms     20,322,233/sec 
reset        100,000 ops      5ms     21,655,755/sec 
scan              20 ops     20ms          1,001/sec 
delete       100,000 ops      4ms     22,516,827/sec 
memory     2,767,704 bytes                  27/entry 
```

## 1,000,000 random string keys

```shell
## STRING KEYS

-- tidwall --
set        1,000,000 ops    236ms      4,234,342/sec 
get        1,000,000 ops    132ms      7,563,164/sec 
reset      1,000,000 ops    125ms      8,025,718/sec 
scan              20 ops    135ms            148/sec 
delete     1,000,000 ops    141ms      7,091,525/sec 
memory    67,108,848 bytes                  67/entry 

-- stdlib --
set        1,000,000 ops    351ms      2,851,085/sec 
get        1,000,000 ops    119ms      8,381,478/sec 
reset      1,000,000 ops    123ms      8,114,941/sec 
scan              20 ops    154ms            130/sec 
delete     1,000,000 ops    137ms      7,322,445/sec 
memory    57,931,472 bytes                  57/entry
```

## 1,000,000 random int keys

```shell
## INT KEYS

-- tidwall --
set        1,000,000 ops    105ms      9,490,462/sec 
get        1,000,000 ops     67ms     14,865,120/sec 
reset      1,000,000 ops     61ms     16,336,651/sec 
scan              20 ops    138ms            144/sec 
delete     1,000,000 ops     66ms     15,078,569/sec 
memory    50,329,272 bytes                  50/entry 

-- stdlib --
set        1,000,000 ops    155ms      6,469,942/sec 
get        1,000,000 ops     69ms     14,390,932/sec 
reset      1,000,000 ops     53ms     18,828,756/sec 
scan              20 ops    153ms            130/sec 
delete     1,000,000 ops     58ms     17,238,929/sec 
memory    40,146,664 bytes                  40/entry 
```

## 10,000,000 random string keys (int values)

```shell
## STRING KEYS

-- tidwall --
set       10,000,000 ops   2438ms      4,102,284/sec 
get       10,000,000 ops   1383ms      7,232,811/sec 
reset     10,000,000 ops   1443ms      6,928,121/sec 
scan              20 ops   1049ms             19/sec 
delete    10,000,000 ops   1691ms      5,915,335/sec 
memory   536,870,896 bytes                  53/entry 

-- stdlib --
set       10,000,000 ops   3572ms      2,799,176/sec 
get       10,000,000 ops   1447ms      6,911,076/sec 
reset     10,000,000 ops   1436ms      6,965,708/sec 
scan              20 ops   1772ms             11/sec 
delete    10,000,000 ops   1560ms      6,412,190/sec 
memory   463,468,240 bytes                  46/entry
```

## 10,000,000 random int keys (int values)

```shell
## INT KEYS

-- tidwall --
set       10,000,000 ops   1383ms      7,229,909/sec 
get       10,000,000 ops    754ms     13,260,571/sec 
reset     10,000,000 ops    778ms     12,860,468/sec 
scan              20 ops   1083ms             18/sec 
delete    10,000,000 ops    910ms     10,988,243/sec 
memory   402,650,808 bytes                  40/entry 

-- stdlib --
set       10,000,000 ops   1635ms      6,116,286/sec 
get       10,000,000 ops    760ms     13,154,962/sec 
reset     10,000,000 ops    852ms     11,740,869/sec 
scan              20 ops   1592ms             12/sec 
delete    10,000,000 ops    889ms     11,253,382/sec 
memory   321,976,040 bytes                  32/entry 
```
