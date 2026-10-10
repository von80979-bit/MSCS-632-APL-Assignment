package rideshare;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTimeoutPreemptively;
import static org.junit.jupiter.api.Assertions.assertTrue;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.time.Duration;
import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Set;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;

class RideSharingSystemTest {

    private static final Duration TIMEOUT = Duration.ofSeconds(10);

    @TempDir
    Path tempDir;

    private RideSharingSystem.Summary runDemo(int drivers, Path resultsPath) throws InterruptedException {
        var system = new RideSharingSystem(Main.DRIVER_NAMES.subList(0, drivers), resultsPath, 0, TIMEOUT);
        return system.run(Main.demoRides());
    }

    private static List<String> rideIdsInFile(Path resultsPath) throws IOException {
        List<String> lines = Files.readAllLines(resultsPath);
        assertEquals("ride_id,driver,rider,ride_type,miles,fare", lines.get(0));
        return lines.subList(1, lines.size()).stream().map(line -> line.split(",")[0]).toList();
    }

    private static void assertEachRideOnce(List<String> rideIds, int expectedCount) {
        assertEquals(expectedCount, rideIds.size());
        assertEquals(expectedCount, new HashSet<>(rideIds).size(), "a ride ID appears twice");
    }

    // Test case 3: exactly once
    @Test
    void demoRecordsEachValidRideExactlyOnce() throws Exception {
        Path resultsPath = tempDir.resolve("results.txt");
        RideSharingSystem.Summary summary = runDemo(3, resultsPath);

        List<String> fileIds = rideIdsInFile(resultsPath);
        assertEachRideOnce(fileIds, 7);
        List<String> listIds = summary.results().stream().map(result -> result.ride().id()).toList();
        assertEachRideOnce(listIds, 7);
        assertEquals(Set.copyOf(fileIds), Set.copyOf(listIds));
        assertEquals(0, summary.exitStatus());
    }

    // Test case 4: failed rides
    @Test
    void demoCountsTwoFailedRidesAndWritesNoLineForThem() throws Exception {
        Path resultsPath = tempDir.resolve("results.txt");
        RideSharingSystem.Summary summary = runDemo(3, resultsPath);

        assertEquals(2, summary.failedRides());
        List<String> fileIds = rideIdsInFile(resultsPath);
        assertFalse(fileIds.contains("RIDE-8"));
        assertFalse(fileIds.contains("RIDE-9"));
    }

    // Test case 5: driver count
    @Test
    void oneDriverAndEightDriversGiveTheSameTotals() throws Exception {
        for (int drivers : new int[] {1, 8}) {
            RideSharingSystem.Summary summary = runDemo(drivers, tempDir.resolve("results-" + drivers + ".txt"));
            assertEquals(7, summary.results().size(), drivers + " drivers");
            assertEquals(2, summary.failedRides(), drivers + " drivers");
            assertEquals(10843, summary.totalFareCents(), drivers + " drivers");
            assertEquals(0, summary.exitStatus(), drivers + " drivers");
        }
    }

    // Test case 6: stress
    @Test
    void thousandRidesAndEightDriversRecordEachRideOnce() throws Exception {
        List<Ride> rides = new ArrayList<>();
        for (int i = 1; i <= 1000; i++) {
            rides.add(new Ride("RIDE-" + i, RideType.STANDARD, "A", "B", 1.0, "Rider", 1));
        }
        Path resultsPath = tempDir.resolve("results.txt");
        var system = new RideSharingSystem(Main.DRIVER_NAMES, resultsPath, 0, TIMEOUT);
        RideSharingSystem.Summary summary = system.run(rides);

        assertEachRideOnce(summary.results().stream().map(result -> result.ride().id()).toList(), 1000);
        assertEachRideOnce(rideIdsInFile(resultsPath), 1000);
        assertEquals(0, summary.exitStatus());
    }

    // Test case 10: file error
    @Test
    void resultsPathInMissingDirectoryGivesErrorAndRidesStillComplete() throws Exception {
        Path resultsPath = tempDir.resolve("no-such-dir").resolve("results.txt");
        RideSharingSystem.Summary summary = runDemo(3, resultsPath);

        assertFalse(Files.exists(resultsPath));
        assertFalse(summary.fileWritten());
        assertEquals(7, summary.results().size());
        assertEquals(10843, summary.totalFareCents());
        assertEquals(1, summary.exitStatus());
    }

    // Test case 11: timeout
    @Test
    void tripLongerThanTimeoutStopsEveryDriver() throws InterruptedException {
        List<Ride> rides = List.of(
                new Ride("RIDE-1", RideType.STANDARD, "A", "B", 1.0, "Rider", 1),
                new Ride("RIDE-2", RideType.STANDARD, "A", "B", 1.0, "Rider", 1));
        var system = new RideSharingSystem(Main.DRIVER_NAMES.subList(0, 3), tempDir.resolve("results.txt"),
                60_000, Duration.ofMillis(200));

        RideSharingSystem.Summary summary = assertTimeoutPreemptively(Duration.ofSeconds(5), () -> system.run(rides));
        assertTrue(summary.timedOut());
        assertEquals(0, summary.results().size());
        assertEquals(1, summary.exitStatus());
        assertNoDriverAlive(Main.DRIVER_NAMES.subList(0, 3));
    }

    // Each driver thread carries its driver name, so a live thread with that name is a driver that did not stop.
    private static void assertNoDriverAlive(List<String> driverNames) throws InterruptedException {
        long deadline = System.nanoTime() + Duration.ofSeconds(2).toNanos();
        while (Thread.getAllStackTraces().keySet().stream()
                .anyMatch(thread -> thread.isAlive() && driverNames.contains(thread.getName()))) {
            assertTrue(System.nanoTime() < deadline, "a driver is still running after the timeout");
            Thread.sleep(10);
        }
    }
}
