package rideshare;

/** A ride that a driver cannot drive. The message is the error text of the specification. */
public class InvalidRideException extends Exception {

    public InvalidRideException(String message) {
        super(message);
    }
}
