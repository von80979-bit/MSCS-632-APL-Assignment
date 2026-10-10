package rideshare;

import java.util.Optional;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.logging.Logger;

/**
 * A driver takes rides from the ride queue until the queue is empty and closed. It never holds the queue lock and the
 * results lock at the same time, so no circular wait can occur.
 */
public class Driver implements Runnable {

    private static final Logger LOG = Logger.getLogger(Driver.class.getName());

    private final String name;
    private final RideQueue queue;
    private final ResultStore store;
    private final AtomicInteger failedRides;
    private final long tripMillisPerMile;

    public Driver(String name, RideQueue queue, ResultStore store, AtomicInteger failedRides, long tripMillisPerMile) {
        this.name = name;
        this.queue = queue;
        this.store = store;
        this.failedRides = failedRides;
        this.tripMillisPerMile = tripMillisPerMile;
    }

    @Override
    public void run() {
        // The log format shows the thread name, so the pool thread takes the driver name.
        Thread.currentThread().setName(name);
        LOG.info("started");
        try {
            while (true) {
                Optional<Ride> next = queue.getTask();
                if (next.isEmpty()) {
                    LOG.info("queue empty and closed, exiting");
                    return;
                }
                drive(next.get());
            }
        } catch (InterruptedException e) {
            LOG.warning("interrupted, exiting");
            Thread.currentThread().interrupt();
        }
    }

    private void drive(Ride ride) throws InterruptedException {
        LOG.info("took " + ride.id());
        try {
            ride.validate();
        } catch (InvalidRideException e) {
            LOG.warning("failed " + ride.id() + ": " + e.getMessage());
            failedRides.incrementAndGet();
            return;
        }
        Thread.sleep(Math.round(ride.miles() * tripMillisPerMile));
        long fareCents = ride.fareCents();
        store.record(new RideResult(ride, name, fareCents));
        LOG.info("completed " + ride.id() + ", fare $" + Ride.formatCents(fareCents));
    }
}
