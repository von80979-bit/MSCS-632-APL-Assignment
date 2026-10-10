# One image with both concurrent ride sharing programs. Run one with `docker run --rm <image> java` or `go`.
# Flags after the language name reach the program, for example `docker run --rm <image> java -drivers 5`.

FROM maven:3.9-eclipse-temurin-21 AS java-build
WORKDIR /build
COPY java/pom.xml .
RUN mvn -q dependency:go-offline
COPY java/src/ src/
RUN mvn -q package -DskipTests

FROM golang:1.24 AS go-build
WORKDIR /build
COPY go/ .
# A static binary needs no C library from the final image.
RUN CGO_ENABLED=0 go build -o rideshare-go .

FROM eclipse-temurin:21-jre
WORKDIR /app
COPY --from=java-build /build/target/rideshare.jar /app/rideshare.jar
COPY --from=go-build /build/rideshare-go /usr/local/bin/rideshare-go
# The JDK already owns the command `java`, so one entry script picks the program from the first argument.
RUN printf '%s\n' '#!/bin/sh' \
        'language="$1"' \
        '[ $# -gt 0 ] && shift' \
        'case "$language" in' \
        '    java) exec java -jar /app/rideshare.jar "$@" ;;' \
        '    go) exec rideshare-go "$@" ;;' \
        '    *) echo "usage: docker run --rm <image> java|go [flags]" >&2; exit 2 ;;' \
        'esac' \
        > /usr/local/bin/rideshare \
    && chmod +x /usr/local/bin/rideshare
ENTRYPOINT ["rideshare"]
