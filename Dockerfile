FROM golang:1.25-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /ada-server ./ada-server

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=build /ada-server /ada-server

ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/ada-server"]
