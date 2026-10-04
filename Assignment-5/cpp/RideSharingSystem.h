#pragma once

#include <cstddef>
#include <map>
#include <memory>
#include <string>
#include <vector>

#include "Driver.h"
#include "Rider.h"
#include "Rides.h"

class RideSharingSystem {
public:
    void addDriver(const std::shared_ptr<Driver>& driver);
    void addRider(const std::shared_ptr<Rider>& rider);
    void requestRide(const std::shared_ptr<Rider>& rider, const std::shared_ptr<Ride>& ride);
    void completeRide(const std::shared_ptr<Ride>& ride);
    void printReport() const;

private:
    std::shared_ptr<Driver> nextDriverInTurn();

    std::vector<std::shared_ptr<Driver>> drivers_;
    std::vector<std::shared_ptr<Rider>> riders_;
    std::vector<std::shared_ptr<Ride>> rides_;
    std::map<std::string, std::shared_ptr<Driver>> matchedDriverByRideID_;
    std::size_t nextDriverIndex_ = 0;
};
