package rideshare;

import java.io.IOException;
import java.nio.file.Path;
import java.time.Duration;
import java.util.ArrayList;
import java.util.List;
import java.util.Locale;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.logging.Logger;

/** Starts the drivers, adds the rides, closes the queue, waits for the drivers, and builds the summary. */
public class RideSharingSystem {

    static final int QUEUE_CAPACITY = 4;

    private static final Logger LOG = Logger.getLogger(RideSharingSystem.class.getName());

    private final List<String> driverNames;
    private final Path resultsPath;
    private final long tripMillisPerMile;
    private final Duration timeout;

    public RideSharingSystem(List<String> driverNames, Path resultsPath, long tripMillisPerMile, Duration timeout) {
        this.driverNames = List.copyOf(driverNames);
        this.resultsPath = resultsPath;
        this.tripMillisPerMile = tripMillisPerMile;
        this.timeout = timeout;
    }

    /** The outcome of one run. */
    public record Summary(List<String> driverNames, List<RideResult> results, int failedRides, boolean fileWritten,
            boolean timedOut, Path resultsPath) {

        public long totalFareCents() {
            return sumFares(results);
        }

        /** 0 when every driver finishes and the results file is complete, otherwise 1. */
        public int exitStatus() {
            return fileWritten && !timedOut ? 0 : 1;
        }

        public List<String> lines() {
            List<String> lines = new ArrayList<>();
            lines.add("--- Summary ---");
            lines.add(String.format(Locale.ROOT, "Rides completed: %d | Rides failed: %d | Total fares: $%s",
                    results.size(), failedRides, Ride.formatCents(totalFareCents())));
            for (String driver : driverNames) {
                List<RideResult> driven = results.stream().filter(result -> result.driver().equals(driver)).toList();
                lines.add(String.format(Locale.ROOT, "%-11s %d rides %7s", driver, driven.size(),
                        "$" + Ride.formatCents(sumFares(driven))));
            }
            lines.add(fileWritten
                    ? "Results file: " + resultsPath + " (" + results.size() + " rides)"
                    : "Results file: not written (see the ERROR line)");
            return lines;
        }

        private static long sumFares(List<RideResult> results) {
            return results.stream().mapToLong(RideResult::fareCents).sum();
        }
    }

    public Summary run(List<Ride> rides) throws InterruptedException {
        LOG.info("starting " + driverNames.size() + " drivers, results file " + resultsPath);
        boolean fileOpened = true;
        ResultStore store;
        try {
            store = ResultStore.open(resultsPath);
        } catch (IOException e) {
            LOG.severe("cannot open results file: " + e.getMessage());
            fileOpened = false;
            store = ResultStore.inMemory();
        }

        AtomicInteger failedRides = new AtomicInteger();
        boolean timedOut = false;
        boolean fileClosed = true;
        try (ResultStore openStore = store) {
            timedOut = runDrivers(rides, openStore, failedRides);
        } catch (IOException e) {
            // The close runs after runDrivers returns, so timedOut keeps its value.
            LOG.severe("cannot close results file: " + e.getMessage());
            fileClosed = false;
        }
        boolean fileWritten = fileOpened && fileClosed && store.fileComplete();
        return new Summary(driverNames, store.results(), failedRides.get(), fileWritten, timedOut, resultsPath);
    }

    /** Returns true when the timeout passed and the drivers were interrupted. */
    private boolean runDrivers(List<Ride> rides, ResultStore store, AtomicInteger failedRides)
            throws InterruptedException {
        RideQueue queue = new RideQueue(QUEUE_CAPACITY);
        ExecutorService pool = Executors.newFixedThreadPool(driverNames.size());
        for (String name : driverNames) {
            pool.execute(new Driver(name, queue, store, failedRides, tripMillisPerMile));
        }
        try {
            for (Ride ride : rides) {
                queue.addTask(ride);
            }
        } catch (InterruptedException e) {
            // Stop the drivers before the caller closes the results file, so no driver writes to a closed file.
            pool.shutdownNow();
            pool.awaitTermination(timeout.toMillis(), TimeUnit.MILLISECONDS);
            throw e;
        } finally {
            queue.close();
            pool.shutdown();
        }
        LOG.info("queue closed");

        if (pool.awaitTermination(timeout.toMillis(), TimeUnit.MILLISECONDS)) {
            LOG.info("all drivers finished");
            return false;
        }
        LOG.warning("timeout after " + timeout.toSeconds() + "s, stopping drivers");
        pool.shutdownNow();
        pool.awaitTermination(timeout.toMillis(), TimeUnit.MILLISECONDS);
        LOG.warning("drivers stopped by the timeout");
        return true;
    }
}
