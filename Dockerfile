FROM golang:1.20-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o server ./cmd/server

FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/server /app/server
RUN apk add --no-cache nuget
EXPOSE 8080
CMD ["./server"]