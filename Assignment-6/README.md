# Assignment 6: Concurrent ride sharing system

A concurrent ride sharing system in Java and in Go. The two programs follow one specification and print the same output, apart from the times and the driver of each ride. Drivers run in parallel, take rides from one shared ride queue, and record each completed ride in one shared results file. The program shows safe access to shared resources, safe termination, and error handling during concurrent work.

## What the application does

Each program runs one fixed demo and exits. It takes no input apart from its flags.

The main thread starts the drivers, adds nine rides to the ride queue, and closes the queue. Each driver takes the next ride, drives it, calculates the fare, and appends one line to the results file. The first free driver takes the next ride, so the driver of each ride changes from run to run. Two of the nine rides are invalid. A driver logs each invalid ride as a failed ride and continues with the next ride. When the queue is empty and closed, each driver exits. The main thread then waits for all drivers and prints a summary.

The program prints two parts:

- The log, with one line for each step of each driver and of the main thread.
- The summary, with the completed rides, the failed rides, the total fares, and the rides and fares of each driver.

## Fares

| Ride type | Fare |
|---|---|
| Standard | $2.00 booking fee + $1.50 per mile |
| Premium | $5.00 booking fee + $3.00 per mile |
| Shared | The standard fare divided by the number of passengers, with 20% off |

## Concurrency design

Both programs use the same design: one bounded ride queue, one results store under one lock, and no driver that holds two locks at the same time. Each program uses the tools of its own language.

### Java

- **Ride queue.** One `ReentrantLock` guards the queue, with the conditions `notEmpty` and `notFull`. `addTask` waits while the queue is full, and `getTask` waits while the queue is empty and open. `getTask` returns no ride when the queue is empty and closed, and the driver then exits.
- **Results store.** A second `ReentrantLock` guards the results list and the results file together, so the two always agree. Only one driver at a time can record a result.
- **Deadlock avoidance.** A driver never holds both locks at the same time. It releases the queue lock before the trip and takes the results lock only after the trip.
- **Safe termination.** The drivers run on a fixed thread pool. After the queue closes, the main thread waits for the pool to finish. If the timeout passes, the main thread interrupts the drivers.
- **Error handling.** An invalid ride throws `InvalidRideException`, and the driver catches it. An interrupt gives `InterruptedException`, and the driver logs it and exits. A results file error gives `IOException`, which the program logs before it continues with the in-memory results list. Try-with-resources closes the results file.

### Go

- **Ride queue.** A buffered channel of capacity 4 is the queue, so the channel does all the synchronization and the queue needs no lock. `AddTask` sends on the channel and waits while the channel is full. `GetTask` receives from the channel and returns `false` when the channel is empty and closed, and the driver then exits. Only the main goroutine sends on the channel and closes it.
- **Results store.** A `sync.Mutex` guards the results list and the results file together. `Record` releases the mutex with `defer`, so every return path unlocks it.
- **Deadlock avoidance.** A driver holds no lock during the trip. The channel operation ends before the trip, and the driver takes the mutex only after the trip.
- **Safe termination.** Each driver is a goroutine, and a `sync.WaitGroup` counts the drivers. Each driver calls `defer wg.Done()`, and the main goroutine waits with `wg.Wait()`. After the queue closes, a `time.AfterFunc` timer cancels a shared context when the timeout passes. A driver waits for the trip with `select` on a timer and `ctx.Done()`, so the cancel stops the trip.
- **Error handling.** `Validate` returns an `error` that wraps `ErrInvalidRide`, and the driver checks it. Each function that can fail returns an `error`, and the caller checks it. A results file error is logged, and the program continues with the in-memory results list. `defer` closes the results file, and the program checks the close error.

## Assumptions

The requirement leaves certain points open, so the program uses these assumptions:

1. Distance is in miles, and money is in US dollars. The program calculates money in whole cents.
2. The simulated trip lasts 100 ms per mile. The drivers then overlap in time, and the demo stays short.
3. A driver validates a ride when it processes the ride, not when the ride enters the queue.
4. The program logs and counts a failed ride, and the driver continues with the next ride. A failed ride writes no line to the results file.
5. The ride queue holds at most 4 rides. When it is full, `addTask` waits until a driver takes a ride.
6. The main thread is the only thread that adds rides. It closes the queue after the last ride.
7. Each run truncates the results file and writes a new header.

## Flags

| Flag | Default | Content |
|---|---|---|
| `-drivers <n>` | `3` | Number of drivers, from 1 to 8 |
| `-out <path>` | `results.txt` | Path of the results file |
| `-timeout <seconds>` | `30` | Time limit for the drivers after the queue closes |

