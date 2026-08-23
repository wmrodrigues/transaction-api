FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /bin/api ./cmd/api
RUN go build -o /bin/db-migrate ./migration/cmd/db-migrate

FROM alpine:3.22

COPY --from=builder /bin/api /bin/api
COPY --from=builder /bin/db-migrate /bin/db-migrate

EXPOSE 8080

CMD ["/bin/api"]
