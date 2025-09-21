FROM golang:1.25-alpine AS build

ENV GOPROXY="https://proxy.golang.org,direct"
ENV PROJECTDIR=/src
ENV CGO_ENABLED=0

RUN apk add --no-cache make
RUN apk add --no-cache tzdata

WORKDIR ${PROJECTDIR}
COPY go.mod go.sum ${PROJECTDIR}/
RUN go mod download

COPY . ${PROJECTDIR}/

RUN make build

# Create a minimal image with only curl and its dependencies
FROM alpine:3.19 AS curl-extract

RUN apk add --no-cache curl
RUN apk add --no-cache tzdata

# Extract curl and its dependencies directly in Dockerfile
RUN mkdir -p /curl-deps && \
    cp /usr/bin/curl /curl-deps/ && \
    ldd /usr/bin/curl | grep "=>" | awk '{print $3}' | xargs -I {} cp {} /curl-deps/ && \
    cp /etc/ssl/certs/ca-certificates.crt /curl-deps/

# Production image - using scratch for maximum security
FROM scratch AS prod

# Copy curl and its dependencies from the extract stage
COPY --from=curl-extract /curl-deps/curl /usr/bin/curl
COPY --from=curl-extract /curl-deps/*.so* /lib/
COPY --from=curl-extract /curl-deps/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Copy timezones
COPY --from=curl-extract /usr/share/zoneinfo /usr/share/zoneinfo

# Copy application binary and migrations
COPY --from=build /src/bin/app /bin/app
COPY --from=build /src/migrations /migrations

CMD ["/bin/app"]
