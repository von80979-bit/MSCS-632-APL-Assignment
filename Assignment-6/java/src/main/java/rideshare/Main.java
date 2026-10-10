package rideshare;

import java.nio.file.Path;
import java.time.Duration;
import java.util.List;

/** Reads the flags, runs the demo, prints the summary, and exits with the status of the run. */
public final class Main {

    static final List<String> DRIVER_NAMES = List.of(
            "Ana Lopez", "Ben Carter", "Chen Wei", "Dana Kim", "Eli Novak", "Fatima Ali", "Gus Moreno", "Hana Sato");

    private static final long TRIP_MILLIS_PER_MILE = 100;
    private static final String USAGE = "usage: rideshare [-drivers 1-8] [-out path] [-timeout seconds]";

    private Main() {
    }

    public static void main(String[] args) throws InterruptedException {
        int drivers = 3;
        String resultsPath = "results.txt";
        int timeoutSeconds = 30;
        try {
            // Each flag is "-name value" or "-name=value", with one or two leading dashes, as in the Go flag package.
            for (int i = 0; i < args.length; i++) {
                if (!args[i].startsWith("-")) {
                    throw new IllegalArgumentException("not a flag: " + args[i]);
                }
                String name = args[i].replaceFirst("^--?", "");
                String value;
                int equals = name.indexOf('=');
                if (equals >= 0) {
                    value = name.substring(equals + 1);
                    name = name.substring(0, equals);
                } else if (i + 1 < args.length) {
                    value = args[++i];
                } else {
                    throw new IllegalArgumentException("flag needs a value: " + name);
                }
                switch (name) {
                    case "drivers" -> drivers = parseInRange(value, 1, DRIVER_NAMES.size());
                    case "out" -> resultsPath = value;
                    case "timeout" -> timeoutSeconds = parseInRange(value, 1, Integer.MAX_VALUE);
                    default -> throw new IllegalArgumentException("unknown flag: " + name);
                }
            }
        } catch (IllegalArgumentException e) {
            System.err.println(USAGE);
            System.exit(2);
        }

        LogFormat.install();
        System.out.println("=== Concurrent Ride Sharing System ===");
        var system = new RideSharingSystem(DRIVER_NAMES.subList(0, drivers), Path.of(resultsPath), TRIP_MILLIS_PER_MILE,
                Duration.ofSeconds(timeoutSeconds));
        RideSharingSystem.Summary summary = system.run(demoRides());
        summary.lines().forEach(System.out::println);
        System.exit(summary.exitStatus());
    }

    /** Throws IllegalArgumentException for a value that is not a number or is out of range. */
    private static int parseInRange(String value, int min, int max) {
        int number = Integer.parseInt(value);
        if (number < min || number > max) {
            throw new IllegalArgumentException("value out of range: " + value);
        }
        return number;
    }

    /** The demo rides in queue order. RIDE-8 and RIDE-9 are invalid and sit between valid rides. */
    static List<Ride> demoRides() {
        return List.of(
                new Ride("RIDE-1", RideType.STANDARD, "Downtown", "Airport", 12.0, "Maria Garcia", 1),
                new Ride("RIDE-2", RideType.PREMIUM, "Hotel", "Convention Center", 4.5, "James Smith", 1),
                new Ride("RIDE-8", RideType.STANDARD, "Office", "Office", 0.0, "James Smith", 1),
                new Ride("RIDE-3", RideType.SHARED, "University", "Stadium", 6.0, "Maria Garcia", 2),
                new Ride("RIDE-4", RideType.STANDARD, "Mall", "Train Station", 3.2, "James Smith", 1),
                new Ride("RIDE-9", RideType.SHARED, "Pier", "Market", 5.0, "Maria Garcia", 1),
                new Ride("RIDE-5", RideType.PREMIUM, "Airport", "Harbor", 15.0, "Maria Garcia", 1),
                new Ride("RIDE-6", RideType.SHARED, "Library", "Park", 8.0, "James Smith", 3),
                new Ride("RIDE-7", RideType.STANDARD, "Museum", "Zoo", 2.0, "Maria Garcia", 1));
    }
}
