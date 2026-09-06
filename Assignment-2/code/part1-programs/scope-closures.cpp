// Section 2: the same three differences in C++, shown by what this program prints. The three are name binding,
// closure capture and type checking, in that order.
//
// Type checking needs its own file. A type error stops the whole translation unit, so type-check.cpp carries it and
// the two parts below still run.
//
// Run: g++ -std=c++17 -g -O0 -o scope-closures scope-closures.cpp && ./scope-closures

#include <functional>
#include <iostream>
#include <string>
#include <vector>

// Name binding. The compiler ties the name inside this function to the file-scope variable, and a variable of the
// same name declared later in main() cannot change that choice.

int measurementLimit = 10;

int reportMeasurementLimit() { return measurementLimit; }

// Closure capture. The capture list decides, one lambda at a time, and the third case keeps a reference to a
// variable that no longer exists.

std::vector<std::function<int()>> buildIndexReportersByValue() {
    std::vector<std::function<int()>> reporters;
    for (int index = 0; index < 3; ++index) {
        reporters.push_back([index]() { return index; });  // each lambda stores its own copy of index
    }
    return reporters;
}

std::function<int()> buildDanglingReporter() {
    int reportedIndex = 2;
    return [&reportedIndex]() { return reportedIndex; };  // the reference outlives reportedIndex; Part 2 returns here
}

int overwriteStackFrame() {
    int filler[64];  // pushes a live frame over the dead one, so the dangling read finds something other than a 2
    for (int position = 0; position < 64; ++position) { filler[position] = -1; }
    return filler[0];
}

int main() {
    int measurementLimit = 99;  // shadows the file-scope name, but only inside main()
    std::cout << "name binding: main() sees " << measurementLimit << ", the function still reports "
              << reportMeasurementLimit() << "\n";

    std::cout << "closure capture: by value ->";
    for (const auto& reporter : buildIndexReportersByValue()) { std::cout << " " << reporter(); }
    std::cout << "\n";

    int sharedIndex = 0;  // one variable for all three lambdas, so they read whatever it holds when they are called
    std::vector<std::function<int()>> referenceReporters;
    for (int index = 0; index < 3; ++index) {
        sharedIndex = index;
        referenceReporters.push_back([&sharedIndex]() { return sharedIndex; });
    }
    std::cout << "closure capture: by reference ->";
    for (const auto& reporter : referenceReporters) { std::cout << " " << reporter(); }
    std::cout << "\n";

    std::function<int()> danglingReporter = buildDanglingReporter();
    overwriteStackFrame();
    std::cout << "closure capture: dangling reference reads " << danglingReporter()
              << ", which is what this run produced rather than a defined answer\n";

    // Type checking. C++ needs the conversion written out. type-check.cpp shows what the compiler says when it is
    // left out.
    std::cout << "type checking: with the conversion written out -> " << std::string("1") + std::to_string(1) << "\n";
    return 0;
}
