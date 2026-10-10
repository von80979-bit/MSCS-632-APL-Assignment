// Starts the drivers, adds the rides, closes the queue, waits for the drivers, and builds the summary.

package main

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"
)

// Config holds the run parameters. The tests set the trip delay to 0 and the results path to a temporary directory.
type Config struct {
	DriverNames      []string
	ResultsPath      string
	TripDelayPerMile time.Duration
	Timeout          time.Duration // time limit for the drivers after the queue closes
}

type Summary struct {
	DriverNames  []string
	Results      []RideResult
	Failed       int
	ResultsPath  string
	FileLines    int
	FileComplete bool
	TimedOut     bool
}

func RunSystem(config Config, rides []Ride, logger *Logger) (summary Summary) {
	logger.Info("main", "starting %d drivers, results file %s", len(config.DriverNames), config.ResultsPath)
	store := &ResultStore{}
	if err := store.OpenFile(config.ResultsPath); err != nil {
		logger.Error("main", "%v", err)
	}
	defer func() {
		if err := store.Close(); err != nil {
			logger.Error("main", "%v", err)
		}
		summary.FileLines = store.FileLines()
		summary.FileComplete = store.FileComplete()
	}()

	queue := NewRideQueue()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	failedCounts := make([]int, len(config.DriverNames))
	for i, name := range config.DriverNames {
		driver := Driver{Name: name, Queue: queue, Store: store, TripDelayPerMile: config.TripDelayPerMile, Log: logger}
		wg.Add(1)
		go func() {
			defer wg.Done()
			failedCounts[i] = driver.Run(ctx) // each goroutine writes only its own element, so no lock is needed
		}()
	}

	for _, ride := range rides {
		queue.AddTask(ride)
	}
	queue.Close()
	logger.Info("main", "queue closed")

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// The timeout starts when the queue closes. When it passes, cancel stops every driver that is on a trip.
	timer := time.NewTimer(config.Timeout)
	defer timer.Stop()
	timedOut := false
	select {
	case <-done:
	case <-timer.C:
		// The drivers can finish at the same moment that the timer fires. Check done again, so a finished run is not a timeout.
		select {
		case <-done:
		default:
			timedOut = true
			logger.Warn("main", "timeout after %gs, stopping drivers", config.Timeout.Seconds())
			cancel()
			<-done
		}
	}

	summary = Summary{DriverNames: config.DriverNames, Results: store.Results(), ResultsPath: config.ResultsPath, TimedOut: timedOut}
	for _, failed := range failedCounts {
		summary.Failed += failed
	}
	if summary.TimedOut {
		logger.Warn("main", "drivers stopped by the timeout")
	} else {
		logger.Info("main", "all drivers finished")
	}
	return summary
}

func (s Summary) TotalCents() int64 {
	var total int64
	for _, result := range s.Results {
		total += result.FareCents
	}
	return total
}

// ExitStatus is 1 after a results file error or a timeout. Invalid rides do not change it.
func (s Summary) ExitStatus() int {
	if !s.FileComplete || s.TimedOut {
		return 1
	}
	return 0
}

func (s Summary) Print(w io.Writer) {
	fmt.Fprintln(w, "--- Summary ---")
	fmt.Fprintf(w, "Rides completed: %d | Rides failed: %d | Total fares: %s\n",
		len(s.Results), s.Failed, formatDollars(s.TotalCents()))
	for _, name := range s.DriverNames {
		rides, cents := 0, int64(0)
		for _, result := range s.Results {
			if result.Driver == name {
				rides++
				cents += result.FareCents
			}
		}
		fmt.Fprintf(w, "%-11s %d rides %7s\n", name, rides, formatDollars(cents))
	}
	if s.FileComplete {
		fmt.Fprintf(w, "Results file: %s (%d rides)\n", s.ResultsPath, s.FileLines)
	} else {
		fmt.Fprintln(w, "Results file: not written (see the ERROR line)")
	}
}
