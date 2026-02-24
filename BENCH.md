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
-- tidwall --
set          100,000 ops     17ms      5,999,685/sec
get          100,000 ops     10ms      9,620,504/sec
reset        100,000 ops     11ms      8,701,108/sec
scan              20 ops      2ms          8,374/sec
delete       100,000 ops     20ms      4,962,668/sec
memory     8,388,608 bytes                  83/entry

-- stdlib --
set          100,000 ops     25ms      4,063,454/sec
get          100,000 ops      9ms     11,170,147/sec
reset        100,000 ops     11ms      9,487,623/sec
scan              20 ops     24ms            846/sec
delete       100,000 ops     11ms      9,132,136/sec
memory     3,494,552 bytes                  34/entry
```

## 100,000 random int keys

```shell
-- tidwall --
set          100,000 ops     13ms      7,645,884/sec
get          100,000 ops      5ms     20,073,198/sec
reset        100,000 ops      3ms     32,332,261/sec
scan              20 ops      3ms          7,385/sec
delete       100,000 ops     14ms      7,275,010/sec
memory     6,291,016 bytes                  62/entry

-- stdlib --
set          100,000 ops     12ms      8,384,588/sec
get          100,000 ops      6ms     16,482,724/sec
reset        100,000 ops      6ms     17,029,795/sec
scan              20 ops     24ms            845/sec
delete       100,000 ops      8ms     13,309,096/sec
memory     2,364,088 bytes                  23/entry
```

## 1,000,000 random string keys

```shell
-- tidwall --
set        1,000,000 ops    266ms      3,756,183/sec
get        1,000,000 ops    132ms      7,589,748/sec
reset      1,000,000 ops    125ms      8,025,198/sec
scan              20 ops      9ms          2,213/sec
delete     1,000,000 ops    157ms      6,351,239/sec
memory    67,108,864 bytes                  67/entry

-- stdlib --
set        1,000,000 ops    342ms      2,925,835/sec
get        1,000,000 ops    142ms      7,059,901/sec
reset      1,000,000 ops    190ms      5,257,840/sec
scan              20 ops    150ms            133/sec
delete     1,000,000 ops    199ms      5,016,100/sec
memory    55,783,672 bytes                  55/entry
```

## 1,000,000 random int keys

```shell
-- tidwall --
set        1,000,000 ops     88ms     11,421,305/sec
get        1,000,000 ops     57ms     17,619,341/sec
reset      1,000,000 ops     56ms     17,914,738/sec
scan              20 ops      9ms          2,212/sec
delete     1,000,000 ops     66ms     15,038,347/sec
memory    50,331,208 bytes                  50/entry

-- stdlib --
set        1,000,000 ops    115ms      8,702,795/sec
get        1,000,000 ops     76ms     13,082,493/sec
reset      1,000,000 ops     62ms     16,050,479/sec
scan              20 ops    149ms            134/sec
delete     1,000,000 ops     71ms     14,162,994/sec
memory    37,703,000 bytes                  37/entry
```

## 10,000,000 random string keys (int values)

```shell
-- tidwall --
set       10,000,000 ops   2910ms      3,436,974/sec
get       10,000,000 ops   1512ms      6,611,683/sec
reset     10,000,000 ops   1545ms      6,472,422/sec
scan              20 ops     73ms            274/sec
delete    10,000,000 ops   2053ms      4,871,632/sec
memory   536,870,912 bytes                  53/entry

-- stdlib --
set       10,000,000 ops   5234ms      1,910,708/sec
get       10,000,000 ops   1798ms      5,561,131/sec
reset     10,000,000 ops   2758ms      3,625,603/sec
scan              20 ops   1444ms             13/sec
delete    10,000,000 ops   2932ms      3,410,416/sec
memory   447,348,248 bytes                  44/entry
```

## 10,000,000 random int keys (int values)

```shell
-- tidwall --
set       10,000,000 ops   1060ms      9,430,586/sec
get       10,000,000 ops    595ms     16,799,748/sec
reset     10,000,000 ops    790ms     12,659,392/sec
scan              20 ops     71ms            282/sec
delete    10,000,000 ops   1045ms      9,572,199/sec
memory   402,652,744 bytes                  40/entry

-- stdlib --
set       10,000,000 ops   1596ms      6,266,321/sec
get       10,000,000 ops   1326ms      7,538,739/sec
reset     10,000,000 ops   1503ms      6,655,292/sec
scan              20 ops   1513ms             13/sec
delete    10,000,000 ops   1546ms      6,468,809/sec
memory   302,644,792 bytes                  30/entry
```
