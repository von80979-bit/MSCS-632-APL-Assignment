#include <cstddef>
#include <iostream>
#include <memory>
#include <vector>

#include "Driver.h"
#include "RideSharingSystem.h"
#include "Rider.h"
#include "Rides.h"

int main() {
    RideSharingSystem system;

    const auto anaLopez = std::make_shared<Driver>("DRV-1", "Ana Lopez");
    const auto benCarter = std::make_shared<Driver>("DRV-2", "Ben Carter");
    const auto chenWei = std::make_shared<Driver>("DRV-3", "Chen Wei");
    system.addDriver(anaLopez);
    system.addDriver(benCarter);
    system.addDriver(chenWei);

    const auto mariaGarcia = std::make_shared<Rider>("RDR-1", "Maria Garcia");
    const auto jamesSmith = std::make_shared<Rider>("RDR-2", "James Smith");
    system.addRider(mariaGarcia);
    system.addRider(jamesSmith);

    const std::vector<std::shared_ptr<Ride>> rides = {
        std::make_shared<StandardRide>("RIDE-1", "Downtown", "Airport", 12.0),
        std::make_shared<PremiumRide>("RIDE-2", "Hotel", "Convention Center", 4.5),
        std::make_shared<SharedRide>("RIDE-3", "University", "Stadium", 6.0, 2),
        std::make_shared<StandardRide>("RIDE-4", "Mall", "Train Station", 3.2),
        std::make_shared<PremiumRide>("RIDE-5", "Airport", "Harbor", 15.0),
        std::make_shared<SharedRide>("RIDE-6", "Library", "Park", 8.0, 3),
    };
    const std::vector<std::shared_ptr<Rider>> ridersInRequestOrder = {
        mariaGarcia, jamesSmith, mariaGarcia, jamesSmith, mariaGarcia, jamesSmith,
    };

    std::cout << "=== Ride Sharing System ===\n\n--- Matching ---\n";
    for (std::size_t index = 0; index < rides.size(); ++index) {
        system.requestRide(ridersInRequestOrder[index], rides[index]);
    }
    for (const auto& ride : rides) system.completeRide(ride);

    const auto rideNotCompleted = std::make_shared<StandardRide>("RIDE-7", "Museum", "Zoo", 2.0);
    if (!anaLopez->addRide(rideNotCompleted)) {
        std::cout << anaLopez->driverID() << " rejected " << rideNotCompleted->rideID()
                  << ": the ride is not completed.\n";
    }

    system.printReport();
    return 0;
}
