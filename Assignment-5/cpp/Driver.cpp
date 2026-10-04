#include "Driver.h"

#include <iostream>
#include <random>
#include <utility>

namespace {

double randomRating() {
    static std::mt19937 generator{std::random_device{}()};
    std::uniform_int_distribution<int> ratingInTenths(30, 50);
    return ratingInTenths(generator) / 10.0;
}

}  // namespace

Driver::Driver(std::string driverID, std::string name)
    : driverID_(std::move(driverID)), name_(std::move(name)), rating_(randomRating()) {}

bool Driver::addRide(const std::shared_ptr<Ride>& ride) {
    if (!ride->isCompleted()) return false;
    assignedRides_.push_back(ride);
    return true;
}

double Driver::totalEarnings() const {
    double total = 0;
    for (const auto& ride : assignedRides_) total += ride->fare();
    return total;
}

void Driver::getDriverInfo() const {
    std::cout << driverID_ << ' ' << name_ << " | Rating " << formatOneDecimal(rating_) << " | "
              << assignedRides_.size() << " completed rides\n";
    for (const auto& ride : assignedRides_) std::cout << "  " << ride->rideDetails() << '\n';
    std::cout << "  Total earned: " << formatMoney(totalEarnings()) << '\n';
}

const std::string& Driver::driverID() const {
    return driverID_;
}

const std::string& Driver::name() const {
    return name_;
}
