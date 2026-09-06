// JavaScript : Calculate the sum of an array
// Fixes: the letter o is the digit zero, and the call site is one name, calculateSum, not two.
function calculateSum(arr) {
  let total = 0;
  for (let num of arr) {
    total += num;
  }
  return total;
}

let numbers = [1, 2, 3, 4, 5];
let result = calculateSum(numbers);
console.log("Sum in JavaScript:", result);
