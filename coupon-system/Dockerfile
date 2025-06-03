FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o coupon-system ./cmd/server/main.go

FROM scratch

WORKDIR /root/

COPY --from=builder /app/coupon-system .

CMD ["./coupon-system"]