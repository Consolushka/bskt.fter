#!/bin/bash

echo "Starting application setup..."

# Выполняем миграции один раз при старте
echo "Running database migrations..."
goose up

# Создаем лог файл для cron
touch /var/log/app-cron.log

# Настраиваем cron задачу (каждую минуту для тестирования)
echo "Setting up cron job..."
echo "*/5 * * * * cd ../ && /build >> /var/log/app-cron.log 2>&1" > /etc/crontabs/root

# Даем права на выполнение cron файлу
chmod 0644 /etc/crontabs/root

# Загружаем crontab для root пользователя
crontab /etc/crontabs/root

# Запускаем cron в фоновом режиме с логированием
echo "Starting cron daemon..."
crond -f &

# Показываем текущую настройку cron
echo "Current cron configuration:"
crontab -l

# Выводим сообщение о том, что контейнер готов
echo "Container is ready! Cron job will run every minute (for testing)."
echo "Logs will be written to /var/log/app-cron.log"

# Ждем несколько секунд, чтобы cron запустился
sleep 5

# Проверяем, запущен ли cron
if pgrep crond > /dev/null; then
    echo "Cron daemon is running successfully!"
else
    echo "Warning: Cron daemon may not be running properly!"
fi

# Держим контейнер живым и показываем логи в реальном времени
tail -f /var/log/app-cron.log
