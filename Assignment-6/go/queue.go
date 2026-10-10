// The ride queue: a buffered channel, so the channel does all the synchronization and no lock is needed.

package main

const queueCapacity = 4

// RideQueue allows one producer only. Go lets only the sender close a channel,
// so the main goroutine is the only caller of AddTask and Close.
type RideQueue struct {
	rides chan Ride
}

func NewRideQueue() *RideQueue {
	return &RideQueue{rides: make(chan Ride, queueCapacity)}
}

// AddTask waits while the queue is full.
func (q *RideQueue) AddTask(ride Ride) {
	q.rides <- ride
}

// GetTask waits while the queue is empty and open. It returns false when the queue is empty and closed.
func (q *RideQueue) GetTask() (Ride, bool) {
	ride, found := <-q.rides
	return ride, found
}

// Close wakes up every driver that waits in GetTask.
func (q *RideQueue) Close() {
	close(q.rides)
}
