// The results store: the in-memory results list and the results file behind one mutex.

package main

import (
	"errors"
	"fmt"
	"os"
	"sync"
)

const resultsHeader = "ride_id,driver,rider,ride_type,miles,fare\n"

// RideResult is one completed ride.
type RideResult struct {
	Ride      Ride
	Driver    string
	FareCents int64
}

// ResultStore keeps the list and the file in agreement, because one mutex guards both.
// Record is the critical section that all drivers compete for.
type ResultStore struct {
	mu           sync.Mutex
	results      []RideResult
	file         *os.File // nil when the results file cannot open
	fileLines    int
	fileComplete bool
}

// OpenFile truncates the results file and writes the header. After an error, the store keeps the list only.
func (s *ResultStore) OpenFile(path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot open results file: %w", err)
	}
	if _, err := file.WriteString(resultsHeader); err != nil {
		return errors.Join(fmt.Errorf("cannot open results file: %w", err), file.Close())
	}
	s.file = file
	s.fileComplete = true
	return nil
}

// Record adds the result to the list and writes one line to the file.
// After a write error, the result stays in the list.
func (s *ResultStore) Record(result RideResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.results = append(s.results, result)
	if s.file == nil {
		return nil
	}
	ride := result.Ride
	_, err := fmt.Fprintf(s.file, "%s,%s,%s,%s,%.1f,%s\n",
		ride.ID, result.Driver, ride.Rider, ride.Type, ride.Miles, formatCents(result.FareCents))
	if err != nil {
		s.fileComplete = false
		return fmt.Errorf("cannot write %s to results file: %w", ride.ID, err)
	}
	s.fileLines++
	return nil
}

// Close closes the results file. A failed close means the file may be incomplete.
func (s *ResultStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.file == nil {
		return nil
	}
	err := s.file.Close()
	s.file = nil
	if err != nil {
		s.fileComplete = false
		return fmt.Errorf("cannot close results file: %w", err)
	}
	return nil
}

// Results returns a copy of the list, so the caller cannot change the store without the lock.
func (s *ResultStore) Results() []RideResult {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]RideResult(nil), s.results...)
}

func (s *ResultStore) FileLines() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.fileLines
}

func (s *ResultStore) FileComplete() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.fileComplete
}
