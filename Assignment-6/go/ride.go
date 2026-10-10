// The ride value, the ride types, the fare rules, and validation.

package main

import (
	"errors"
	"fmt"
	"math"
)

type RideType int

const (
	Standard RideType = iota
	Premium
	Shared
)

func (t RideType) String() string {
	switch t {
	case Standard:
		return "Standard"
	case Premium:
		return "Premium"
	case Shared:
		return "Shared"
	}
	return fmt.Sprintf("RideType(%d)", int(t))
}

type Ride struct {
	ID         string
	Type       RideType
	Pickup     string
	Dropoff    string
	Miles      float64
	Rider      string
	Passengers int
}

// ErrInvalidRide is wrapped by every validation error, so a caller can test for it with errors.Is.
var ErrInvalidRide = errors.New("invalid ride")

type invalidRideError struct {
	reason string
}

func (e invalidRideError) Error() string { return e.reason }
func (e invalidRideError) Unwrap() error { return ErrInvalidRide }

func (r Ride) Validate() error {
	if r.Miles <= 0 {
		return invalidRideError{"distance must be greater than 0"}
	}
	if r.Type == Shared && r.Passengers < 2 {
		return invalidRideError{"a shared ride needs 2 or more passengers"}
	}
	return nil
}

// FareCents calculates in whole cents. Miles have at most one decimal place, so the
// calculation runs on tenths of a mile in integers and rounds half up with no float error.
func (r Ride) FareCents() int64 {
	tenths := int64(math.Round(r.Miles * 10))
	standard := 200 + (150*tenths+5)/10
	switch r.Type {
	case Premium:
		return 500 + (300*tenths+5)/10
	case Shared:
		// round(standard × 0.8 / passengers) = round(8 × standard / (10 × passengers))
		passengers := int64(r.Passengers)
		return (8*standard + 5*passengers) / (10 * passengers)
	}
	return standard
}

// formatCents gives 1850 as 18.50, the fare format of the results file.
func formatCents(cents int64) string {
	return fmt.Sprintf("%d.%02d", cents/100, cents%100)
}

func formatDollars(cents int64) string {
	return "$" + formatCents(cents)
}
