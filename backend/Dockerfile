FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o chat cmd/chat/main.go

FROM alpine:latest
RUN apk --no-cache add libc6-compat

WORKDIR /app

COPY --from=builder /app/chat /app/chat
COPY --from=builder /app/files/ /app/files/
RUN chmod -R 755 /app/files

CMD ["/app/chat"]