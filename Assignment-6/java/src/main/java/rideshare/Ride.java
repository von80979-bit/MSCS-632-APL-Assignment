package rideshare;

import java.util.Locale;

/** One ride request. Money is in whole cents, so the Java and Go programs match to the cent. */
public record Ride(String id, RideType rideType, String pickup, String dropoff, double miles, String rider,
        int passengers) {

    public void validate() throws InvalidRideException {
        if (miles <= 0) {
            throw new InvalidRideException("distance must be greater than 0");
        }
        if (rideType == RideType.SHARED && passengers < 2) {
            throw new InvalidRideException("a shared ride needs 2 or more passengers");
        }
    }

    public long fareCents() {
        return switch (rideType) {
            case STANDARD -> standardFareCents();
            case PREMIUM -> 500 + Math.round(300 * miles);
            // round(standard fare × 0.8 / passengers), half up, in integer arithmetic with no floating point error
            case SHARED -> (standardFareCents() * 8 + 5L * passengers) / (10L * passengers);
        };
    }

    private long standardFareCents() {
        return 200 + Math.round(150 * miles);
    }

    /** Formats cents as dollars with two decimals, without the dollar sign, for example {@code 18.50}. */
    public static String formatCents(long cents) {
        return String.format(Locale.ROOT, "%d.%02d", cents / 100, cents % 100);
    }
}
