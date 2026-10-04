Object subclass: Driver [
    | driverID name rating assignedRides |

    Driver class >> driverID: anID name: aName [
        ^self new setDriverID: anID name: aName
    ]

    setDriverID: anID name: aName [
        driverID := anID.
        name := aName.
        rating := (Random between: 30 and: 50) / 10.
        assignedRides := OrderedCollection new
    ]

    addRide: aRide [
        aRide isCompleted ifFalse: [^false].
        assignedRides add: aRide.
        ^true
    ]

    totalEarnings [
        ^assignedRides inject: 0 into: [:total :ride | total + ride fare]
    ]

    getDriverInfo [
        Transcript showCr: driverID, ' ', name, ' | Rating ', rating asOneDecimalString, ' | ',
            assignedRides size printString, ' completed rides'.
        assignedRides do: [:ride | Transcript showCr: '  ', ride rideDetails].
        Transcript showCr: '  Total earned: ', self totalEarnings asMoneyString
    ]

    driverID [
        ^driverID
    ]

    name [
        ^name
    ]
]
