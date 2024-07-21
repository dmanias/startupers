# Development stage
ARG GO_VERSION=1.22
FROM golang:${GO_VERSION}-alpine AS development

WORKDIR /app
#COPY internal/docker .
COPY . .

RUN go mod download

# Command to run the application in development mode
CMD ["go", "run", "app/main.go"]

# Debug stage
FROM development AS debug

# Install Delve for debugging
RUN go install github.com/go-delve/delve/cmd/dlv@latest

ENTRYPOINT [ "dlv", "--listen=:2345", "--api-version=2", "--headless", "--accept-multiclient", "debug", "app/main.go" ]

FROM development AS testing
RUN go test -v ./...

FROM development AS build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o startupers app/main.go

FROM gcr.io/distroless/cc-debian12@sha256:899570acf85a1f1362862a9ea4d9e7b1827cb5c62043ba5b170b21de89618608 AS production
COPY --from=build /app/startupers /bin/startupers
CMD ["/bin/startupers"]
