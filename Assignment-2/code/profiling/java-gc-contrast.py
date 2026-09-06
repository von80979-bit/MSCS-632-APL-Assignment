#!/usr/bin/env python3
"""Draws the two Java runs as one chart, so the effect of a single option can be read off it.

The combined chart compares three languages. This one compares Java against itself: same program, same pinned
heap, and one extra option, -XX:G1PeriodicGCInterval=1000. It exists because the three-language chart alone
invites the wrong reading, that the JVM cannot give the memory back.

    ./container.sh run python3 profiling/java-gc-contrast.py

Run it after java.txt and java-periodic-gc.txt exist.
"""

import matplotlib

matplotlib.use("Agg")  # no display inside the container

import matplotlib.pyplot as plt

PROFILING_DIRECTORY = "profiling"
OUTPUT_IMAGE = f"{PROFILING_DIRECTORY}/java-periodic-gc-contrast.png"

# Both runs are the same green, because they are the same language and the same workload. The option is carried
# by the line style rather than by the colour, which is what makes the pair read as one comparison.
JAVA_RUNS = [
    ("java", "No periodic collection asked for, the measured run", "-"),
    ("java-periodic-gc", "-XX:G1PeriodicGCInterval=1000, the one difference", "--"),
]
JAVA_COLOUR = "#2e7d32"


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


def draw_contrast_chart():
    figure, axes = plt.subplots(figsize=(9.5, 5.5), dpi=200)

    for run_name, legend_label, line_style in JAVA_RUNS:
        elapsed_seconds, resident_megabytes = read_psrecord_log(f"{PROFILING_DIRECTORY}/{run_name}.txt")
        peak_megabytes = max(resident_megabytes)
        final_megabytes = resident_megabytes[-1]
        axes.plot(elapsed_seconds, resident_megabytes, color=JAVA_COLOUR, linewidth=2.2, linestyle=line_style,
                  label=f"{legend_label} - peak {peak_megabytes:.1f} MB, {final_megabytes:.1f} MB at exit")

    axes.set_title("The same Java program and the same pinned heap, changed by one option", fontsize=14, pad=12)
    axes.set_xlabel("Elapsed time (seconds)", fontsize=12)
    axes.set_ylabel("Resident memory (MB)", fontsize=12)
    axes.tick_params(labelsize=11)
    # The x limit matches combined-memory.png. The y limit is higher only to hold the legend clear of the lines,
    # since both runs are the same tall Java staircase and neither leaves a wide empty band.
    axes.set_xlim(0, 13.2)
    axes.set_ylim(0, 335)
    axes.grid(True, linewidth=0.5, alpha=0.4)
    axes.legend(loc="upper left", fontsize=10.5, framealpha=0.95)

    # The moment the memory comes back is the whole point of the chart, so the arrow points at it and the label
    # names the collection that caused it. The GC log timestamp is 10.061 s and the fall is the next sample pair.
    axes.annotate("G1 Periodic Collection at 10.061 s,\n197M to 1M in a 2.279 ms pause:\n258.1 MB falls to 50.9 MB"
                  " in one sample",
                  xy=(10.14, 150), xytext=(3.9, 72), fontsize=11,
                  arrowprops=dict(arrowstyle="->", color=JAVA_COLOUR, linewidth=1.4))
    axes.text(0.3, 268, "Without the option the memory is never returned before exit", fontsize=11)

    figure.tight_layout()
    figure.savefig(OUTPUT_IMAGE)
    print(f"wrote {OUTPUT_IMAGE}")


if __name__ == "__main__":
    draw_contrast_chart()
