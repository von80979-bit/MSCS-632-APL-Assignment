# One Linux environment for every part of MSCS-632 Assignment 2.
#
# Two reasons this image exists. valgrind has no working port on macOS 15, so the C++ leak and the dangling
# read cannot be proven on the host at all. The stronger reason is comparability: measuring a JVM against a
# native binary is only honest when both ran on the same operating system and processor, and the host Mac
# runs its toolchain under Rosetta.
#
# No platform is pinned, so the image matches whatever the Docker daemon offers. On the machine this was
# built on that is linux/amd64, because Docker Desktop's Linux virtual machine is the Intel build. Absolute
# timings are therefore not host-native, but every language sits in the same environment, which is the
# property the memory comparison actually needs.
FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive
# matplotlib has no display attached here, so psrecord needs the file-writing backend to produce its PNG.
ENV MPLBACKEND=Agg

RUN apt-get update && apt-get install --yes --no-install-recommends \
        g++ \
        rustc \
        openjdk-21-jdk-headless \
        python3 \
        python3-pip \
        python3-matplotlib \
        python3-psutil \
        nodejs \
        valgrind \
        ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Ubuntu marks its system Python as externally managed, so pip needs the override. psrecord is the one
# profiling tool that measures all three languages the same way and the distribution does not package it.
RUN pip3 install --break-system-packages --no-cache-dir psrecord

# Baked in as well as mounted, so a grader can read the environment straight off the image with no mount.
COPY toolchain-versions.sh /usr/local/bin/toolchain-versions
RUN chmod +x /usr/local/bin/toolchain-versions

# container.sh bind-mounts the repository root at /workspace, so host edits are visible without a rebuild.
WORKDIR /workspace/Assignment-2/code

CMD ["/bin/bash"]
