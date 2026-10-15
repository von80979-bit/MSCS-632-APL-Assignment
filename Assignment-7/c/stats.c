/*
 * Statistics calculator in procedural C: plain functions over an int array, with manual memory management.
 * Build: gcc -std=c11 -Wall -Wextra -O2 -o stats stats.c
 */
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static void exit_out_of_memory(void) {
    fprintf(stderr, "error: out of memory\n");
    exit(EXIT_FAILURE);
}

static int *copy_array(const int *values, size_t count) {
    int *copy = malloc(count * sizeof *copy);
    if (copy == NULL) {
        exit_out_of_memory();
    }
    memcpy(copy, values, count * sizeof *copy);
    return copy;
}

/* Sorts in place, so callers pass a copy to keep the demo list in its original order. */
static void insertion_sort(int *values, size_t count) {
    for (size_t i = 1; i < count; i++) {
        int current = values[i];
        size_t j = i;
        while (j > 0 && values[j - 1] > current) {
            values[j] = values[j - 1];
            j--;
        }
        values[j] = current;
    }
}

static size_t count_values(size_t count) {
    return count;
}

static double mean(const int *values, size_t count) {
    long long sum = 0;
    for (size_t i = 0; i < count; i++) {
        sum += values[i];
    }
    return (double)sum / (double)count;
}

static double median(const int *values, size_t count) {
    int *sorted = copy_array(values, count);
    insertion_sort(sorted, count);
    double result;
    if (count % 2 == 1) {
        result = sorted[count / 2];
    } else {
        result = (sorted[count / 2 - 1] + sorted[count / 2]) / 2.0;
    }
    free(sorted);
    return result;
}

/* Returns a new array of every value that ties for the longest run in the sorted copy. The caller frees it. */
static int *mode(const int *values, size_t count, size_t *mode_count) {
    int *sorted = copy_array(values, count);
    insertion_sort(sorted, count);

    size_t longest_run = 0;
    for (size_t start = 0; start < count;) {
        size_t end = start;
        while (end < count && sorted[end] == sorted[start]) {
            end++;
        }
        if (end - start > longest_run) {
            longest_run = end - start;
        }
        start = end;
    }

    int *modes = malloc(count * sizeof *modes);
    if (modes == NULL) {
        free(sorted);
        exit_out_of_memory();
    }
    *mode_count = 0;
    for (size_t start = 0; start < count;) {
        size_t end = start;
        while (end < count && sorted[end] == sorted[start]) {
            end++;
        }
        if (end - start == longest_run) {
            modes[(*mode_count)++] = sorted[start];
        }
        start = end;
    }

    free(sorted);
    return modes;
}

static int min_value(const int *values, size_t count) {
    int smallest = values[0];
    for (size_t i = 1; i < count; i++) {
        if (values[i] < smallest) {
            smallest = values[i];
        }
    }
    return smallest;
}

static int max_value(const int *values, size_t count) {
    int largest = values[0];
    for (size_t i = 1; i < count; i++) {
        if (values[i] > largest) {
            largest = values[i];
        }
    }
    return largest;
}

static int range(const int *values, size_t count) {
    return max_value(values, count) - min_value(values, count);
}

static void print_list(const int *values, size_t count) {
    printf("[");
    for (size_t i = 0; i < count; i++) {
        printf("%s%d", i == 0 ? "" : ", ", values[i]);
    }
    printf("]");
}

static void print_report(int list_number, const int *demo_list, size_t count) {
    int *values = copy_array(demo_list, count);
    size_t mode_count;
    int *modes = mode(values, count, &mode_count);

    printf("List %d: ", list_number);
    print_list(values, count);
    printf("\n");
    printf("  Count:  %zu\n", count_values(count));
    printf("  Mean:   %.2f\n", mean(values, count));
    printf("  Median: %.2f\n", median(values, count));
    printf("  Mode:   ");
    print_list(modes, mode_count);
    printf("\n");
    printf("  Min:    %d\n", min_value(values, count));
    printf("  Max:    %d\n", max_value(values, count));
    printf("  Range:  %d\n", range(values, count));

    free(modes);
    free(values);
}

int main(void) {
    const int list_1[] = {4, 8, 15, 16, 23, 42, 8};
    const int list_2[] = {7, 3, 9, 3, 7, 1, 12, 2};
    const int list_3[] = {5, -2, 10, 0, -7};

    printf("=== Statistics Calculator: C (procedural) ===\n");
    printf("\n");
    print_report(1, list_1, sizeof list_1 / sizeof list_1[0]);
    printf("\n");
    print_report(2, list_2, sizeof list_2 / sizeof list_2[0]);
    printf("\n");
    print_report(3, list_3, sizeof list_3 / sizeof list_3[0]);
    return EXIT_SUCCESS;
}
