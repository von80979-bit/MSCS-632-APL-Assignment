#include "RideSharingSystem.h"

#include <iostream>

void RideSharingSystem::addDriver(const std::shared_ptr<Driver>& driver) {
    drivers_.push_back(driver);
}

void RideSharingSystem::addRider(const std::shared_ptr<Rider>& rider) {
    riders_.push_back(rider);
}

void RideSharingSystem::requestRide(const std::shared_ptr<Rider>& rider, const std::shared_ptr<Ride>& ride) {
    rider->requestRide(ride);
    rides_.push_back(ride);
    const auto driver = nextDriverInTurn();
    matchedDriverByRideID_[ride->rideID()] = driver;
    ride->assignTo(driver->name());
    std::cout << ride->rideID() << " requested by " << rider->name() << ", assigned to " << driver->name() << '\n';
}

void RideSharingSystem::completeRide(const std::shared_ptr<Ride>& ride) {
    const auto matchedDriver = matchedDriverByRideID_.find(ride->rideID());
    if (matchedDriver == matchedDriverByRideID_.end()) return;
    const auto driver = matchedDriver->second;
    ride->complete();
    driver->addRide(ride);
    matchedDriverByRideID_.erase(matchedDriver);
    std::cout << ride->rideID() << " completed by " << driver->name() << '\n';
}

void RideSharingSystem::printReport() const {
    std::cout << "\n--- All rides ---\n";
    for (const auto& ride : rides_) std::cout << ride->rideDetails() << '\n';

    std::cout << "\n--- Riders ---\n";
    for (const auto& rider : riders_) rider->viewRides();

    std::cout << "\n--- Drivers ---\n";
    for (const auto& driver : drivers_) driver->getDriverInfo();
}

std::shared_ptr<Driver> RideSharingSystem::nextDriverInTurn() {
    const auto driver = drivers_[nextDriverIndex_];
    nextDriverIndex_ = (nextDriverIndex_ + 1) % drivers_.size();
    return driver;
}
