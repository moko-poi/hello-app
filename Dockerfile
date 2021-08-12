##
## Build
##

FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN go build -v -o main

##
## Deploy
##

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /app/main /app/main

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/main"]
