"""Statistics calculator in object-oriented Python: a class keeps the data and has one method for each statistic.

Run: python3 stats.py
"""


class StatisticsCalculator:
    def __init__(self, values):
        # A copy keeps the calculator independent of later changes to the caller's list.
        self._values = list(values)

    def values(self):
        return list(self._values)

    def count(self):
        return len(self._values)

    def mean(self):
        return sum(self._values) / self.count()

    def median(self):
        ordered = sorted(self._values)
        middle = self.count() // 2
        if self.count() % 2 == 1:
            return float(ordered[middle])
        return (ordered[middle - 1] + ordered[middle]) / 2

    def mode(self):
        occurrences = {}
        for value in self._values:
            occurrences[value] = occurrences.get(value, 0) + 1
        highest = max(occurrences.values())
        return sorted(value for value, times in occurrences.items() if times == highest)

    def min(self):
        return min(self._values)

    def max(self):
        return max(self._values)

    def range(self):
        return self.max() - self.min()


def format_list(values):
    return "[" + ", ".join(str(value) for value in values) + "]"


def print_report(list_number, calculator):
    print(f"List {list_number}: {format_list(calculator.values())}")
    print(f"  Count:  {calculator.count()}")
    print(f"  Mean:   {calculator.mean():.2f}")
    print(f"  Median: {calculator.median():.2f}")
    print(f"  Mode:   {format_list(calculator.mode())}")
    print(f"  Min:    {calculator.min()}")
    print(f"  Max:    {calculator.max()}")
    print(f"  Range:  {calculator.range()}")


DEMO_LISTS = [
    [4, 8, 15, 16, 23, 42, 8],
    [7, 3, 9, 3, 7, 1, 12, 2],
    [5, -2, 10, 0, -7],
]


def main():
    print("=== Statistics Calculator: Python (object-oriented) ===")
    for list_number, values in enumerate(DEMO_LISTS, start=1):
        print()
        print_report(list_number, StatisticsCalculator(values))


if __name__ == "__main__":
    main()
