Eval [
    | system anaLopez benCarter chenWei mariaGarcia jamesSmith rides ridersInRequestOrder rideNotCompleted |
    system := RideSharingSystem new.

    anaLopez := Driver driverID: 'DRV-1' name: 'Ana Lopez'.
    benCarter := Driver driverID: 'DRV-2' name: 'Ben Carter'.
    chenWei := Driver driverID: 'DRV-3' name: 'Chen Wei'.
    system addDriver: anaLopez; addDriver: benCarter; addDriver: chenWei.

    mariaGarcia := Rider riderID: 'RDR-1' name: 'Maria Garcia'.
    jamesSmith := Rider riderID: 'RDR-2' name: 'James Smith'.
    system addRider: mariaGarcia; addRider: jamesSmith.

    rides := {
        StandardRide rideID: 'RIDE-1' from: 'Downtown' to: 'Airport' distance: 12.0.
        PremiumRide rideID: 'RIDE-2' from: 'Hotel' to: 'Convention Center' distance: 4.5.
        SharedRide rideID: 'RIDE-3' from: 'University' to: 'Stadium' distance: 6.0 passengers: 2.
        StandardRide rideID: 'RIDE-4' from: 'Mall' to: 'Train Station' distance: 3.2.
        PremiumRide rideID: 'RIDE-5' from: 'Airport' to: 'Harbor' distance: 15.0.
        SharedRide rideID: 'RIDE-6' from: 'Library' to: 'Park' distance: 8.0 passengers: 3
    }.
    ridersInRequestOrder := {mariaGarcia. jamesSmith. mariaGarcia. jamesSmith. mariaGarcia. jamesSmith}.

    Transcript showCr: '=== Ride Sharing System ==='; cr; showCr: '--- Matching ---'.
    rides with: ridersInRequestOrder do: [:ride :rider | system requestRide: ride by: rider].
    rides do: [:ride | system completeRide: ride].

    rideNotCompleted := StandardRide rideID: 'RIDE-7' from: 'Museum' to: 'Zoo' distance: 2.0.
    (anaLopez addRide: rideNotCompleted) ifFalse: [
        Transcript showCr: anaLopez driverID, ' rejected ', rideNotCompleted rideID, ': the ride is not completed.'].

    system printReport
]
