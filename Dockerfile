FROM golang:1.20.0-alpine

ENV CGO_ENABLED=1

WORKDIR /app

RUN apk add gcc musl-dev

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o erlog

CMD ["./erlog"]