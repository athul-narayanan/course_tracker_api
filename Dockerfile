FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o course-tracker-api .

EXPOSE 8080

FROM alpine:3.19
WORKDIR /app
COPY  --from=builder /app/course-tracker-api .

EXPOSE 8080

CMD ["./course-tracker-api"]



