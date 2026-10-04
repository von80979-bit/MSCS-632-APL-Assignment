#include "Rider.h"

#include <iostream>
#include <utility>

Rider::Rider(std::string riderID, std::string name) : riderID_(std::move(riderID)), name_(std::move(name)) {}

void Rider::requestRide(const std::shared_ptr<Ride>& ride) {
    requestedRides_.push_back(ride);
}

double Rider::totalSpending() const {
    double total = 0;
    for (const auto& ride : requestedRides_) total += ride->fare();
    return total;
}

void Rider::viewRides() const {
    std::cout << riderID_ << ' ' << name_ << '\n';
    for (const auto& ride : requestedRides_) std::cout << "  " << ride->rideDetails() << '\n';
    std::cout << "  Total spent: " << formatMoney(totalSpending()) << '\n';
}

const std::string& Rider::name() const {
    return name_;
}
