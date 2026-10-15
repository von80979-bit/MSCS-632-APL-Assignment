# One image with the three statistics calculator programs. `docker run --rm <image>` runs all three in order.
# Run one with `docker run --rm <image> c`, `ocaml`, or `python`.

FROM ocaml/opam:debian-12-ocaml-5.2 AS ocaml-build
WORKDIR /home/opam/build
COPY --chown=opam ocaml/stats.ml .
RUN opam exec -- ocamlopt -o stats-ocaml stats.ml

FROM gcc:14 AS c-build
WORKDIR /build
COPY c/stats.c .
RUN gcc -std=c11 -Wall -Wextra -Werror -O2 -o stats-c stats.c

# The C and OCaml binaries link only against glibc, so a Python image on the same Debian release runs all three.
FROM python:3.12-slim-bookworm
COPY --from=c-build /build/stats-c /usr/local/bin/stats-c
COPY --from=ocaml-build /home/opam/build/stats-ocaml /usr/local/bin/stats-ocaml
COPY python/stats.py /app/stats.py
RUN printf '%s\n' '#!/bin/sh' \
        'case "${1:-all}" in' \
        '    c) exec stats-c ;;' \
        '    ocaml) exec stats-ocaml ;;' \
        '    python) exec python3 /app/stats.py ;;' \
        '    all) stats-c && echo && stats-ocaml && echo && python3 /app/stats.py ;;' \
        '    *) echo "usage: docker run --rm <image> [c|ocaml|python]" >&2; exit 2 ;;' \
        'esac' \
        > /usr/local/bin/stats \
    && chmod +x /usr/local/bin/stats
ENTRYPOINT ["stats"]
