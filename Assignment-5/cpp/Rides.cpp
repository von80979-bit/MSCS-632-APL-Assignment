#include "Rides.h"

#include <cmath>
#include <iomanip>
#include <sstream>
#include <utility>

namespace {

constexpr double kStandardBookingFee = 2.00;
constexpr double kStandardPricePerMile = 1.50;
constexpr double kPremiumBookingFee = 5.00;
constexpr double kPremiumPricePerMile = 3.00;
constexpr double kSharedRideDiscount = 0.20;

constexpr int kRideTypeColumnWidth = 18;
constexpr int kRouteColumnWidth = 29;
constexpr int kDistanceColumnWidth = 4;
constexpr int kMoneyColumnWidth = 6;

std::string statusName(RideStatus status) {
    switch (status) {
        case RideStatus::Requested: return "Requested";
        case RideStatus::Assigned: return "Assigned";
        case RideStatus::Completed: return "Completed";
    }
    return "";
}

}  // namespace

// Rounds to whole cents and prints the digits directly, so the output matches the Smalltalk program exactly.
std::string formatMoney(double amountInDollars) {
    const long totalCents = std::lround(amountInDollars * 100);
    std::ostringstream text;
    text << '$' << totalCents / 100 << '.' << std::setw(2) << std::setfill('0') << totalCents % 100;
    return text.str();
}

std::string formatOneDecimal(double value) {
    const long totalTenths = std::lround(value * 10);
    return std::to_string(totalTenths / 10) + '.' + std::to_string(totalTenths % 10);
}

Ride::Ride(std::string rideID, std::string pickupLocation, std::string dropoffLocation, double distance)
    : distance_(distance),
      rideID_(std::move(rideID)),
      pickupLocation_(std::move(pickupLocation)),
      dropoffLocation_(std::move(dropoffLocation)) {}

double Ride::fare() const {
    return kStandardBookingFee + kStandardPricePerMile * distance_;
}

std::string Ride::typeName() const {
    return "Ride";
}

std::string Ride::rideDetails() const {
    std::ostringstream line;
    line << rideID_ << " | "
         << std::left << std::setw(kRideTypeColumnWidth) << typeName() << " | "
         << std::setw(kRouteColumnWidth) << pickupLocation_ + " -> " + dropoffLocation_ << " | "
         << std::right << std::setw(kDistanceColumnWidth) << formatOneDecimal(distance_) << " mi | "
         << std::setw(kMoneyColumnWidth) << formatMoney(fare()) << " | "
         << statusName(status_);
    return line.str();
}

void Ride::assignTo(const std::string& driverName) {
    if (status_ != RideStatus::Requested) return;
    driverName_ = driverName;
    status_ = RideStatus::Assigned;
}

void Ride::complete() {
    if (status_ != RideStatus::Assigned) return;
    status_ = RideStatus::Completed;
}

bool Ride::isCompleted() const {
    return status_ == RideStatus::Completed;
}

const std::string& Ride::rideID() const {
    return rideID_;
}

const std::string& Ride::pickupLocation() const {
    return pickupLocation_;
}

const std::string& Ride::dropoffLocation() const {
    return dropoffLocation_;
}

double Ride::distance() const {
    return distance_;
}

double StandardRide::fare() const {
    return Ride::fare();
}

std::string StandardRide::typeName() const {
    return "Standard";
}

double PremiumRide::fare() const {
    return kPremiumBookingFee + kPremiumPricePerMile * distance_;
}

std::string PremiumRide::typeName() const {
    return "Premium";
}

SharedRide::SharedRide(std::string rideID, std::string pickupLocation, std::string dropoffLocation, double distance,
                       int passengerCount)
    : StandardRide(std::move(rideID), std::move(pickupLocation), std::move(dropoffLocation), distance),
      passengerCount_(passengerCount) {}

double SharedRide::fare() const {
    return StandardRide::fare() / passengerCount_ * (1 - kSharedRideDiscount);
}

std::string SharedRide::typeName() const {
    return "Shared (" + std::to_string(passengerCount_) + " riders)";
}
