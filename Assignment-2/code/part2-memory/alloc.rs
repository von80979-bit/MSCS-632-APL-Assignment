// Section 3, Rust. Allocates integers on the heap in phases, sums them through a borrow, and lets ownership
// release them. The workload is identical to alloc.java and alloc.cpp, so the three checksums must match.
//
//   rustc -O -o alloc alloc.rs && ./alloc

use std::thread::sleep;
use std::time::{Duration, Instant};

const INTEGERS_PER_BLOCK: usize = 10_000_000;
const BLOCK_COUNT: usize = 5;

// All three pauses are here for the profiler, not for the algorithm. Without them the run finishes in well
// under a second, psrecord collects a handful of samples, and the chart has no readable shape. The first pause
// gives the chart a baseline before anything is allocated, because psrecord's first sample otherwise lands
// after phase one and the staircase loses its bottom step.
const BASELINE_HOLD: Duration = Duration::from_millis(1500);
const PHASE_HOLD: Duration = Duration::from_millis(1500);
const TAIL_HOLD: Duration = Duration::from_secs(3);

// Borrows the block instead of taking it, so the caller keeps ownership and nothing is copied or freed here.
fn sum_block(block: &[i32]) -> i64 {
    block.iter().map(|&value| value as i64).sum()
}

fn main() {
    let started_at = Instant::now();
    let mut checksum: i64 = 0;
    sleep(BASELINE_HOLD);

    // Inner scope: `blocks` owns every allocation, and the compiler inserts the deallocation where the scope ends.
    {
        let mut blocks: Vec<Vec<i32>> = Vec::with_capacity(BLOCK_COUNT);
        for block_index in 0..BLOCK_COUNT {
            let mut block: Vec<i32> = Vec::with_capacity(INTEGERS_PER_BLOCK);
            for element_index in 0..INTEGERS_PER_BLOCK {
                block.push((element_index % 1000) as i32 + block_index as i32);
            }
            checksum += sum_block(&block);
            blocks.push(block);
            let phase = block_index + 1;
            println!("phase {} allocated {} integers, running checksum {}", phase, INTEGERS_PER_BLOCK, checksum);
            sleep(PHASE_HOLD);
        }
        println!("releasing {} blocks", blocks.len());
    }

    println!("checksum {}", checksum);
    sleep(TAIL_HOLD);
    println!("done in {:.1} s", started_at.elapsed().as_secs_f64());
}
