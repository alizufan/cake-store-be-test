## Build
FROM golang:1.19 AS build

WORKDIR /go/src/cake-store

COPY . .

RUN go mod download
RUN go vet -v
RUN go test -v ./...

RUN CGO_ENABLED=0 go build -o /go/bin/cake-store

## Deploy using distroless to squeeze image size (~2MB)

## Dev
# FROM gcr.io/distroless/static-debian11:debug

## Prod
FROM gcr.io/distroless/static-debian11

WORKDIR /

COPY --from=build /go/bin/cake-store /

ENTRYPOINT ["/cake-store"]