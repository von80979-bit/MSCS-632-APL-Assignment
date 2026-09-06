// The last reference to a large block is dropped and the collector takes the space back. The weak reference reads
// null afterwards, which is the collector reporting the block is gone rather than this program assuming it.
import java.lang.ref.WeakReference;

public class ErrorDemo {
    public static void main(String[] arguments) {
        byte[] largeBlock = new byte[64 * 1024 * 1024];
        WeakReference<byte[]> blockWatcher = new WeakReference<>(largeBlock);
        System.out.println("used heap while the block is held: " + usedHeapInMegabytes() + " MB");

        largeBlock = null;  // the only strong reference is gone, so the block is now unreachable
        System.gc();        // a request rather than an order, but G1 acts on it here

        System.out.println("used heap after collection:       " + usedHeapInMegabytes() + " MB");
        System.out.println("weak reference to the block:      " + blockWatcher.get());
    }

    private static long usedHeapInMegabytes() {
        Runtime runtime = Runtime.getRuntime();
        return (runtime.totalMemory() - runtime.freeMemory()) / (1024 * 1024);
    }
}
