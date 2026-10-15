# The family tree program on the official SWI-Prolog image. `docker run --rm <image>` runs the demo.
# `docker run --rm -it <image> swipl family.pl` opens the prompt for queries by hand.

FROM swipl:10.0.2
WORKDIR /app
COPY prolog/family.pl .
CMD ["swipl", "-g", "demo", "-t", "halt", "family.pl"]
