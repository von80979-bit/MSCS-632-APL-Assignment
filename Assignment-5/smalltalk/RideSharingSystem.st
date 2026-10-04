Object subclass: RideSharingSystem [
    | drivers riders rides matchedDriverByRideID nextDriverIndex |

    RideSharingSystem class >> new [
        ^super new setUp
    ]

    setUp [
        drivers := OrderedCollection new.
        riders := OrderedCollection new.
        rides := OrderedCollection new.
        matchedDriverByRideID := Dictionary new.
        nextDriverIndex := 1
    ]

    addDriver: aDriver [
        drivers add: aDriver
    ]

    addRider: aRider [
        riders add: aRider
    ]

    requestRide: aRide by: aRider [
        | driver |
        aRider requestRide: aRide.
        rides add: aRide.
        driver := self nextDriverInTurn.
        matchedDriverByRideID at: aRide rideID put: driver.
        aRide assignTo: driver name.
        Transcript showCr: aRide rideID, ' requested by ', aRider name, ', assigned to ', driver name
    ]

    completeRide: aRide [
        | driver |
        driver := matchedDriverByRideID removeKey: aRide rideID ifAbsent: [^self].
        aRide complete.
        driver addRide: aRide.
        Transcript showCr: aRide rideID, ' completed by ', driver name
    ]

    printReport [
        Transcript cr; showCr: '--- All rides ---'.
        rides do: [:ride | Transcript showCr: ride rideDetails].

        Transcript cr; showCr: '--- Riders ---'.
        riders do: [:rider | rider viewRides].

        Transcript cr; showCr: '--- Drivers ---'.
        drivers do: [:driver | driver getDriverInfo]
    ]

    nextDriverInTurn [
        | driver |
        driver := drivers at: nextDriverIndex.
        nextDriverIndex := nextDriverIndex \\ drivers size + 1.
        ^driver
    ]
]
