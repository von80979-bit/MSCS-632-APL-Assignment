package rideshare;

/** The three ride types of Assignment 5. */
public enum RideType {
    STANDARD("Standard"),
    PREMIUM("Premium"),
    SHARED("Shared");

    private final String label;

    RideType(String label) {
        this.label = label;
    }

    @Override
    public String toString() {
        return label;
    }
}
