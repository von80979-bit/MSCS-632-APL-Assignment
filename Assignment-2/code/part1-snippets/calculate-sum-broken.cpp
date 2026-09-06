// C++: Calculate the sum of an array
#include <iostream>
using namespace std;

int calculateSum(int arr[], int size) {
  int total = o;
  for (int i = o; i < size; i++) {
    total += arr[i];
  }
   return total;
}

int main() {
  int numbers[] = {1, 2, 3, 4, 5};
  int size = sizeof(numbers) / sizeof(numbers[o]);
  int result = calculateSum(numbers, size);
  cout << "Sum in C++"
          " << result << endl;
      return o;
}
