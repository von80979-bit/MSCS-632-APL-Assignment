package rideshare;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTimeoutPreemptively;
import static org.junit.jupiter.api.Assertions.assertTrue;

import java.time.Duration;
import java.util.Optional;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;
import org.junit.jupiter.api.Test;

class RideQueueTest {

    private static final Duration LIMIT = Duration.ofSeconds(2);

    private static Ride ride(String id) {
        return new Ride(id, RideType.STANDARD, "A", "B", 1.0, "Rider", 1);
    }

    // Waits until the thread blocks on a lock condition, so the test knows the call is waiting.
    private static void awaitWaiting(Thread thread) throws InterruptedException {
        long deadline = System.nanoTime() + LIMIT.toNanos();
        while (thread.getState() != Thread.State.WAITING) {
            assertTrue(System.nanoTime() < deadline, "the thread did not start to wait");
            Thread.sleep(5);
        }
    }

    // Test case 7: empty and closed queue
    @Test
    void getTaskOnEmptyClosedQueueReturnsNoRideAtOnce() {
        RideQueue queue = new RideQueue(4);
        queue.close();
        assertTimeoutPreemptively(LIMIT, () -> assertEquals(Optional.empty(), queue.getTask()));
    }

    // Test case 8: blocked getTask wakes up on addTask
    @Test
    void blockedGetTaskWakesUpOnAddTask() throws Exception {
        RideQueue queue = new RideQueue(4);
        CompletableFuture<Optional<Ride>> taken = new CompletableFuture<>();
        Thread driver = new Thread(() -> {
            try {
                taken.complete(queue.getTask());
            } catch (InterruptedException e) {
                taken.completeExceptionally(e);
            }
        });
        driver.start();
        awaitWaiting(driver);
        assertFalse(taken.isDone());

        queue.addTask(ride("RIDE-1"));
        assertEquals("RIDE-1", taken.get(2, TimeUnit.SECONDS).orElseThrow().id());
    }

    // Test case 8: blocked getTask wakes up on close
    @Test
    void blockedGetTaskWakesUpOnClose() throws Exception {
        RideQueue queue = new RideQueue(4);
        CompletableFuture<Optional<Ride>> taken = new CompletableFuture<>();
        Thread driver = new Thread(() -> {
            try {
                taken.complete(queue.getTask());
            } catch (InterruptedException e) {
                taken.completeExceptionally(e);
            }
        });
        driver.start();
        awaitWaiting(driver);
        assertFalse(taken.isDone());

        queue.close();
        assertEquals(Optional.empty(), taken.get(2, TimeUnit.SECONDS));
    }

    // Test case 9: full queue
    @Test
    void addTaskOnFullQueueWaitsUntilARideIsTaken() throws Exception {
        RideQueue queue = new RideQueue(4);
        for (int i = 1; i <= 4; i++) {
            queue.addTask(ride("RIDE-" + i));
        }
        CompletableFuture<Void> added = new CompletableFuture<>();
        Thread producer = new Thread(() -> {
            try {
                queue.addTask(ride("RIDE-5"));
                added.complete(null);
            } catch (InterruptedException e) {
                added.completeExceptionally(e);
            }
        });
        producer.start();
        awaitWaiting(producer);
        assertFalse(added.isDone());

        assertEquals("RIDE-1", queue.getTask().orElseThrow().id());
        added.get(2, TimeUnit.SECONDS);
        queue.close();
        for (int i = 2; i <= 5; i++) {
            assertEquals("RIDE-" + i, queue.getTask().orElseThrow().id());
        }
        assertEquals(Optional.empty(), queue.getTask());
    }
}
