# One image with both ride sharing programs. Run one with `docker run --rm <image> smalltalk` or `cpp`.

FROM ubuntu:22.04 AS cpp-build
RUN apt-get update && apt-get install -y --no-install-recommends g++ && rm -rf /var/lib/apt/lists/*
WORKDIR /build
COPY cpp/ .
RUN g++ -std=c++17 -Wall -Wextra -o ride_sharing *.cpp

# Ubuntu 22.04 is the newest release with a GNU Smalltalk package.
FROM ubuntu:22.04
RUN apt-get update && apt-get install -y --no-install-recommends gnu-smalltalk && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY smalltalk/ smalltalk/
COPY --from=cpp-build /build/ride_sharing /usr/local/bin/cpp
RUN printf '#!/bin/sh\ncd /app/smalltalk && exec gst Rides.st Driver.st Rider.st RideSharingSystem.st main.st\n' \
        > /usr/local/bin/smalltalk \
    && chmod +x /usr/local/bin/smalltalk
