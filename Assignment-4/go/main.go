// Reads the slot maximum, runs the input, the scheduler, and the output.

package main

import (
	"fmt"
	"math/rand/v2"
	"os"
)

func main() {
	maxPerSlot, isValid := slotMaximum(os.LookupEnv("MAX_EMPLOYEES_PER_SLOT"))
	if !isValid {
		fmt.Fprintf(os.Stderr, "MAX_EMPLOYEES_PER_SLOT must be a whole number of 2 or more. Using %d.\n", maxPerSlot)
	}

	employees := readEmployees(os.Stdin, os.Stdout)
	if len(employees) == 0 {
		fmt.Println("No employees entered.")
		return
	}

	random := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	schedule := buildSchedule(employees, maxPerSlot, random)
	fmt.Println()
	writeSchedule(os.Stdout, schedule)
}
