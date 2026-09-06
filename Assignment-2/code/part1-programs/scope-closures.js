// Section 2: the same three differences in JavaScript, shown by what this program prints. The three are name
// binding, closure capture and type checking, in that order.
//
// Run: node scope-closures.js

// Name binding. Both declarations below are known before this line runs, but var and let disagree about what the
// name holds until its own line is reached.

function readHoistedLimit() {
  return hoistedLimit;
}

function readBlockScopedLimit() {
  return blockScopedLimit;
}

console.log("name binding: var before its declaration ->", readHoistedLimit());

try {
  console.log("name binding: let before its declaration ->", readBlockScopedLimit());
} catch (accessBeforeInitialisation) {
  console.log("name binding: ReferenceError ->", accessBeforeInitialisation.message);
}

var hoistedLimit = 10;
let blockScopedLimit = 99;
console.log("name binding: after both declaration ->", readHoistedLimit(), readBlockScopedLimit());

// Closure capture. The same loop gives two different answers, decided by the loop keyword alone.

function buildIndexReportersWithVar() {
  const reporters = [];
  for (var index = 0; index < 3; index += 1) {
    reporters.push(() => index); // one function-scoped index, shared by all three, left holding 3 by the exit test
  }
  return reporters;
}

function buildIndexReportersWithLet() {
  const reporters = [];
  for (let index = 0; index < 3; index += 1) {
    reporters.push(() => index); // let gives each iteration its own binding, so each closure keeps its own index
  }
  return reporters;
}

console.log("closure capture: var ->", buildIndexReportersWithVar().map((reporter) => reporter()));
console.log("closure capture: let ->", buildIndexReportersWithLet().map((reporter) => reporter()));

// Type checking. No error surfaces at all. The operator decides which way to convert, and both lines run.

console.log('type checking: "1" + 1 ->', "1" + 1, "of type", typeof ("1" + 1));
console.log('type checking: "1" - 1 ->', "1" - 1, "of type", typeof ("1" - 1));
