FROM alpine:latest

# Optional: install bash if you prefer bash over sh
RUN apk add --no-cache bash

WORKDIR /app

# Entry point simply passes any command received to sh
ENTRYPOINT ["/bin/sh", "-c"]
