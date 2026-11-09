FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /vanta-api ./cmd/api
FROM alpine:3.20
WORKDIR /app
COPY --from=build /vanta-api /app/vanta-api
COPY --from=build /app/migrations /app/migrations
EXPOSE 8080
CMD ["/app/vanta-api"]
