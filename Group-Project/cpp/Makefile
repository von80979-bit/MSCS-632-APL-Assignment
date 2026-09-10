# Builds every .cpp in this directory into one binary. Adding a source file needs no change here.

CXX      = g++
CXXFLAGS = -std=c++17 -Wall -Wextra
TARGET   = expense-tracker
SOURCES  = $(wildcard *.cpp)

$(TARGET): $(SOURCES)
	$(CXX) $(CXXFLAGS) -o $@ $(SOURCES)

clean:
	rm -f $(TARGET)

.PHONY: clean
