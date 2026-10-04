Object subclass: Rider [
    | riderID name requestedRides |

    Rider class >> riderID: anID name: aName [
        ^self new setRiderID: anID name: aName
    ]

    setRiderID: anID name: aName [
        riderID := anID.
        name := aName.
        requestedRides := OrderedCollection new
    ]

    requestRide: aRide [
        requestedRides add: aRide
    ]

    totalSpending [
        ^requestedRides inject: 0 into: [:total :ride | total + ride fare]
    ]

    viewRides [
        Transcript showCr: riderID, ' ', name.
        requestedRides do: [:ride | Transcript showCr: '  ', ride rideDetails].
        Transcript showCr: '  Total spent: ', self totalSpending asMoneyString
    ]

    name [
        ^name
    ]
]
