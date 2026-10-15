# Assignment 7: Statistics calculator

A statistics calculator in C, OCaml, and Python. The three programs follow one specification and print the same output, apart from the first line. Each program uses a different paradigm: C is procedural, OCaml is functional, and Python is object-oriented.

## What the application does

Each program calculates statistics on three fixed demo lists of integers and prints them. It takes no input, runs once, and exits with status 0.

| Statistic | Rule |
|---|---|
| Count | The number of integers in the list |
| Mean | The sum divided by the count, with 2 decimal places |
| Median | The middle value of the sorted list. With an even count, the average of the two middle values |
| Mode | Every value that ties for the highest number of occurrences, in ascending order |
| Min | The smallest value |
| Max | The largest value |
| Range | Max minus min |

When every value occurs once, every value is a mode. Each demo list shows one case:

| List | Case | Integers |
|---|---|---|
| 1 | An odd count with one mode | `[4, 8, 15, 16, 23, 42, 8]` |
| 2 | An even count with two tied modes | `[7, 3, 9, 3, 7, 1, 12, 2]` |
| 3 | Negative values, and every value occurs once | `[5, -2, 10, 0, -7]` |

## Paradigm of each language

### C (procedural)

- **Functions over an array.** `stats.c` has one function for each statistic. Each function takes an array and its length.
- **Manual memory management.** The program copies each demo list into an array from `malloc`. It calls `free` on every array that it allocates, and Valgrind reports no leak.
- **Sort.** A hand-written insertion sort works on a copy of the array, so the demo list keeps its order. The program does not use `qsort`.
- **Mode.** The `mode` function counts the runs of equal values in the sorted copy and finds the longest run. It then collects every value with a run of that length into a new array from `malloc`. The caller frees that array.

### OCaml (functional)

- **Immutable data.** The data is an immutable `int list`. The program has no `ref`, no array, no mutable record field, and no loop.
- **Higher-order functions.** The statistics use `List.sort`, `List.fold_left`, `List.map`, and `List.filter`.
- **Mode.** A fold over the sorted list builds the pairs of each value and its number of occurrences. The `mode` function then keeps the values with the highest number.

### Python (object-oriented)

- **Class.** The `StatisticsCalculator` class keeps the list in the private attribute `_values`. Each statistic is a method: `count`, `mean`, `median`, `mode`, `min`, `max`, and `range`.
- **Mode.** The `mode` method counts the values in a `dict`.
- **Calculation and output apart.** The class does the calculation, and the separate function `print_report` does the output.
- **Built-in types only.** The program does not use the `statistics` module or `collections.Counter`.

## Run with Docker

Docker is the primary run path. One image holds the three programs. Run these commands in this folder.

Build the image:

```sh
docker build -t statistics-calculator .
```

Run the three programs in order:

```sh
docker run --rm statistics-calculator
```

Run one program:

```sh
docker run --rm statistics-calculator c
docker run --rm statistics-calculator ocaml
docker run --rm statistics-calculator python
```

## Run without Docker

### C

The C program needs GCC or Clang with C11 support and uses the C standard library only. Use these commands in the `c` folder to build the program and run it:

```sh
cd c
gcc -std=c11 -Wall -Wextra -O2 -o stats stats.c
./stats
```

To check the memory management on Linux, run the program under Valgrind:

```sh
valgrind --leak-check=full ./stats
```

### OCaml

The OCaml program needs OCaml 4.14 or later and uses the standard library only. Use these commands in the `ocaml` folder to build the program and run it:

```sh
cd ocaml
ocamlopt -o stats stats.ml
./stats
```

### Python

The Python program needs Python 3.8 or later and uses built-in types only. Use these commands in the `python` folder to run the program:

```sh
cd python
python3 stats.py
```

## Sample output

The output of the C program:

```
=== Statistics Calculator: C (procedural) ===

List 1: [4, 8, 15, 16, 23, 42, 8]
  Count:  7
  Mean:   16.57
  Median: 15.00
  Mode:   [8]
  Min:    4
  Max:    42
  Range:  38

List 2: [7, 3, 9, 3, 7, 1, 12, 2]
  Count:  8
  Mean:   5.50
  Median: 5.00
  Mode:   [3, 7]
  Min:    1
  Max:    12
  Range:  11

List 3: [5, -2, 10, 0, -7]
  Count:  5
  Mean:   1.20
  Median: 0.00
  Mode:   [-7, -2, 0, 5, 10]
  Min:    -7
  Max:    10
  Range:  17
```

The OCaml program and the Python program print the same lines, apart from the first line:

```
=== Statistics Calculator: OCaml (functional) ===
=== Statistics Calculator: Python (object-oriented) ===
```

## Project layout

```
Assignment-7/
  Dockerfile      the C, OCaml, and Python programs in one image
  README.md
  screenshots/
  c/
    stats.c       the procedural program: one function for each statistic, insertion sort, malloc, and free
  ocaml/
    stats.ml      the functional program: immutable lists, folds, and no mutable state
  python/
    stats.py      the object-oriented program: the StatisticsCalculator class and the report function
```
