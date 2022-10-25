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
set          100,000 ops     15ms      6,535,956/sec 
get          100,000 ops      7ms     15,009,648/sec 
reset        100,000 ops      5ms     20,811,745/sec 
scan              20 ops      8ms          2,539/sec 
delete       100,000 ops      7ms     14,557,933/sec 
memory     4,194,288 bytes                  41/entry 

-- stdlib --
set          100,000 ops     17ms      5,892,223/sec 
get          100,000 ops      8ms     12,148,359/sec 
reset        100,000 ops      4ms     24,779,419/sec 
scan              20 ops     14ms          1,395/sec 
delete       100,000 ops      8ms     11,915,708/sec 
memory     3,966,288 bytes                  39/entry
```

## 100,000 random int keys

```shell
## INT KEYS

-- tidwall --
set          100,000 ops      8ms     12,573,069/sec 
get          100,000 ops      4ms     24,821,181/sec 
reset        100,000 ops      4ms     25,324,412/sec 
scan              20 ops      8ms          2,430/sec 
delete       100,000 ops      5ms     22,156,034/sec 
memory     3,143,352 bytes                  31/entry 

-- stdlib --
set          100,000 ops      7ms     13,547,121/sec 
get          100,000 ops      4ms     26,458,302/sec 
reset        100,000 ops      4ms     28,379,163/sec 
scan              20 ops     16ms          1,214/sec 
delete       100,000 ops      4ms     24,771,495/sec 
memory     2,784,264 bytes                  27/entry 
```

## 1,000,000 random string keys

```shell
## STRING KEYS

-- tidwall --
set        1,000,000 ops    217ms      4,607,387/sec 
get        1,000,000 ops    127ms      7,872,817/sec 
reset      1,000,000 ops    130ms      7,709,027/sec 
scan              20 ops    136ms            147/sec 
delete     1,000,000 ops    149ms      6,716,045/sec 
memory    67,108,848 bytes                  67/entry 

-- stdlib --
set        1,000,000 ops    325ms      3,078,132/sec 
get        1,000,000 ops    122ms      8,217,771/sec 
reset      1,000,000 ops    133ms      7,510,273/sec 
scan              20 ops    163ms            122/sec 
delete     1,000,000 ops    148ms      6,761,332/sec 
memory    57,931,472 bytes                  57/entry
```

## 1,000,000 random int keys

```shell
## INT KEYS

-- tidwall --
set        1,000,000 ops    101ms      9,901,395/sec 
get        1,000,000 ops     63ms     15,928,770/sec 
reset      1,000,000 ops     66ms     15,107,262/sec 
scan              20 ops    139ms            144/sec 
delete     1,000,000 ops     66ms     15,216,322/sec 
memory    50,329,272 bytes                  50/entry 

-- stdlib --
set        1,000,000 ops    119ms      8,431,961/sec 
get        1,000,000 ops     61ms     16,376,595/sec 
reset      1,000,000 ops     59ms     17,032,395/sec 
scan              20 ops    153ms            130/sec 
delete     1,000,000 ops     67ms     15,026,654/sec 
memory    40,146,760 bytes                  40/entry 
```

## 10,000,000 random string keys (int values)

```shell
## STRING KEYS

-- tidwall --
set       10,000,000 ops   2584ms      3,869,389/sec 
get       10,000,000 ops   1418ms      7,051,328/sec 
reset     10,000,000 ops   1469ms      6,807,487/sec 
scan              20 ops   1049ms             19/sec 
delete    10,000,000 ops   1694ms      5,901,787/sec 
memory   536,870,896 bytes                  53/entry 

-- stdlib --
set       10,000,000 ops   3771ms      2,651,828/sec 
get       10,000,000 ops   1494ms      6,695,021/sec 
reset     10,000,000 ops   1480ms      6,758,881/sec 
scan              20 ops   1855ms             10/sec 
delete    10,000,000 ops   1629ms      6,138,209/sec 
memory   463,468,240 bytes                  46/entry
```

## 10,000,000 random int keys (int values)

```shell
## INT KEYS

-- tidwall --
set       10,000,000 ops   1428ms      7,002,173/sec 
get       10,000,000 ops    733ms     13,636,196/sec 
reset     10,000,000 ops    787ms     12,710,144/sec 
scan              20 ops   1098ms             18/sec 
delete    10,000,000 ops    900ms     11,108,541/sec 
memory   402,650,808 bytes                  40/entry 

-- stdlib --
set       10,000,000 ops   1709ms      5,850,969/sec 
get       10,000,000 ops    797ms     12,551,221/sec 
reset     10,000,000 ops    874ms     11,437,820/sec 
scan              20 ops   1629ms             12/sec 
delete    10,000,000 ops    910ms     10,994,436/sec 
memory   321,976,032 bytes                  32/entry 
```