A flag takes the form `-drivers 5` or `-drivers=5`, with one or two leading dashes, as in the Go `flag` package. An unknown flag, an argument that is not a flag, `-h`, or a value out of range prints a usage line and exits with status 2. The program exits with status 0 when every driver finishes and the results file is complete. It exits with status 1 after a results file error or a timeout.

## Run with Docker

Docker is the primary run path. Run these commands in this folder.

Build the image:

```sh
docker build -t concurrent-ride-sharing .
```

Run the Java program or the Go program:

```sh
docker run --rm concurrent-ride-sharing java
docker run --rm concurrent-ride-sharing go
```

Flags go after the language name. For example, this command runs 5 drivers:

```sh
docker run --rm concurrent-ride-sharing java -drivers 5
```

The results file stays inside the container. To keep it, mount a folder. Then give a path in that folder to `-out`:

```sh
docker run --rm -v "$PWD:/out" concurrent-ride-sharing java -out /out/results.txt
```

## Run without Docker

### Java

The Java program needs Java 21 or later and Maven 3.9 or later. Use these commands in the `java` folder to build the program and run it:

```sh
cd java
mvn package
java -jar target/rideshare.jar
```

To see the two error paths, run the program with a results file in a folder that does not exist, or with a timeout that is too short:

```sh
java -jar target/rideshare.jar -out /no/such/dir/results.txt
java -jar target/rideshare.jar -timeout 1
```

### Go

The Go program needs Go 1.24 or later and uses the standard library only. Use these commands in the `go` folder to build the program and run it:

```sh
cd go
go build -o rideshare-go .
./rideshare-go
```

The two error paths use the same flags:

```sh
./rideshare-go -out /no/such/dir/results.txt
./rideshare-go -timeout 1
```

## Run the tests

Run the Java tests in the `java` folder:

```sh
cd java
mvn test
```

Run the Go tests in the `go` folder. The `-race` flag turns on the Go race detector, which reports any data race during the tests:

```sh
cd go
go vet ./...
go test -race ./...
```

The two programs have the same test cases. The tests cover the fares, the validation, the ride queue, and full runs of the system. They set the trip delay to 0 and write the results file to a temporary folder.

## Sample output

The times and the driver of each ride change from run to run. The completed rides, the failed rides, and the total fares are the same in every run.

### Java

```
=== Concurrent Ride Sharing System ===
17:11:20.027 INFO  main       | starting 3 drivers, results file results.txt
17:11:20.051 INFO  Ana Lopez  | started
17:11:20.051 INFO  Chen Wei   | started
17:11:20.051 INFO  Ben Carter | started
17:11:20.052 INFO  Chen Wei   | took RIDE-2
17:11:20.052 INFO  Ana Lopez  | took RIDE-1
17:11:20.052 INFO  Ben Carter | took RIDE-8
17:11:20.055 WARN  Ben Carter | failed RIDE-8: distance must be greater than 0
17:11:20.055 INFO  Ben Carter | took RIDE-3
17:11:20.517 INFO  Chen Wei   | completed RIDE-2, fare $18.50
17:11:20.518 INFO  Chen Wei   | took RIDE-4
17:11:20.518 INFO  main       | queue closed
17:11:20.659 INFO  Ben Carter | completed RIDE-3, fare $4.40
17:11:20.659 INFO  Ben Carter | took RIDE-9
17:11:20.659 WARN  Ben Carter | failed RIDE-9: a shared ride needs 2 or more passengers
17:11:20.659 INFO  Ben Carter | took RIDE-5
17:11:20.843 INFO  Chen Wei   | completed RIDE-4, fare $6.80
17:11:20.843 INFO  Chen Wei   | took RIDE-6
17:11:21.257 INFO  Ana Lopez  | completed RIDE-1, fare $20.00
17:11:21.257 INFO  Ana Lopez  | took RIDE-7
17:11:21.462 INFO  Ana Lopez  | completed RIDE-7, fare $5.00
17:11:21.462 INFO  Ana Lopez  | queue empty and closed, exiting
17:11:21.647 INFO  Chen Wei   | completed RIDE-6, fare $3.73
17:11:21.647 INFO  Chen Wei   | queue empty and closed, exiting
17:11:22.163 INFO  Ben Carter | completed RIDE-5, fare $50.00
17:11:22.164 INFO  Ben Carter | queue empty and closed, exiting
17:11:22.164 INFO  main       | all drivers finished
--- Summary ---
Rides completed: 7 | Rides failed: 2 | Total fares: $108.43
Ana Lopez   2 rides  $25.00
Ben Carter  2 rides  $54.40
Chen Wei    3 rides  $29.03
Results file: results.txt (7 rides)
```

