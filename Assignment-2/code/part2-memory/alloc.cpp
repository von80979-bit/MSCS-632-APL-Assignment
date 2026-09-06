// Section 3, C++. Allocates integers on the heap in phases, sums them through a const pointer, and frees every
// block with an explicit delete[]. The workload is identical to alloc.rs and alloc.java, so the three checksums
// must match.
//
//   g++ -O2 -o alloc alloc.cpp && ./alloc

#include <chrono>
#include <cstddef>
#include <cstdint>
#include <iomanip>
#include <iostream>
#include <thread>
#include <vector>

constexpr std::size_t INTEGERS_PER_BLOCK = 10000000;
constexpr std::size_t BLOCK_COUNT = 5;

// All three pauses are here for the profiler, not for the algorithm. Without them the run finishes in well
// under a second, psrecord collects a handful of samples, and the chart has no readable shape. The first pause
// gives the chart a baseline before anything is allocated, because psrecord's first sample otherwise lands
// after phase one and the staircase loses its bottom step.
constexpr std::chrono::milliseconds BASELINE_HOLD(1500);
constexpr std::chrono::milliseconds PHASE_HOLD(1500);
constexpr std::chrono::seconds TAIL_HOLD(3);

// Takes a pointer it does not own, so responsibility for the delete[] stays with the caller.
std::int64_t sumBlock(const int* block, std::size_t elementCount) {
    std::int64_t total = 0;
    for (std::size_t elementIndex = 0; elementIndex < elementCount; ++elementIndex) {
        total += block[elementIndex];
    }
    return total;
}

int main() {
    const auto startedAt = std::chrono::steady_clock::now();
    std::int64_t checksum = 0;
    std::vector<int*> blocks;
    blocks.reserve(BLOCK_COUNT);
    std::this_thread::sleep_for(BASELINE_HOLD);

    for (std::size_t blockIndex = 0; blockIndex < BLOCK_COUNT; ++blockIndex) {
        int* block = new int[INTEGERS_PER_BLOCK];
        for (std::size_t elementIndex = 0; elementIndex < INTEGERS_PER_BLOCK; ++elementIndex) {
            block[elementIndex] = static_cast<int>(elementIndex % 1000 + blockIndex);
        }
        checksum += sumBlock(block, INTEGERS_PER_BLOCK);
        blocks.push_back(block);
        std::cout << "phase " << blockIndex + 1 << " allocated " << INTEGERS_PER_BLOCK
                  << " integers, running checksum " << checksum << std::endl;
        std::this_thread::sleep_for(PHASE_HOLD);
    }

    std::cout << "releasing " << blocks.size() << " blocks" << std::endl;
    for (int* block : blocks) {
        delete[] block; // Nothing frees this block unless the program says so, on the line the program chooses.
    }
    blocks.clear();

    std::cout << "checksum " << checksum << std::endl;
    std::this_thread::sleep_for(TAIL_HOLD);
    const std::chrono::duration<double> elapsed = std::chrono::steady_clock::now() - startedAt;
    std::cout << "done in " << std::fixed << std::setprecision(1) << elapsed.count() << " s" << std::endl;
    return 0;
}
