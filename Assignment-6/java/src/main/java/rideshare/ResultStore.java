package rideshare;

import java.io.BufferedWriter;
import java.io.FileWriter;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.locks.ReentrantLock;
import java.util.logging.Logger;

/**
 * The shared results: the in-memory results list and the open results file. One lock guards both, so the list and
 * the file always agree. {@link #record} is the critical section that all drivers compete for.
 */
public class ResultStore implements AutoCloseable {

    private static final Logger LOG = Logger.getLogger(ResultStore.class.getName());

    private final List<RideResult> results = new ArrayList<>();
    private final BufferedWriter writer;
    private final ReentrantLock lock = new ReentrantLock();
    private boolean writeFailed;

    private ResultStore(BufferedWriter writer) {
        this.writer = writer;
    }

    /** Opens the results file, truncates it, and writes the header. */
    public static ResultStore open(Path path) throws IOException {
        BufferedWriter writer = new BufferedWriter(new FileWriter(path.toFile(), StandardCharsets.UTF_8));
        try {
            writer.write(RideResult.CSV_HEADER);
            writer.newLine();
            writer.flush();
        } catch (IOException e) {
            writer.close();
            throw e;
        }
        return new ResultStore(writer);
    }

    /** A store with no results file, for a run after the file could not open. */
    public static ResultStore inMemory() {
        return new ResultStore(null);
    }

    /** Adds the result to the list, and writes and flushes one line to the results file. */
    public void record(RideResult result) {
        lock.lock();
        try {
            results.add(result);
            if (writer != null) {
                try {
                    writer.write(result.csvLine());
                    writer.newLine();
                    writer.flush();
                } catch (IOException e) {
                    writeFailed = true;
                    LOG.severe("cannot write " + result.ride().id() + " to results file: " + e.getMessage());
                }
            }
        } finally {
            lock.unlock();
        }
    }

    public List<RideResult> results() {
        lock.lock();
        try {
            return List.copyOf(results);
        } finally {
            lock.unlock();
        }
    }

    /** True when the results file holds every result in the list. */
    public boolean fileComplete() {
        lock.lock();
        try {
            return writer != null && !writeFailed;
        } finally {
            lock.unlock();
        }
    }

    @Override
    public void close() throws IOException {
        lock.lock();
        try {
            if (writer != null) {
                writer.close();
            }
        } finally {
            lock.unlock();
        }
    }
}