The results file of the same run:

```
ride_id,driver,rider,ride_type,miles,fare
RIDE-2,Chen Wei,James Smith,Premium,4.5,18.50
RIDE-3,Ben Carter,Maria Garcia,Shared,6.0,4.40
RIDE-4,Chen Wei,James Smith,Standard,3.2,6.80
RIDE-1,Ana Lopez,Maria Garcia,Standard,12.0,20.00
RIDE-7,Ana Lopez,Maria Garcia,Standard,2.0,5.00
RIDE-6,Chen Wei,James Smith,Shared,8.0,3.73
RIDE-5,Ben Carter,Maria Garcia,Premium,15.0,50.00
```

### Go

```
=== Concurrent Ride Sharing System ===
17:24:43.090 INFO  main       | starting 3 drivers, results file results.txt
17:24:43.091 INFO  Chen Wei   | started
17:24:43.091 INFO  Chen Wei   | took RIDE-1
17:24:43.091 INFO  Ana Lopez  | started
17:24:43.091 INFO  Ana Lopez  | took RIDE-2
17:24:43.091 INFO  Ben Carter | started
17:24:43.091 INFO  Ben Carter | took RIDE-8
17:24:43.091 WARN  Ben Carter | failed RIDE-8: distance must be greater than 0
17:24:43.091 INFO  Ben Carter | took RIDE-3
17:24:43.552 INFO  Ana Lopez  | completed RIDE-2, fare $18.50
17:24:43.552 INFO  Ana Lopez  | took RIDE-4
17:24:43.552 INFO  main       | queue closed
17:24:43.691 INFO  Ben Carter | completed RIDE-3, fare $4.40
17:24:43.691 INFO  Ben Carter | took RIDE-9
17:24:43.691 WARN  Ben Carter | failed RIDE-9: a shared ride needs 2 or more passengers
17:24:43.691 INFO  Ben Carter | took RIDE-5
17:24:43.873 INFO  Ana Lopez  | completed RIDE-4, fare $6.80
17:24:43.873 INFO  Ana Lopez  | took RIDE-6
17:24:44.291 INFO  Chen Wei   | completed RIDE-1, fare $20.00
17:24:44.291 INFO  Chen Wei   | took RIDE-7
17:24:44.492 INFO  Chen Wei   | completed RIDE-7, fare $5.00
17:24:44.492 INFO  Chen Wei   | queue empty and closed, exiting
17:24:44.673 INFO  Ana Lopez  | completed RIDE-6, fare $3.73
17:24:44.673 INFO  Ana Lopez  | queue empty and closed, exiting
17:24:45.192 INFO  Ben Carter | completed RIDE-5, fare $50.00
17:24:45.192 INFO  Ben Carter | queue empty and closed, exiting
17:24:45.192 INFO  main       | all drivers finished
--- Summary ---
Rides completed: 7 | Rides failed: 2 | Total fares: $108.43
Ana Lopez   3 rides  $29.03
Ben Carter  2 rides  $54.40
Chen Wei    2 rides  $25.00
Results file: results.txt (7 rides)
```

The results file of the Go run has the same format as the results file of the Java run.

## Project layout

```
Assignment-6/
  Dockerfile                     the Java program and the Go program in one image
  README.md
  java/
    pom.xml                      Java 21, JUnit 5, and the jar target/rideshare.jar
    src/main/java/rideshare/
      Ride.java                  the ride, the fare rules, and the validation
      RideType.java              Standard, Premium, and Shared
      InvalidRideException.java  the checked exception for an invalid ride
      RideQueue.java             the ride queue
      ResultStore.java           the results list and the results file, under one lock
      RideResult.java            one completed ride
      Driver.java                the driver loop
      RideSharingSystem.java     starts the drivers, adds the rides, waits, and builds the summary
      LogFormat.java             the log line format
      Main.java                  the flags, the demo rides, and the exit status
    src/test/java/rideshare/
      FareTest.java              the fares and the validation
      RideQueueTest.java         the ride queue
      RideSharingSystemTest.java full runs of the system
  go/
    go.mod                       module rideshare-go, Go 1.24
    ride.go                      the ride, the ride types, the fare rules, and the validation
    queue.go                     the ride queue over a buffered channel
    store.go                     the results list and the results file, under one mutex
    driver.go                    the driver loop
    system.go                    starts the drivers, adds the rides, waits, and builds the summary
    logging.go                   the log line format
    main.go                      the flags, the demo rides, and the exit status
    ride_test.go                 the fares and the validation
    queue_test.go                the ride queue
    system_test.go               full runs of the system
```
