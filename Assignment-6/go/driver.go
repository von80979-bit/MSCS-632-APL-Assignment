// The driver loop. Each driver runs in its own goroutine.

package main

import (
	"context"
	"errors"
	"time"
)

type Driver struct {
	Name             string
	Queue            *RideQueue
	Store            *ResultStore
	TripDelayPerMile time.Duration
	Log              *Logger
}

// Run takes rides until the queue is empty and closed, or until ctx is cancelled.
// It returns the number of failed rides.
// The driver holds no lock during the trip, and never holds the queue and the store at the same time.
func (d Driver) Run(ctx context.Context) int {
	failed := 0
	d.Log.Info(d.Name, "started")
	for {
		if ctx.Err() != nil {
			d.Log.Warn(d.Name, "interrupted, exiting")
			return failed
		}
		ride, found := d.Queue.GetTask()
		if !found {
			d.Log.Info(d.Name, "queue empty and closed, exiting")
			return failed
		}
		d.Log.Info(d.Name, "took %s", ride.ID)

		if err := ride.Validate(); errors.Is(err, ErrInvalidRide) {
			d.Log.Warn(d.Name, "failed %s: %v", ride.ID, err)
			failed++
			continue
		}

		trip := time.Duration(ride.Miles * float64(d.TripDelayPerMile))
		select {
		case <-time.After(trip):
		case <-ctx.Done():
			d.Log.Warn(d.Name, "interrupted, exiting")
			return failed
		}

		fare := ride.FareCents()
		if err := d.Store.Record(RideResult{Ride: ride, Driver: d.Name, FareCents: fare}); err != nil {
			d.Log.Error(d.Name, "%v", err)
		}
		d.Log.Info(d.Name, "completed %s, fare %s", ride.ID, formatDollars(fare))
	}
}
