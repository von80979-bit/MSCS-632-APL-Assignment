# Running this assignment

Every program here runs inside one Ubuntu container. 

Docker Desktop is required to run any command below.

## What each folder holds

| Folder | Holds |
| --- | --- |
| `part1-snippets/` | the three provided broken snippets and three corrected ones, one per language. Section 1 of the report |
| `part1-programs/` | `scope-closures.{py,js,cpp}`, one file per language in three named parts, plus `type-check.cpp`, which fails to compile on purpose. Section 2 |
| `part2-memory/` | `alloc.{rs,java,cpp}`, the measured workload, and `error-demo.{rs,java,cpp}`, one memory mistake per language. Section 3 |
| `profiling/` | psrecord sample logs, the two charts and the scripts that draw them, and the valgrind and GC log transcripts |

`container.sh`, `Dockerfile` and `toolchain-versions.sh` build and enter the environment. `run-all.sh` produces every output; `run-all.log` is the text record of the last run.

## Reproducing every output

```sh
./run-all.sh              # all four sections, in report order, into a fresh run-all.log
./run-all.sh section1     # the three broken snippets, and the three corrected ones
./run-all.sh section2     # scopes and closures in Python, JavaScript and C++
./run-all.sh section3     # what each language does about a memory mistake
./run-all.sh profiling    # the measured runs, the memory summaries, the charts, the clean valgrind runs
```

Add `--pause` to stop before each command and wait for a key, which is how the screenshots were taken one run at a time.

Running one section appends to `run-all.log`. Running the whole script starts the log fresh, so a full run leaves one file holding every output in report order.
