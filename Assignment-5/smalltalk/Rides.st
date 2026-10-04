Number extend [
    asMoneyString [
        "Rounds to whole cents and prints the digits directly, so the output matches the C++ program exactly."

        | totalCents |
        totalCents := (self * 100) rounded.
        ^'$', (totalCents // 100) printString, '.', ((totalCents \\ 100) printPaddedWith: $0 to: 2)
    ]

    asOneDecimalString [
        | totalTenths |
        totalTenths := (self * 10) rounded.
        ^(totalTenths // 10) printString, '.', (totalTenths \\ 10) printString
    ]
]

String extend [
    paddedRightTo: width [
        ^self, (String new: (width - self size max: 0) withAll: $ )
    ]

    paddedLeftTo: width [
        ^(String new: (width - self size max: 0) withAll: $ ), self
    ]
]

Object subclass: Ride [
    | rideID pickupLocation dropoffLocation distance status driverName |

    Ride class >> rideID: anID from: aPickupLocation to: aDropoffLocation distance: miles [
        ^self new setRideID: anID from: aPickupLocation to: aDropoffLocation distance: miles
    ]

    setRideID: anID from: aPickupLocation to: aDropoffLocation distance: miles [
        rideID := anID.
        pickupLocation := aPickupLocation.
        dropoffLocation := aDropoffLocation.
        distance := miles.
        status := #Requested.
        driverName := ''
    ]

    fare [
        ^2.00 + (1.50 * distance)
    ]

    typeName [
        ^'Ride'
    ]

    rideDetails [
        ^rideID, ' | ',
            (self typeName paddedRightTo: 18), ' | ',
            (pickupLocation, ' -> ', dropoffLocation paddedRightTo: 29), ' | ',
            (distance asOneDecimalString paddedLeftTo: 4), ' mi | ',
            (self fare asMoneyString paddedLeftTo: 6), ' | ',
            status asString
    ]

    assignTo: aDriverName [
        status = #Requested ifFalse: [^self].
        driverName := aDriverName.
        status := #Assigned
    ]

    complete [
        status = #Assigned ifFalse: [^self].
        status := #Completed
    ]

    isCompleted [
        ^status = #Completed
    ]

    rideID [
        ^rideID
    ]

    pickupLocation [
        ^pickupLocation
    ]

    dropoffLocation [
        ^dropoffLocation
    ]

    distance [
        ^distance
    ]
]

Ride subclass: StandardRide [
    fare [
        ^super fare
    ]

    typeName [
        ^'Standard'
    ]
]

Ride subclass: PremiumRide [
    fare [
        ^5.00 + (3.00 * distance)
    ]

    typeName [
        ^'Premium'
    ]
]

StandardRide subclass: SharedRide [
    | passengerCount |

    SharedRide class >> rideID: anID from: aPickupLocation to: aDropoffLocation distance: miles
            passengers: aPassengerCount [
        ^(self rideID: anID from: aPickupLocation to: aDropoffLocation distance: miles)
            setPassengerCount: aPassengerCount
    ]

    setPassengerCount: aPassengerCount [
        passengerCount := aPassengerCount
    ]

    fare [
        "Each passenger pays an equal share of the standard fare, with 20% off."

        ^super fare / passengerCount * 0.80
    ]

    typeName [
        ^'Shared (', passengerCount printString, ' riders)'
    ]
]
