FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates bash ffmpeg tzdata

ENV TZ=Asia/Seoul

WORKDIR /app

COPY --from=builder /app /app
COPY . .

EXPOSE 5100

CMD ["./main"]
