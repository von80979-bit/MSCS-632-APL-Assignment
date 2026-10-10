package rideshare;

import java.util.ArrayDeque;
import java.util.Optional;
import java.util.concurrent.locks.Condition;
import java.util.concurrent.locks.ReentrantLock;

/**
 * The bounded ride queue. The main thread adds rides and closes the queue, and the drivers take rides from it.
 * One lock guards the deque and the closed flag. Every method releases the lock in a finally block.
 */
public class RideQueue {

    private final ArrayDeque<Ride> rides = new ArrayDeque<>();
    private final int capacity;
    private final ReentrantLock lock = new ReentrantLock();
    private final Condition notEmpty = lock.newCondition();
    private final Condition notFull = lock.newCondition();
    private boolean closed;

    public RideQueue(int capacity) {
        this.capacity = capacity;
    }

    /** Adds a ride, and waits while the queue is full. */
    public void addTask(Ride ride) throws InterruptedException {
        lock.lockInterruptibly();
        try {
            if (closed) {
                throw new IllegalStateException("the ride queue is closed");
            }
            while (rides.size() == capacity) {
                notFull.await();
            }
            rides.addLast(ride);
            notEmpty.signal();
        } finally {
            lock.unlock();
        }
    }

    /**
     * Returns the next ride, and waits while the queue is empty and open. Returns no ride when the queue is empty and
     * closed, which tells the driver to exit.
     */
    public Optional<Ride> getTask() throws InterruptedException {
        // lockInterruptibly throws at once when the driver is already interrupted, so it takes no ride after a timeout.
        lock.lockInterruptibly();
        try {
            while (rides.isEmpty() && !closed) {
                notEmpty.await();
            }
            Ride ride = rides.pollFirst();
            if (ride != null) {
                notFull.signal();
            }
            return Optional.ofNullable(ride);
        } finally {
            lock.unlock();
        }
    }

    /** Closes the queue and wakes every waiting driver, so each one can see the empty and closed queue. */
    public void close() {
        lock.lock();
        try {
            closed = true;
            notEmpty.signalAll();
        } finally {
            lock.unlock();
        }
    }
}
