#!/usr/bin/env python3
"""Draws the three psrecord logs as one chart.

psrecord writes one PNG per run, each on its own vertical scale, so three separate images would hide the very
difference this assignment measures. This script puts all three runs on one pair of axes.

    ./container.sh run python3 profiling/combined-memory.py

Run it after cpp.txt, rust.txt and java.txt exist.
"""

import matplotlib

matplotlib.use("Agg")  # no display inside the container

import matplotlib.pyplot as plt

PROFILING_DIRECTORY = "profiling"
OUTPUT_IMAGE = f"{PROFILING_DIRECTORY}/combined-memory.png"

# Colour and label per language, in the order the legend should read them. The Java heap is pinned, so the
# legend names the two options rather than leaving a reader to guess which heap produced the line.
LANGUAGE_SERIES = [
    ("rust", "Rust (rustc 1.75.0)", "#c1440e"),
    ("cpp", "C++ (g++ 13.3.0)", "#1f6fb2"),
    ("java", "Java (OpenJDK 21, G1, -Xms32m -Xmx512m)", "#2e7d32"),
]


def read_psrecord_log(log_path):
    """Returns elapsed seconds and resident megabytes from a psrecord log, skipping its header line."""
    elapsed_seconds = []
    resident_megabytes = []
    with open(log_path) as log_file:
        for line in log_file:
            if line.startswith("#"):
                continue
            elapsed, _cpu_percent, resident, _virtual = (float(field) for field in line.split())
            elapsed_seconds.append(elapsed)
            resident_megabytes.append(resident)
    return elapsed_seconds, resident_megabytes


def draw_combined_chart():
    figure, axes = plt.subplots(figsize=(9.5, 5.5), dpi=200)

    for language, legend_label, colour in LANGUAGE_SERIES:
        elapsed_seconds, resident_megabytes = read_psrecord_log(f"{PROFILING_DIRECTORY}/{language}.txt")
        peak_megabytes = max(resident_megabytes)
        axes.plot(elapsed_seconds, resident_megabytes, color=colour, linewidth=2.2,
                  label=f"{legend_label} - peak {peak_megabytes:.1f} MB")

    axes.set_title("Resident memory of the same 200 MB workload in three languages", fontsize=14, pad=12)
    axes.set_xlabel("Elapsed time (seconds)", fontsize=12)
    axes.set_ylabel("Resident memory (MB)", fontsize=12)
    axes.tick_params(labelsize=11)
    axes.set_xlim(0, 13.2)
    axes.set_ylim(0, 292)
    axes.grid(True, linewidth=0.5, alpha=0.4)
    axes.legend(loc="upper left", fontsize=11, framealpha=0.95)

    # The tail is the finding, so it is labelled on the chart itself. Both labels sit in bands the lines leave
    # empty, which keeps them readable at page width without arrows crossing the staircase.
    axes.text(0.4, 216, "Java holds about 245 MB to exit", fontsize=11)
    axes.text(0.4, 152, "Rust and C++ fall back to about\n5 MB when the blocks are released", fontsize=11)

    figure.tight_layout()
    figure.savefig(OUTPUT_IMAGE)
    print(f"wrote {OUTPUT_IMAGE}")


if __name__ == "__main__":
    draw_combined_chart()
