# Assignment 5: Ride sharing system

Two class-based ride sharing programs built in GNU Smalltalk and C++17. Both programs show encapsulation, inheritance, and polymorphism, and both print the same output.

## What the application does

The program runs one fixed demo script and exits. It takes no input.

The script creates three drivers, two riders, and six rides of three ride types: standard, premium, and shared. Each rider requests rides, and the system matches each ride to a driver. Drivers take turns in a fixed order. Each driver then completes the rides that the system matched to them.

The program prints four parts:

- The matching, which shows each request, each completion, and one ride that a driver rejects because it is not completed.
- All rides, from one list that holds the three ride types. Each ride calculates its own fare.
- Each rider's ride history and total spending.
- Each driver's details, completed rides, and total earnings.

## Fares

| Ride type | Fare |
|---|---|
| Standard | $2.00 booking fee + $1.50 per mile |
| Premium | $5.00 booking fee + $3.00 per mile |
| Shared | The standard fare divided by the number of passengers, with 20% off |

## Object-oriented principles

- **Encapsulation.** Each class keeps its state private and changes it only through its own methods. A driver accepts a ride into its list only when the ride is completed.
- **Inheritance.** `StandardRide` and `PremiumRide` inherit from `Ride`, and `SharedRide` inherits from `StandardRide`.
- **Polymorphism.** One list holds rides of all three ride types. Each ride answers `fare()` with the formula of its own ride type.

## Assumptions

The requirement leaves certain points open, so the two programs share these assumptions:

1. Distance is in miles, and money is in US dollars.
2. A ride calculates its fare each time it is asked, from its distance and its ride type. A stored fare could go out of date, so the program stores none.
3. A driver's ride list holds completed rides only, because the requirement defines it as "rides completed by the driver". The driver rejects a ride that is not completed.
4. A shared ride is one rider's seat in a shared trip. Its fare is that rider's share, and the driver earns the sum of the seats.
5. The rating of a driver is random at creation, from 3.0 to 5.0 with one decimal. The rating does not affect matching or fares, so it is the only line of output that changes between runs.

## Run with Docker

Docker is the primary run path. One image holds both programs. Run these commands in this folder.

Build the image:

```sh
docker build -t ride-sharing .
```

Run the Smalltalk program:

```sh
docker run --rm ride-sharing smalltalk
```

Run the C++ program:

```sh
docker run --rm ride-sharing cpp
```

## Run without Docker

The Smalltalk program needs GNU Smalltalk 3.2 or later. Run it in the `smalltalk` folder:

```sh
cd smalltalk
gst Rides.st Driver.st Rider.st RideSharingSystem.st main.st
```

The C++ program needs a compiler with C++17 support. Build and run it in the `cpp` folder:

```sh
cd cpp
g++ -std=c++17 -Wall -Wextra -o ride_sharing *.cpp
./ride_sharing
```

## Sample output

Both programs print this output. Only the driver ratings change between runs.

```
=== Ride Sharing System ===

--- Matching ---
RIDE-1 requested by Maria Garcia, assigned to Ana Lopez
RIDE-2 requested by James Smith, assigned to Ben Carter
RIDE-3 requested by Maria Garcia, assigned to Chen Wei
RIDE-4 requested by James Smith, assigned to Ana Lopez
RIDE-5 requested by Maria Garcia, assigned to Ben Carter
RIDE-6 requested by James Smith, assigned to Chen Wei
RIDE-1 completed by Ana Lopez
RIDE-2 completed by Ben Carter
RIDE-3 completed by Chen Wei
RIDE-4 completed by Ana Lopez
RIDE-5 completed by Ben Carter
RIDE-6 completed by Chen Wei
DRV-1 rejected RIDE-7: the ride is not completed.

--- All rides ---
RIDE-1 | Standard           | Downtown -> Airport           | 12.0 mi | $20.00 | Completed
RIDE-2 | Premium            | Hotel -> Convention Center    |  4.5 mi | $18.50 | Completed
RIDE-3 | Shared (2 riders)  | University -> Stadium         |  6.0 mi |  $4.40 | Completed
RIDE-4 | Standard           | Mall -> Train Station         |  3.2 mi |  $6.80 | Completed
RIDE-5 | Premium            | Airport -> Harbor             | 15.0 mi | $50.00 | Completed
RIDE-6 | Shared (3 riders)  | Library -> Park               |  8.0 mi |  $3.73 | Completed

--- Riders ---
RDR-1 Maria Garcia
  RIDE-1 | Standard           | Downtown -> Airport           | 12.0 mi | $20.00 | Completed
  RIDE-3 | Shared (2 riders)  | University -> Stadium         |  6.0 mi |  $4.40 | Completed
  RIDE-5 | Premium            | Airport -> Harbor             | 15.0 mi | $50.00 | Completed
  Total spent: $74.40
RDR-2 James Smith
  RIDE-2 | Premium            | Hotel -> Convention Center    |  4.5 mi | $18.50 | Completed
  RIDE-4 | Standard           | Mall -> Train Station         |  3.2 mi |  $6.80 | Completed
  RIDE-6 | Shared (3 riders)  | Library -> Park               |  8.0 mi |  $3.73 | Completed
  Total spent: $29.03

--- Drivers ---
DRV-1 Ana Lopez | Rating 4.1 | 2 completed rides
  RIDE-1 | Standard           | Downtown -> Airport           | 12.0 mi | $20.00 | Completed
  RIDE-4 | Standard           | Mall -> Train Station         |  3.2 mi |  $6.80 | Completed
  Total earned: $26.80
DRV-2 Ben Carter | Rating 4.3 | 2 completed rides
  RIDE-2 | Premium            | Hotel -> Convention Center    |  4.5 mi | $18.50 | Completed
  RIDE-5 | Premium            | Airport -> Harbor             | 15.0 mi | $50.00 | Completed
  Total earned: $68.50
DRV-3 Chen Wei | Rating 4.3 | 2 completed rides
  RIDE-3 | Shared (2 riders)  | University -> Stadium         |  6.0 mi |  $4.40 | Completed
  RIDE-6 | Shared (3 riders)  | Library -> Park               |  8.0 mi |  $3.73 | Completed
  Total earned: $8.13
```

## Project layout

```
Assignment-5/
  Dockerfile               both programs in one image
  README.md
  smalltalk/
    Rides.st               Ride, StandardRide, PremiumRide, SharedRide, and the money format
    Driver.st              Driver
    Rider.st               Rider
    RideSharingSystem.st   the matching, the completions, and the report
    main.st                the demo script
  cpp/
    Rides.h  Rides.cpp     the same five parts as the Smalltalk program, one header and one source file each
    Driver.h  Driver.cpp
    Rider.h  Rider.cpp
    RideSharingSystem.h  RideSharingSystem.cpp
    main.cpp
```

Both programs format money from whole cents, so the two outputs match to the cent.
