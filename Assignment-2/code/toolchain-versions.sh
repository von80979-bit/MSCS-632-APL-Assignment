#!/usr/bin/env bash
# Prints one line per toolchain this assignment depends on, plus the two checks macOS could not satisfy.
# The report's environment sentence is written from this output, so a missing tool has to be visible here
# rather than surfacing later as a broken ticket.

set -uo pipefail

print_labelled_line() {
    printf '  %-12s %s\n' "$1" "$2"
}

# Most of these tools print several lines and only the first names the version.
report_version() {
    local label=$1
    shift
    local output
    if output=$("$@" 2>&1); then
        print_labelled_line "$label" "$(printf '%s\n' "$output" | head -n 1)"
    else
        print_labelled_line "$label" "MISSING: '$*' failed"
    fi
}

report_python_package_version() {
    local label=$1
    local package=$2
    report_version "$label" python3 -c "import importlib.metadata as metadata; print(metadata.version('${package}'))"
}

# valgrind and psrecord are the reason this container exists, so both are exercised rather than only queried.
report_command_runs() {
    local label=$1
    shift
    if "$@" >/dev/null 2>&1; then
        print_labelled_line "$label" "ok"
    else
        print_labelled_line "$label" "FAILED: '$*'"
    fi
}

echo 'Environment'
print_labelled_line 'distribution' "$(. /etc/os-release && printf '%s' "$PRETTY_NAME")"
print_labelled_line 'kernel' "$(uname -sr)"
print_labelled_line 'architecture' "$(uname -m)"

echo
echo 'Toolchains'
report_version 'g++' g++ --version
report_version 'rustc' rustc --version
report_version 'javac' javac --version
report_version 'java' java --version
report_version 'python3' python3 --version
report_version 'node' node --version

echo
echo 'Profiling'
report_version 'valgrind' valgrind --version
report_python_package_version 'psrecord' psrecord
report_python_package_version 'matplotlib' matplotlib

echo
echo 'Checks'
report_command_runs 'valgrind' valgrind --version
report_command_runs 'psrecord' psrecord --help
report_command_runs 'matplotlib' python3 -c 'import matplotlib; matplotlib.use("Agg")'
