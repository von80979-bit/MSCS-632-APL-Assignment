// Type checking, the third difference, in the one form C++ can show it. This file exists to fail, and its result is
// the compiler's message rather than a program.
//
// std::string is deliberate. A bare "1" + 1 is pointer arithmetic on a const char*, which g++ compiles silently
// even under -Wall -Wextra and then prints an empty string, so it would say nothing about type checking. C++ has two
// different things that look like strings, and only one of them is checked.
//
// Run: g++ -std=c++17 -o type-check type-check.cpp

#include <iostream>
#include <string>

int main() {
    std::string digitAsText = "1";
    std::cout << digitAsText + 1 << "\n";  // no operator+ takes a std::string and an int
    return 0;
}
