# Base build image
FROM golang:1.11-alpine AS build_base

# Install some dependencies needed to build the project
RUN apk add bash ca-certificates git gcc g++ libc-dev
# Cache dependeny

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o azup

FROM alpine AS dist
RUN apk add ca-certificates
WORKDIR /app
COPY --from=build_base /app/azup /app/azup
ENTRYPOINT ["/app/azup"]
