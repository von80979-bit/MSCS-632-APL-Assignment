// Reads the flags, runs the demo in section 7 of the specification, and sets the exit status.

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

const (
	tripDelayPerMile = 100 * time.Millisecond
	usageLine        = "usage: rideshare-go [-drivers 1-8] [-out path] [-timeout seconds]"
)

var driverNames = []string{
	"Ana Lopez", "Ben Carter", "Chen Wei", "Dana Kim", "Eli Novak", "Fatima Ali", "Gus Moreno", "Hana Sato",
}

func main() {
	config, err := parseFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, usageLine)
		os.Exit(2)
	}

	fmt.Println("=== Concurrent Ride Sharing System ===")
	summary := RunSystem(config, demoRides(), NewLogger(os.Stdout))
	summary.Print(os.Stdout)
	os.Exit(summary.ExitStatus())
}

func parseFlags(args []string) (Config, error) {
	flags := flag.NewFlagSet("rideshare-go", flag.ContinueOnError)
	flags.SetOutput(io.Discard) // main prints one usage line instead of the flag package defaults
	driverCount := flags.Int("drivers", 3, "number of drivers, from 1 to 8")
	resultsPath := flags.String("out", "results.txt", "path of the results file")
	timeoutSeconds := flags.Int("timeout", 30, "time limit in seconds for the drivers after the queue closes")
	if err := flags.Parse(args); err != nil {
		return Config{}, err
	}
	switch {
	case flags.NArg() > 0:
		return Config{}, fmt.Errorf("unexpected argument %q", flags.Arg(0))
	case *driverCount < 1 || *driverCount > len(driverNames):
		return Config{}, fmt.Errorf("-drivers must be from 1 to %d", len(driverNames))
	case *timeoutSeconds < 1:
		return Config{}, errors.New("-timeout must be 1 or more")
	}
	return Config{
		DriverNames:      driverNames[:*driverCount],
		ResultsPath:      *resultsPath,
		TripDelayPerMile: tripDelayPerMile,
		Timeout:          time.Duration(*timeoutSeconds) * time.Second,
	}, nil
}

// demoRides gives the rides in queue order. The fields follow the table in section 7: ID, ride type, pickup, dropoff,
// miles, rider, and passengers.
func demoRides() []Ride {
	return []Ride{
		{"RIDE-1", Standard, "Downtown", "Airport", 12.0, "Maria Garcia", 1},
		{"RIDE-2", Premium, "Hotel", "Convention Center", 4.5, "James Smith", 1},
		{"RIDE-8", Standard, "Office", "Office", 0.0, "James Smith", 1},
		{"RIDE-3", Shared, "University", "Stadium", 6.0, "Maria Garcia", 2},
		{"RIDE-4", Standard, "Mall", "Train Station", 3.2, "James Smith", 1},
		{"RIDE-9", Shared, "Pier", "Market", 5.0, "Maria Garcia", 1},
		{"RIDE-5", Premium, "Airport", "Harbor", 15.0, "Maria Garcia", 1},
		{"RIDE-6", Shared, "Library", "Park", 8.0, "James Smith", 3},
		{"RIDE-7", Standard, "Museum", "Zoo", 2.0, "Maria Garcia", 1},
	}
}
