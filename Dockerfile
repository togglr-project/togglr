FROM golang:1.25-alpine AS build

ENV GOPROXY="https://proxy.golang.org,direct"
ENV PROJECTDIR=/src
ENV CGO_ENABLED=0

RUN apk add --no-cache make

WORKDIR ${PROJECTDIR}
COPY go.mod go.sum ${PROJECTDIR}/
RUN go mod download

COPY . ${PROJECTDIR}/

RUN make build

# Production image
FROM scratch AS prod

COPY --from=build /src/bin/app /bin/app
COPY --from=build /src/migrations /migrations

CMD ["/bin/app"]
