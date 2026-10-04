#pragma once

#include <memory>
#include <string>
#include <vector>

#include "Rides.h"

class Driver {
public:
    Driver(std::string driverID, std::string name);

    bool addRide(const std::shared_ptr<Ride>& ride);
    double totalEarnings() const;
    void getDriverInfo() const;

    const std::string& driverID() const;
    const std::string& name() const;

private:
    std::string driverID_;
    std::string name_;
    double rating_;
    std::vector<std::shared_ptr<Ride>> assignedRides_;
};
