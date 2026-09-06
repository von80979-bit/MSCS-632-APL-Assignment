#!/usr/bin/env bash
# Produces every output the report needs, in the order the report presents them, so the screenshots come from
# real runs rather than from saved recordings. Everything printed is also written to run-all.log, which is the
# text record the report quotes exact error wording from, since a screenshot cannot be copied out of.
#
#   ./run-all.sh              every section, in report order
#   ./run-all.sh section1     the three broken snippets, and the three corrected ones
#   ./run-all.sh section2     scopes and closures across Python, JavaScript and C++
#   ./run-all.sh section3     memory errors in Rust, C++ and Java
#   ./run-all.sh profiling    the measured allocation runs and the charts
#
# Add --pause to stop before each command and wait for a key, which is what makes screenshotting one run at a
# time possible without scrolling.
#
# Running one section appends to the log. Running the whole script starts the log fresh, so a full run leaves
# one file holding every output in report order.
#
# Five steps exit non-zero on purpose: the three broken snippets, the C++ type check and the Rust move. The
# script therefore does not use 'set -e', which would stop it on the first broken snippet and produce nothing.
# The screenshot guide says which those five are; the log shows the exit code and lets the message speak.

set -uo pipefail

CODE_DIRECTORY=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
cd "${CODE_DIRECTORY}"

CONTAINER=./container.sh
LOG_FILE="${CODE_DIRECTORY}/run-all.log"
SEPARATOR="======================================================================"

pause_before_each_run=false

section() {
    { printf '\n%s\n# %s\n' "${SEPARATOR}" "$1"; } | tee -a "${LOG_FILE}"
}

# run [--screen-lines N] <display command> -- <argv for container.sh run>
# Echoes the command the way a reader would type it, runs it in the container, and appends both to the log.
# Output is piped rather than shown on a terminal, which keeps carriage returns out of the file.
#
# --screen-lines keeps the whole output in the log but shows only the first N lines on screen, for the one
# step whose output runs past a screenful and cannot be photographed whole.
run() {
    local screen_lines=0
    if [ "${1:-}" = "--screen-lines" ]; then screen_lines=$2; shift 2; fi

    local display=$1; shift
    [ "${1:-}" = "--" ] && shift

    if [ "${pause_before_each_run}" = true ]; then
        printf '\n  next: %s\n  press return to run it, or ctrl-c to stop.' "${display}"
        read -r _
        clear
    fi

    { printf '\n%s\n$ %s\n' "${SEPARATOR}" "${display}"; } | tee -a "${LOG_FILE}"

    local status
    if [ "${screen_lines}" -gt 0 ]; then
        local output
        output=$("${CONTAINER}" run "$@" 2>&1); status=$?
        printf '%s\n' "${output}" >> "${LOG_FILE}"
        printf '%s\n' "${output}" | head -n "${screen_lines}"
        printf '\n[%d lines in total; the rest is in run-all.log]\n' "$(printf '%s\n' "${output}" | wc -l)"
    else
        "${CONTAINER}" run "$@" 2>&1 | tee -a "${LOG_FILE}"
        status=${PIPESTATUS[0]}
    fi
    printf '[exit %d]\n' "${status}" | tee -a "${LOG_FILE}"
    return 0
}

run_section_1() {
    section "section 1: the three supplied snippets, each broken in its own way"

    run "python3 part1-snippets/calculate-sum-broken.py" \
        -- python3 part1-snippets/calculate-sum-broken.py

    run "node part1-snippets/calculate-sum-broken.js" \
        -- node part1-snippets/calculate-sum-broken.js

    run "g++ -std=c++17 -o part1-snippets/calculate-sum-broken part1-snippets/calculate-sum-broken.cpp" \
        -- g++ -std=c++17 -o part1-snippets/calculate-sum-broken part1-snippets/calculate-sum-broken.cpp

    run "python3 part1-snippets/calculate-sum-fixed.py" -- python3 part1-snippets/calculate-sum-fixed.py
    run "node part1-snippets/calculate-sum-fixed.js" -- node part1-snippets/calculate-sum-fixed.js
    run "g++ -std=c++17 -o /tmp/calculate-sum-fixed part1-snippets/calculate-sum-fixed.cpp && /tmp/calculate-sum-fixed" \
        -- bash -c 'g++ -std=c++17 -o /tmp/calculate-sum-fixed part1-snippets/calculate-sum-fixed.cpp && /tmp/calculate-sum-fixed'
}

run_section_2() {
    section "section 2: scopes and closures in Python, JavaScript and C++"

    run "python3 part1-programs/scope-closures.py" -- python3 part1-programs/scope-closures.py
    run "node part1-programs/scope-closures.js" -- node part1-programs/scope-closures.js
    run "g++ -std=c++17 -g -O0 -o /tmp/scope-closures scope-closures.cpp && /tmp/scope-closures" \
        -- bash -c 'cd part1-programs && g++ -std=c++17 -g -O0 -o /tmp/scope-closures scope-closures.cpp && /tmp/scope-closures'

    # g++ follows the error with every operator+ candidate it rejected, which is why only the head is shown.
    run --screen-lines 7 "g++ -std=c++17 -o /tmp/type-check type-check.cpp" \
        -- bash -c 'cd part1-programs && g++ -std=c++17 -o /tmp/type-check type-check.cpp'
}


