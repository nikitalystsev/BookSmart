FROM golang:latest

WORKDIR /usr/local/src

COPY ["./go.mod", "./go.sum", "./"]

# для разработки
COPY ./components/component-repo-mongo/go.mod ./components/component-repo-mongo/go.mod
COPY ./components/component-repo-mongo/go.sum ./components/component-repo-mongo/go.sum
COPY ./components/component-repo-postgres/go.mod ./components/component-repo-postgres/go.mod
COPY ./components/component-repo-postgres/go.sum ./components/component-repo-postgres/go.sum
COPY ./components/component-services/go.mod ./components/component-services/go.mod
COPY ./components/component-services/go.sum ./components/component-services/go.sum
COPY ./components/component-web-api/go.mod ./components/component-web-api/go.mod
COPY ./components/component-web-api/go.sum ./components/component-web-api/go.sum
COPY ./components/component-tech-ui/go.mod ./components/component-tech-ui/go.mod
COPY ./components/component-tech-ui/go.sum ./components/component-tech-ui/go.sum

RUN go mod download

# для разработки
COPY ./components/component-repo-mongo/ ./components/component-repo-mongo/
COPY ./components/component-repo-postgres/ ./components/component-repo-postgres/
COPY ./components/component-services/ ./components/component-services/
COPY ./components/component-web-api/ ./components/component-web-api/

COPY ./cmd/app ./cmd/app
COPY ./docs_swagger ./docs_swagger
COPY ./internal/app ./internal/app
COPY ./internal/config ./internal/config
COPY ./pkg ./pkg

RUN go build -o ./app cmd/app/main.go

CMD ["./app"]
