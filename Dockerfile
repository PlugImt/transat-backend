FROM golang:1.23.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o transat-backend main.go

EXPOSE 3000

ENV PORT=3000

ENTRYPOINT [ "./transat-backend" ]