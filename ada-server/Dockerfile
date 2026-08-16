FROM golang:1.25-alpine AS build
WORKDIR /src

# Cache module downloads: copy manifests first.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source and build a static binary.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /ada-server ./ada-server

# Minimal runtime image.
FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=build /ada-server /ada-server

ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/ada-server"]
