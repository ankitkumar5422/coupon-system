FROM golang:1.20-alpine AS builder

# Install build dependencies for SQLite
RUN apk add --no-cache --update gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build with CGO enabled for SQLite support
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o coupon-system ./cmd/server/main.go

FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache --update ca-certificates sqlite-libs

WORKDIR /root/

COPY --from=builder /app/coupon-system .

EXPOSE 8080

CMD ["./coupon-system"]