package rideshare;

import java.time.LocalTime;
import java.time.ZoneId;
import java.time.format.DateTimeFormatter;
import java.util.logging.Formatter;
import java.util.logging.Level;
import java.util.logging.LogRecord;
import java.util.logging.Logger;
import java.util.logging.StreamHandler;

/** Formats each log line as {@code HH:MM:SS.mmm LEVEL WHO        | message}. WHO is the thread name. */
public class LogFormat extends Formatter {

    private static final DateTimeFormatter TIME = DateTimeFormatter.ofPattern("HH:mm:ss.SSS");

    /** Sends every log line to standard output in this format, in place of the default handler on standard error. */
    public static void install() {
        Logger root = Logger.getLogger("");
        for (var handler : root.getHandlers()) {
            root.removeHandler(handler);
        }
        root.addHandler(new StreamHandler(System.out, new LogFormat()) {
            @Override
            public synchronized void publish(LogRecord record) {
                super.publish(record);
                flush();
            }
        });
    }

    @Override
    public String format(LogRecord record) {
        String time = LocalTime.ofInstant(record.getInstant(), ZoneId.systemDefault()).format(TIME);
        // The handler formats on the thread that logs, so the current thread is the driver or main.
        return String.format("%s %-5s %-10s | %s%n", time, level(record.getLevel()), Thread.currentThread().getName(),
                formatMessage(record));
    }

    private static String level(Level level) {
        if (level.intValue() >= Level.SEVERE.intValue()) {
            return "ERROR";
        }
        return level.intValue() >= Level.WARNING.intValue() ? "WARN" : "INFO";
    }
}
