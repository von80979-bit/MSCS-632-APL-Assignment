#pragma once

#include <memory>
#include <string>
#include <vector>

#include "Rides.h"

class Rider {
public:
    Rider(std::string riderID, std::string name);

    void requestRide(const std::shared_ptr<Ride>& ride);
    double totalSpending() const;
    void viewRides() const;

    const std::string& name() const;

private:
    std::string riderID_;
    std::string name_;
    std::vector<std::shared_ptr<Ride>> requestedRides_;
};
