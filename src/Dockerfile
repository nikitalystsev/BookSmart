FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./ ./

RUN go install github.com/golang/mock/mockgen@v1.6.0

CMD ["go", "test", "-v", "./internal/tests/unitTests/..."]
