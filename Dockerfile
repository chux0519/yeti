
FROM golang:1.20-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY pkg ./pkg
COPY cmd ./cmd

RUN go build -o /yeti

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /yeti /yeti

EXPOSE 25700

# USER nonroot:nonroot

ENTRYPOINT ["/yeti"]