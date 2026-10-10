package rideshare;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

import java.util.Map;
import java.util.Set;
import org.junit.jupiter.api.Test;

class FareTest {

    // Test case 1: fares
    @Test
    void eachDemoRideGetsTheExpectedFare() {
        Map<String, Long> expectedCents = Map.of(
                "RIDE-1", 2000L, "RIDE-2", 1850L, "RIDE-3", 440L, "RIDE-4", 680L,
                "RIDE-5", 5000L, "RIDE-6", 373L, "RIDE-7", 500L);
        for (Ride ride : Main.demoRides()) {
            if (expectedCents.containsKey(ride.id())) {
                assertEquals(expectedCents.get(ride.id()), ride.fareCents(), ride.id());
            }
        }
    }

    // Test case 2: validation
    @Test
    void invalidDemoRidesFailWithTheErrorText() {
        Map<String, String> expectedErrors = Map.of(
                "RIDE-8", "distance must be greater than 0",
                "RIDE-9", "a shared ride needs 2 or more passengers");
        for (Ride ride : Main.demoRides()) {
            if (expectedErrors.containsKey(ride.id())) {
                InvalidRideException error = assertThrows(InvalidRideException.class, ride::validate, ride.id());
                assertEquals(expectedErrors.get(ride.id()), error.getMessage());
            } else {
                assertDoesNotThrow(ride::validate, ride.id());
            }
        }
        assertEquals(Set.of("RIDE-1", "RIDE-2", "RIDE-3", "RIDE-4", "RIDE-5", "RIDE-6", "RIDE-7", "RIDE-8", "RIDE-9"),
                Main.demoRides().stream().map(Ride::id).collect(java.util.stream.Collectors.toSet()));
    }
}
