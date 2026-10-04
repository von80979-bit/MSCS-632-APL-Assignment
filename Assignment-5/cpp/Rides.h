#pragma once

#include <string>

enum class RideStatus { Requested, Assigned, Completed };

std::string formatMoney(double amountInDollars);
std::string formatOneDecimal(double value);

class Ride {
public:
    Ride(std::string rideID, std::string pickupLocation, std::string dropoffLocation, double distance);
    virtual ~Ride() = default;

    virtual double fare() const;
    virtual std::string typeName() const;
    std::string rideDetails() const;

    void assignTo(const std::string& driverName);
    void complete();
    bool isCompleted() const;

    const std::string& rideID() const;
    const std::string& pickupLocation() const;
    const std::string& dropoffLocation() const;
    double distance() const;

protected:
    double distance_;

private:
    std::string rideID_;
    std::string pickupLocation_;
    std::string dropoffLocation_;
    RideStatus status_ = RideStatus::Requested;
    std::string driverName_;
};

class StandardRide : public Ride {
public:
    using Ride::Ride;

    double fare() const override;
    std::string typeName() const override;
};

class PremiumRide : public Ride {
public:
    using Ride::Ride;

    double fare() const override;
    std::string typeName() const override;
};

class SharedRide : public StandardRide {
public:
    SharedRide(std::string rideID, std::string pickupLocation, std::string dropoffLocation, double distance,
               int passengerCount);

    double fare() const override;
    std::string typeName() const override;

private:
    int passengerCount_;
};
