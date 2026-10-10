// Test cases 1 and 2 of the specification: fares and validation.

package main

import (
	"errors"
	"testing"
)

func demoRide(t *testing.T, id string) Ride {
	t.Helper()
	for _, ride := range demoRides() {
		if ride.ID == id {
			return ride
		}
	}
	t.Fatalf("no demo ride %s", id)
	return Ride{}
}

func TestFares(t *testing.T) {
	want := map[string]int64{
		"RIDE-1": 2000,
		"RIDE-2": 1850,
		"RIDE-3": 440,
		"RIDE-4": 680,
		"RIDE-5": 5000,
		"RIDE-6": 373,
		"RIDE-7": 500,
	}
	for id, cents := range want {
		if got := demoRide(t, id).FareCents(); got != cents {
			t.Errorf("%s fare = %d cents, want %d", id, got, cents)
		}
	}
}

func TestValidation(t *testing.T) {
	want := map[string]string{
		"RIDE-8": "distance must be greater than 0",
		"RIDE-9": "a shared ride needs 2 or more passengers",
	}
	for _, ride := range demoRides() {
		err := ride.Validate()
		text, isInvalid := want[ride.ID]
		if !isInvalid {
			if err != nil {
				t.Errorf("%s: unexpected error %v", ride.ID, err)
			}
			continue
		}
		if !errors.Is(err, ErrInvalidRide) {
			t.Errorf("%s: error %v does not wrap ErrInvalidRide", ride.ID, err)
		} else if err.Error() != text {
			t.Errorf("%s: error text %q, want %q", ride.ID, err.Error(), text)
		}
	}
}
