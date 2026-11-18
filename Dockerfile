# syntax=docker/dockerfile:1

# ---- build ----
FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/go-geojson-api .

# ---- runtime ----
FROM alpine:3.20
RUN addgroup -S app && adduser -S app -G app
USER app
COPY --from=build /out/go-geojson-api /usr/local/bin/go-geojson-api
EXPOSE 8080
ENV PORT=8080
ENTRYPOINT ["/usr/local/bin/go-geojson-api"]
