// Two mistakes C++ allows: a block that is never released, and a read through a pointer whose memory has already
// gone back to the allocator. Compile with -g -O0, otherwise valgrind cannot name the line that allocated.
#include <iostream>

void leakABlock() {
    int *forgottenValues = new int[1000];  // the matching delete[] is missing on purpose
    forgottenValues[0] = 42;
}

int main() {
    leakABlock();

    int *releasedValue = new int(7);
    delete releasedValue;
    // Undefined behaviour. The number printed is what this run produced, not a value the language defines.
    std::cout << "read through the released pointer: " << *releasedValue << std::endl;
    return 0;
}
