// Test cases 7, 8, and 9 of the specification: the ride queue.

package main

import (
	"testing"
	"time"
)

// blockWait is how long a test waits to decide that a call blocks.
const blockWait = 50 * time.Millisecond

// wakeWait is the longest time a test waits for a blocked call to return.
const wakeWait = time.Second

type taskResult struct {
	ride  Ride
	found bool
}

func getTaskAsync(queue *RideQueue) <-chan taskResult {
	results := make(chan taskResult, 1)
	go func() {
		ride, found := queue.GetTask()
		results <- taskResult{ride, found}
	}()
	return results
}

func TestGetTaskOnEmptyClosedQueue(t *testing.T) {
	queue := NewRideQueue()
	queue.Close()

	select {
	case result := <-getTaskAsync(queue):
		if result.found {
			t.Fatalf("GetTask returned %s, want no ride", result.ride.ID)
		}
	case <-time.After(wakeWait):
		t.Fatal("GetTask on an empty, closed queue did not return at once")
	}
}

func TestBlockedGetTaskWakesOnAddTask(t *testing.T) {
	queue := NewRideQueue()
	results := getTaskAsync(queue)

	select {
	case <-results:
		t.Fatal("GetTask on an empty, open queue did not wait")
	case <-time.After(blockWait):
	}
	queue.AddTask(Ride{ID: "RIDE-1"})

	select {
	case result := <-results:
		if !result.found || result.ride.ID != "RIDE-1" {
			t.Fatalf("GetTask returned %+v, want RIDE-1", result)
		}
	case <-time.After(wakeWait):
		t.Fatal("GetTask did not wake up on AddTask")
	}
}

func TestBlockedGetTaskWakesOnClose(t *testing.T) {
	queue := NewRideQueue()
	results := getTaskAsync(queue)

	select {
	case <-results:
		t.Fatal("GetTask on an empty, open queue did not wait")
	case <-time.After(blockWait):
	}
	queue.Close()

	select {
	case result := <-results:
		if result.found {
			t.Fatalf("GetTask returned %s, want no ride", result.ride.ID)
		}
	case <-time.After(wakeWait):
		t.Fatal("GetTask did not wake up on Close")
	}
}

func TestAddTaskWaitsOnFullQueue(t *testing.T) {
	queue := NewRideQueue()
	for i := range queueCapacity {
		queue.AddTask(Ride{ID: string(rune('A' + i))})
	}
	added := make(chan struct{})
	go func() {
		queue.AddTask(Ride{ID: "LAST"})
		close(added)
	}()

	select {
	case <-added:
		t.Fatal("AddTask on a full queue did not wait")
	case <-time.After(blockWait):
	}
	if _, found := queue.GetTask(); !found {
		t.Fatal("GetTask on a full queue returned no ride")
	}

	select {
	case <-added:
	case <-time.After(wakeWait):
		t.Fatal("AddTask did not continue after GetTask took a ride")
	}
}
