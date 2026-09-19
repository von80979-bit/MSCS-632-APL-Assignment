// Reads the slot maximum, runs the input, the scheduler, and the output.

package main

import (
	"fmt"
	"math/rand/v2"
	"os"
)

func main() {
	slotMaximum, isValid := parseSlotMaximum(os.LookupEnv(slotMaximumVariable))
	if !isValid {
		fmt.Fprintf(os.Stderr, "%s must be a whole number of %d or more. Using %d.\n",
			slotMaximumVariable, minEmployeesPerSlot, slotMaximum)
	}

	employees := readEmployees(os.Stdin, os.Stdout)
	if len(employees) == 0 {
		fmt.Println("No employees entered.")
		return
	}

	random := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	schedule := buildSchedule(employees, slotMaximum, random)
	fmt.Println()
	writeSchedule(os.Stdout, schedule)
}
