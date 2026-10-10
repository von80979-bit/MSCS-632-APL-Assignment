package rideshare;

import java.util.Locale;

/** One completed ride: the ride, the driver who drove it, and the fare. */
public record RideResult(Ride ride, String driver, long fareCents) {

    static final String CSV_HEADER = "ride_id,driver,rider,ride_type,miles,fare";

    public String csvLine() {
        return String.format(Locale.ROOT, "%s,%s,%s,%s,%.1f,%s", ride.id(), driver, ride.rider(), ride.rideType(),
                ride.miles(), Ride.formatCents(fareCents));
    }
}
