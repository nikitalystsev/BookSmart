#!/bin/bash

# Проверяем, был ли передан аргумент
if [ -z "$1" ]; then
    echo "Использование: $0 <номер_теста>"
    echo "Доступные тесты: test_01, test_02, test_03"
    exit 1
fi

# Определяем папку на основе переданного аргумента
TEST_FOLDER="./internal/tests_for_testing/benchmarks/test_$1"

# Проверяем, существует ли указанная папка
if [ ! -d "$TEST_FOLDER" ]; then
    echo "Ошибка: Папка '$TEST_FOLDER' не найдена."
    exit 1
fi

# Удаляем старые логи
rm -rf "$TEST_FOLDER/phout_gin.log"
rm -rf "$TEST_FOLDER/phout_echo.log"

# Указываем количество повторений
REPEAT_COUNT=1  # Замените 5 на нужное вам количество повторений

make run-app

for ((i=1; i<=REPEAT_COUNT; i++)); do

    echo "Запуск итерации $i..."

    # Запускаем первый процесс в фоновом режиме
    ./internal/tests_for_testing/benchmarks/pandora.exe "$TEST_FOLDER/load_gin.yml" &

    # Запускаем второй процесс в фоновом режиме
    ./internal/tests_for_testing/benchmarks/pandora.exe "$TEST_FOLDER/load_echo.yml" &

    # Ждем завершения обоих процессов
    wait

    echo "Итерация $i завершена."

    # Добавляем задержку в 10 секунд
    sleep 10
done

echo "Все процессы завершены."

