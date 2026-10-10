// Test cases 3, 4, 5, 6, 10, and 11 of the specification: full runs of the system.

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const demoTotalCents = 10843

func testConfig(t *testing.T, driverCount int) Config {
	t.Helper()
	return Config{
		DriverNames:      driverNames[:driverCount],
		ResultsPath:      filepath.Join(t.TempDir(), "results.txt"),
		TripDelayPerMile: 0,
		Timeout:          10 * time.Second,
	}
}

func runQuietly(config Config, rides []Ride) Summary {
	return RunSystem(config, rides, NewLogger(io.Discard))
}

// fileRideIDs returns the ride ID of each line in the results file, after the header.
func fileRideIDs(t *testing.T, path string) []string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read results file: %v", err)
	}
	lines := strings.Split(strings.TrimSuffix(string(content), "\n"), "\n")
	if lines[0] != strings.TrimSuffix(resultsHeader, "\n") {
		t.Fatalf("results file header = %q", lines[0])
	}
	ids := make([]string, 0, len(lines)-1)
	for _, line := range lines[1:] {
		ids = append(ids, strings.SplitN(line, ",", 2)[0])
	}
	return ids
}

func resultRideIDs(summary Summary) []string {
	ids := make([]string, len(summary.Results))
	for i, result := range summary.Results {
		ids[i] = result.Ride.ID
	}
	return ids
}

// assertEachOnce checks that ids holds every ID in want exactly once, and nothing else.
func assertEachOnce(t *testing.T, source string, ids []string, want []string) {
	t.Helper()
	counts := map[string]int{}
	for _, id := range ids {
		counts[id]++
	}
	for _, id := range want {
		if counts[id] != 1 {
			t.Errorf("%s has %s %d times, want once", source, id, counts[id])
		}
		delete(counts, id)
	}
	for id := range counts {
		t.Errorf("%s has unexpected ride %s", source, id)
	}
}

var demoCompletedIDs = []string{"RIDE-1", "RIDE-2", "RIDE-3", "RIDE-4", "RIDE-5", "RIDE-6", "RIDE-7"}

func TestEachRideRecordedExactlyOnce(t *testing.T) {
	config := testConfig(t, 3)

	summary := runQuietly(config, demoRides())

	assertEachOnce(t, "results file", fileRideIDs(t, config.ResultsPath), demoCompletedIDs)
	assertEachOnce(t, "results list", resultRideIDs(summary), demoCompletedIDs)
	if summary.ExitStatus() != 0 {
		t.Errorf("exit status = %d, want 0", summary.ExitStatus())
	}
}

func TestFailedRidesCountedAndNotRecorded(t *testing.T) {
	config := testConfig(t, 3)

	summary := runQuietly(config, demoRides())

	if summary.Failed != 2 {
		t.Errorf("failed rides = %d, want 2", summary.Failed)
	}
	for _, id := range fileRideIDs(t, config.ResultsPath) {
		if id == "RIDE-8" || id == "RIDE-9" {
			t.Errorf("results file has failed ride %s", id)
		}
	}
}

func TestTotalsSameForAnyDriverCount(t *testing.T) {
	for _, driverCount := range []int{1, 8} {
		t.Run(fmt.Sprintf("%d drivers", driverCount), func(t *testing.T) {
			summary := runQuietly(testConfig(t, driverCount), demoRides())

			if len(summary.Results) != 7 || summary.Failed != 2 || summary.TotalCents() != demoTotalCents {
				t.Errorf("completed %d, failed %d, total %d cents; want 7, 2, %d",
					len(summary.Results), summary.Failed, summary.TotalCents(), demoTotalCents)
			}
		})
	}
}

func TestStressThousandRidesEightDrivers(t *testing.T) {
	const rideCount = 1000
	rides := make([]Ride, rideCount)
	ids := make([]string, rideCount)
	for i := range rides {
		ids[i] = fmt.Sprintf("RIDE-%d", i+1)
		rides[i] = Ride{ID: ids[i], Type: Standard, Pickup: "A", Dropoff: "B", Miles: 1.5, Rider: "R", Passengers: 1}
	}
	config := testConfig(t, 8)

	summary := runQuietly(config, rides)

	assertEachOnce(t, "results list", resultRideIDs(summary), ids)
	assertEachOnce(t, "results file", fileRideIDs(t, config.ResultsPath), ids)
}

func TestResultsFileError(t *testing.T) {
	config := testConfig(t, 3)
	config.ResultsPath = filepath.Join(t.TempDir(), "no-such-dir", "results.txt")

	summary := runQuietly(config, demoRides())

	assertEachOnce(t, "results list", resultRideIDs(summary), demoCompletedIDs)
	if summary.ExitStatus() != 1 {
		t.Errorf("exit status = %d, want 1", summary.ExitStatus())
	}
}

func TestTimeoutStopsEveryDriver(t *testing.T) {
	config := testConfig(t, 3)
	config.TripDelayPerMile = 100 * time.Millisecond
	config.Timeout = 50 * time.Millisecond
	// Without the timeout, RIDE-1 alone takes 1.2 seconds.
	const fullRunTime = 1200 * time.Millisecond

	start := time.Now()
	summary := runQuietly(config, demoRides())

	if elapsed := time.Since(start); elapsed >= fullRunTime {
		t.Errorf("run took %v, want the timeout to stop the drivers before %v", elapsed, fullRunTime)
	}
	if len(summary.Results) >= 7 {
		t.Errorf("completed %d rides, want the timeout to stop some", len(summary.Results))
	}
	if summary.ExitStatus() != 1 {
		t.Errorf("exit status = %d, want 1", summary.ExitStatus())
	}
}
