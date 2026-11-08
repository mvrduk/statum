ARG GO_VERSION=1.27.3

FROM golang:${GO_VERSION} AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/api ./cmd/api
RUN CGO_ENABLED=0 go build -o /out/migrator ./cmd/migrator

FROM gcr.io/distroless/static-debian12:nonroot AS api
WORKDIR /app
COPY --from=builder /out/api /app/api
COPY config/app.ini.example /app/config/app.ini
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app/api"]

FROM gcr.io/distroless/static-debian12:nonroot AS migrator
WORKDIR /app
COPY --from=builder /out/migrator /app/migrator
COPY migrations /app/migrations
COPY config/app.ini.example /app/config/app.ini
USER nonroot:nonroot
ENTRYPOINT ["/app/migrator"]
CMD ["up"]