run_section_3() {
    section "section 3: what each language does about a memory mistake"

    # Each step compiles and runs in one invocation, because every container.sh run is a fresh container and a
    # binary written to /tmp by one step is gone by the next.

    run "rustc -o /tmp/error-demo-rs part2-memory/error-demo.rs" \
        -- rustc -o /tmp/error-demo-rs part2-memory/error-demo.rs

    # The Java heap is pinned in every Java run, so the figures do not depend on the machine's memory.
    run "java -Xms32m -Xmx512m -Xlog:gc part2-memory/error-demo.java" \
        -- java -Xms32m -Xmx512m -Xlog:gc part2-memory/error-demo.java

    # -g keeps the source lines for valgrind and -O0 keeps the optimiser from deleting the leaking allocation.
    run "g++ -g -O0 -o /tmp/error-demo part2-memory/error-demo.cpp && /tmp/error-demo" \
        -- bash -c 'g++ -g -O0 -o /tmp/error-demo part2-memory/error-demo.cpp && /tmp/error-demo'

    run "valgrind --leak-check=full /tmp/error-demo" \
        -- bash -c 'g++ -g -O0 -o /tmp/error-demo part2-memory/error-demo.cpp && valgrind --leak-check=full /tmp/error-demo'
}

run_profiling() {
    section "profiling: the same workload in three languages, measured the same way"

    run "g++ -O2 -o part2-memory/alloc-cpp part2-memory/alloc.cpp" \
        -- g++ -O2 -o part2-memory/alloc-cpp part2-memory/alloc.cpp
    run "rustc -O -o part2-memory/alloc-rust part2-memory/alloc.rs" \
        -- rustc -O -o part2-memory/alloc-rust part2-memory/alloc.rs
    run "javac -d part2-memory part2-memory/alloc.java" \
        -- javac -d part2-memory part2-memory/alloc.java

    # psrecord writes its readings to a file and prints nothing to the terminal, so each profiled run is
    # followed by memory-summary.py. Without it a screenshot of these runs carries no memory figure at all.
    # --include-children is required: without it psrecord follows the wrapper and reports a flat 1.645 MB.
    # The runs go one at a time, because a second profiled workload perturbs the first.
    profile_one_run "./part2-memory/alloc-cpp" cpp
    profile_one_run "./part2-memory/alloc-rust" rust
    profile_one_run "java -Xms32m -Xmx512m -Xlog:gc -cp part2-memory Alloc" java
    profile_one_run "java -Xms32m -Xmx512m -XX:G1PeriodicGCInterval=1000 -Xlog:gc -cp part2-memory Alloc" java-periodic-gc

    run "python3 profiling/combined-memory.py" -- python3 profiling/combined-memory.py
    run "python3 profiling/java-gc-contrast.py" -- python3 profiling/java-gc-contrast.py

    # The same tool and flags on the measured programs. Without these two the leak in section 3 cannot be told
    # apart from a property of C++ or of valgrind. alloc.rs is built -O because memcheck runs the full 200 MB.
    run "valgrind --leak-check=full /tmp/alloc-rust" \
        -- bash -c 'rustc -g -O -o /tmp/alloc-rust part2-memory/alloc.rs && valgrind --leak-check=full /tmp/alloc-rust'
    run "valgrind --leak-check=full /tmp/alloc-cpp" \
        -- bash -c 'g++ -g -O0 -o /tmp/alloc-cpp part2-memory/alloc.cpp && valgrind --leak-check=full /tmp/alloc-cpp'
}

# profile_one_run <command psrecord watches> <name used for the log and the chart>
profile_one_run() {
    local watched_command=$1 name=$2
    local display
    printf -v display 'psrecord "%s" \\\n    --interval 0.1 --include-children --log profiling/%s.txt --plot profiling/%s.png \\\n    && python3 profiling/memory-summary.py profiling/%s.txt' \
        "${watched_command}" "${name}" "${name}" "${name}"

    run "${display}" -- bash -c "psrecord '${watched_command}' --interval 0.1 --include-children \
        --log profiling/${name}.txt --plot profiling/${name}.png \
        && python3 profiling/memory-summary.py profiling/${name}.txt"
}

not_written_yet() {
    printf '\n%s is not written yet.\n' "$1"
}

subcommand=all
for argument in "$@"; do
    case "${argument}" in
        --pause) pause_before_each_run=true ;;
        *)       subcommand=${argument} ;;
    esac
done

case "${subcommand}" in
    section1|section2|section3|profiling) ;;
    all) : > "${LOG_FILE}" ;;
    *)
        echo "unknown section '${subcommand}'; expected section1, section2, section3, profiling or all" >&2
        exit 2
        ;;
esac

case "${subcommand}" in
    section1)  run_section_1 ;;
    section2)  run_section_2 ;;
    section3)  run_section_3 ;;
    profiling) run_profiling ;;
    all)
        run_section_1
        run_section_2
        run_section_3
        run_profiling
        ;;
esac

printf '\nlog: %s\n' "${LOG_FILE}"
