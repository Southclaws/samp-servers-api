# -
# Build workspace
# -
FROM golang:1.11 AS compile

RUN apt-get update -y && apt-get install --no-install-recommends -y -q build-essential ca-certificates

WORKDIR /samp-servers-api
ADD . .
RUN make static

# -
# Runtime
# -
FROM scratch

COPY --from=compile /samp-servers-api/samp-servers-api /bin/samp-servers-api
COPY --from=compile /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["samp-servers-api"]
