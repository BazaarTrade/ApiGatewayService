FROM golang:1.23.1 AS builder

WORKDIR /app
COPY . .

RUN go mod tidy

WORKDIR /app/cmd
RUN go build -o apiGateway .

FROM gcr.io/distroless/base

COPY --from=builder /app/cmd/apiGateway /apiGateway

EXPOSE 8080

CMD ["/apiGateway"]