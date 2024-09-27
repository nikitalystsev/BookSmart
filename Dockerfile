FROM golang:latest

WORKDIR /usr/local/src

COPY ["./go.mod", "./go.sum", "./"]

COPY ./components ./components

RUN go mod download

COPY ./cmd/app ./cmd/app
COPY ./docs ./docs
COPY ./internal/app ./internal/app
COPY ./internal/config ./internal/config
COPY ./pkg ./pkg

RUN go build -o ./app cmd/app/main.go

CMD ["./app"]
