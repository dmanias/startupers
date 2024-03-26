# docker build -t ghcr.io/dmanias/creativenarration -f ./deployments/Dockerfile --target development .

ARG GO_VERSION=1.21.6
FROM golang:${GO_VERSION}-alpine AS development

WORKDIR /app
COPY . ./

RUN go mod download
CMD ["go", "run", "./app/main.go"]

# docker build -t ghcr.io/dmanias/creativenarration -f ./deployments/Dockerfile --target testing .
FROM development AS testing
RUN go test -v ./...

FROM development AS build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o creativenarration ./app/main.go

FROM gcr.io/distroless/cc-debian12@sha256:899570acf85a1f1362862a9ea4d9e7b1827cb5c62043ba5b170b21de89618608 AS production
COPY --from=build /app/creativenarration /bin/creativenarration
CMD ["/bin/creativenarration"]
