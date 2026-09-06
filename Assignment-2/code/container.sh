#!/usr/bin/env bash
# Builds and enters the single Linux environment every part of this assignment runs in.
#
#   ./container.sh build            build or rebuild the image
#   ./container.sh versions         print every toolchain version
#   ./container.sh shell            interactive shell inside the container
#   ./container.sh run <command>    run one command inside the container
#
# The repository is bind-mounted rather than copied, so a file edited on the host is already changed inside
# the container and no rebuild is needed between an edit and a run.

set -euo pipefail

IMAGE_NAME=mscs632-assignment2
CODE_DIRECTORY=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
REPOSITORY_ROOT=$(cd "${CODE_DIRECTORY}/../.." && pwd)
WORKING_DIRECTORY_INSIDE_CONTAINER=/workspace/Assignment-2/code

build_image() {
    docker build --tag "${IMAGE_NAME}" "${CODE_DIRECTORY}"
}

run_in_container() {
    # Assembled as one array because the terminal flags are conditional, and macOS still ships bash 3.2,
    # where expanding an empty array under 'set -u' aborts the script.
    local docker_command=(docker run --rm)
    if [ -t 0 ] && [ -t 1 ]; then
        docker_command+=(--interactive --tty)
    fi
    docker_command+=(--volume "${REPOSITORY_ROOT}:/workspace")
    docker_command+=(--workdir "${WORKING_DIRECTORY_INSIDE_CONTAINER}")
    docker_command+=("${IMAGE_NAME}" "$@")

    "${docker_command[@]}"
}

subcommand=${1:-shell}
shift || true

case "${subcommand}" in
    build)
        build_image
        ;;
    versions)
        run_in_container toolchain-versions
        ;;
    shell)
        run_in_container /bin/bash
        ;;
    run)
        if [ "$#" -eq 0 ]; then
            echo "container.sh run needs a command to run" >&2
            exit 2
        fi
        run_in_container "$@"
        ;;
    *)
        echo "unknown subcommand '${subcommand}'; expected build, versions, shell or run" >&2
        exit 2
        ;;
esac
