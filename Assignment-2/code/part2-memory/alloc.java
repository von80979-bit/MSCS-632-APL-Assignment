// Section 3, Java. Allocates integers on the heap in phases, sums them, and drops the reference so the
// collector decides when to reclaim. The workload is identical to alloc.rs and alloc.cpp, so the three
// checksums must match. The class is package private, which lets javac accept the lower-case file name.
//
//   javac alloc.java && java -Xlog:gc Alloc

import java.util.ArrayList;
import java.util.List;

class Alloc {
    private static final int INTEGERS_PER_BLOCK = 10_000_000;
    private static final int BLOCK_COUNT = 5;

    // All three pauses are here for the profiler, not for the algorithm. Without them the run finishes in well
    // under a second, psrecord collects a handful of samples, and the chart has no readable shape. The first
    // pause gives the chart a baseline before anything is allocated, because psrecord's first sample otherwise
    // lands after phase one and the staircase loses its bottom step.
    private static final long BASELINE_HOLD_MILLISECONDS = 1500;
    private static final long PHASE_HOLD_MILLISECONDS = 1500;
    private static final long TAIL_HOLD_MILLISECONDS = 3000;

    // Reads the block without storing it anywhere, so this call adds no reference that could delay collection.
    private static long sumBlock(int[] block) {
        long total = 0;
        for (int value : block) {
            total += value;
        }
        return total;
    }

    public static void main(String[] arguments) throws InterruptedException {
        long startedAt = System.nanoTime();
        long checksum = 0;
        List<int[]> blocks = new ArrayList<>(BLOCK_COUNT);
        Thread.sleep(BASELINE_HOLD_MILLISECONDS);

        for (int blockIndex = 0; blockIndex < BLOCK_COUNT; blockIndex++) {
            int[] block = new int[INTEGERS_PER_BLOCK];
            for (int elementIndex = 0; elementIndex < INTEGERS_PER_BLOCK; elementIndex++) {
                block[elementIndex] = elementIndex % 1000 + blockIndex;
            }
            checksum += sumBlock(block);
            blocks.add(block);
            System.out.println("phase " + (blockIndex + 1) + " allocated " + INTEGERS_PER_BLOCK
                    + " integers, running checksum " + checksum);
            Thread.sleep(PHASE_HOLD_MILLISECONDS);
        }

        System.out.println("releasing " + blocks.size() + " blocks");
        blocks = null; // The whole release. Nothing forces a collection here; the collector chooses its moment.

        System.out.println("checksum " + checksum);
        Thread.sleep(TAIL_HOLD_MILLISECONDS);
        System.out.printf("done in %.1f s%n", (System.nanoTime() - startedAt) / 1_000_000_000.0);
    }
}
