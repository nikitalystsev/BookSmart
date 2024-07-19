FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

# Загружаем зависимостиs
RUN go mod download

# Копируем остальную часть приложения
COPY ./ ./

# Устанавливаем нужный модуль для тестирования
RUN go install github.com/golang/mock/mockgen@v1.6.0

# Команда по умолчанию для запуска тестов
CMD ["go", "test", "-v", "./internal/tests/unitTests/..."]
