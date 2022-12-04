FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -v -o /bin/app .

RUN rm -rf /app

ENV GIN_MODE=release

CMD ["app"]