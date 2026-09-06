"""Prints the four figures that matter from a psrecord sample log.

psrecord writes its readings to a file and prints nothing to the terminal, so a screenshot of a profiled run
carries no memory figure at all. This turns the log back into four lines a reader can see.

    python3 profiling/memory-summary.py profiling/java.txt
"""

import sys

BASELINE_HOLD_ENDS_AT_SECONDS = 1.4    # the allocation programs hold for 1.5 s before the first phase


def read_samples(log_path):
    samples = []
    with open(log_path) as log_file:
        for line in log_file:
            if line.startswith("#"):
                continue
            columns = line.split()
            if len(columns) >= 3:
                samples.append((float(columns[0]), float(columns[2])))
    return samples


def summarise(log_path):
    samples = read_samples(log_path)
    if not samples:
        print(f"{log_path}: no samples")
        return

    baseline = [reading for elapsed, reading in samples if elapsed <= BASELINE_HOLD_ENDS_AT_SECONDS]
    peak_elapsed, peak_reading = max(samples, key=lambda sample: sample[1])
    last_elapsed, last_reading = samples[-1]

    # The largest fall between two consecutive samples is what separates the languages: Rust and C++ give
    # nearly everything back in one sampling interval, and Java gives back nothing at all.
    falls = [(before[1] - after[1], before[0], after[0]) for before, after in zip(samples, samples[1:])]
    largest_fall, fall_from, fall_to = max(falls, key=lambda fall: fall[0])

    print(f"  memory summary of {log_path}, {len(samples)} samples")
    print(f"    baseline, before the first phase   {max(baseline) if baseline else float('nan'):7.1f} MB")
    print(f"    peak                               {peak_reading:7.1f} MB at {peak_elapsed:.1f} s")
    print(f"    largest fall in one interval       {largest_fall:7.1f} MB between {fall_from:.1f} s and {fall_to:.1f} s")
    print(f"    last sample                        {last_reading:7.1f} MB at {last_elapsed:.1f} s")


if __name__ == "__main__":
    if len(sys.argv) != 2:
        sys.exit("usage: memory-summary.py <psrecord log>")
    summarise(sys.argv[1])
