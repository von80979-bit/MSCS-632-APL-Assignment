# Assignment 8: Family tree program in Prolog

A family tree program in Prolog. Facts hold an invented family of 4 generations, and rules derive the other relationships from these facts. The program answers questions about parents, grandparents, siblings, cousins, and descendants through the pattern matching and the logical inference of Prolog. The descendant rule uses recursion to find the relationships across more than one generation.

## What the program does

The program runs one fixed demo and exits. The `demo/0` rule runs 15 sample queries in a fixed order and prints each query with its answer, in the format of the SWI-Prolog prompt. You can also load the file at the SWI-Prolog prompt and type your own queries.

## Family

The family has 13 people in 4 generations.

| Generation | People | Parents |
|---|---|---|
| 1 | `george` (male), `helen` (female) | None |
| 2 | `david` (male), `susan` (female), `mark` (male) | `george` and `helen` |
| 2, married in | `laura` (female), the partner of `david`. `paul` (male), the partner of `susan` | None |
| 3 | `emma` (female), `jack` (male) | `david` and `laura` |
| 3 | `olivia` (female) | `susan` and `paul` |
| 3, married in | `ryan` (male), the partner of `emma` | None |
| 4 | `lily` (female), `noah` (male) | `emma` and `ryan` |

## Facts and rules

| Predicate | Kind | Meaning |
|---|---|---|
| `parent(Parent, Child)` | Fact | `Parent` is a parent of `Child` |
| `male(Person)`, `female(Person)` | Fact | The gender of `Person` |
| `grandparent(X, Y)` | Rule | `X` is a parent of a parent of `Y` |
| `sibling(X, Y)` | Rule | `X` and `Y` share the same parent, and `X` is not `Y` |
| `cousin(X, Y)` | Rule | The parents of `X` and `Y` are siblings |
| `descendant(D, A)` | Recursive rule | `D` is a child of `A`, or a descendant of a child of `A` |
| `ancestor(A, D)` | Rule | `D` is a descendant of `A` |
| `child(C, P)` | Rule | `P` is a parent of `C` |
| `children_of(Person, Children)` | Query helper | The sorted list of the children of `Person` |
| `siblings_of(Person, Siblings)` | Query helper | The sorted list of the siblings of `Person` |
| `descendants_of(Person, Descendants)` | Query helper | The sorted list of all descendants of `Person` |

The `descendant/2` rule has a base case and a recursive step. The base case states that a child is a descendant. The recursive step states that a descendant of a child is a descendant. The recursion stops because the family tree has no cycle, and the people in the last generation have no children.

The query helpers collect the answers with `setof/3`, which sorts the answers and removes the duplicates. The sibling rule and the cousin rule find the same answer once for each shared parent, so the helpers give each person only once.

## Assumptions

The requirement leaves certain points open, so the program uses these assumptions:

1. The family is invented. The names are lowercase atoms, for example `george`.
2. The family has no remarriage. Every child has two parents in the facts, and siblings always share both parents. The sibling rule follows the requirement text, "share the same parent", and the data gives no half-siblings.
3. A person who marries into the family has no parents in the facts.
4. A query helper that finds no answer returns the empty list `[]`, not `false`.

## Run with Docker

Docker is the primary run path. Run these commands in this folder.

Build the image:

```sh
docker build -t family-tree .
```

Run the demo:

```sh
docker run --rm family-tree
```

Open the SWI-Prolog prompt to type queries by hand:

```sh
docker run --rm -it family-tree swipl family.pl
```

## Run without Docker

The program needs SWI-Prolog 9 or later and uses no library outside the SWI-Prolog distribution. Use these commands in the `prolog` folder.

Run the demo:

```sh
cd prolog
swipl -g demo -t halt family.pl
```

Open the SWI-Prolog prompt to type queries by hand:

```sh
swipl family.pl
```

## Queries to type by hand

Type each query at the `?-` prompt. When Prolog shows one answer and waits, type `;` for the next answer or press Enter to stop. Type `halt.` to exit.

```prolog
?- parent(george, X).
X = david ;
X = susan ;
X = mark.

?- ancestor(A, lily).
A = emma ;
A = ryan ;
A = george ;
A = helen ;
A = david ;
A = laura ;
false.

?- sibling(jack, S).
S = emma ;
S = emma ;
false.

?- siblings_of(jack, Siblings).
Siblings = [emma].

?- cousin(lily, C).
false.
```

The query `sibling(jack, S)` gives `emma` two times, once for each shared parent. The helper `siblings_of/2` gives `emma` only once. The query `cousin(lily, C)` returns `false` because the parents of `lily` have no siblings with children.

## Sample output

The demo prints the same output on each run.

```
=== Family Tree in Prolog ===

?- parent(george, david).
true.

?- female(olivia).
true.

?- children_of(david, Children).
Children = [emma, jack].

?- children_of(mark, Children).
Children = [].

?- grandparent(george, emma).
true.

?- setof(G, grandparent(G, lily), Grandparents).
Grandparents = [david, laura].

?- siblings_of(emma, Siblings).
Siblings = [jack].

?- siblings_of(david, Siblings).
Siblings = [mark, susan].

?- siblings_of(olivia, Siblings).
Siblings = [].

?- cousin(emma, olivia).
true.

?- cousin(emma, jack).
false.

?- setof(C, cousin(olivia, C), Cousins).
Cousins = [emma, jack].

?- ancestor(george, noah).
true.

?- descendants_of(emma, Descendants).
Descendants = [lily, noah].

?- descendants_of(george, Descendants).
Descendants = [david, emma, jack, lily, mark, noah, olivia, susan].
```

## Project layout

```
Assignment-8/
  Dockerfile        the official SWI-Prolog image, which runs the demo by default
  README.md
  screenshots/
  prolog/
    family.pl       the facts, the rules, the query helpers, and the demo
```
