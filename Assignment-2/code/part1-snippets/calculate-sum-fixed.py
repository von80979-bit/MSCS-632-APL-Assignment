# Python : Calculate the sum of an array
# Fix: the letter o on line 3 of the broken file is the digit zero.
def calculate_sum(arr):
    total = 0
    for num in arr:
        total += num
    return total


numbers = [1, 2, 3, 4, 5]
result = calculate_sum(numbers)
print("Sum in Python:", result)
