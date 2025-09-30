#!/bin/bash

echo "Starting application setup..."

# Выполняем миграции один раз при старте
echo "Running database migrations..."
goose up

# Создаем лог файл для cron
touch /var/log/app-cron.log

# Настраиваем cron задачу (каждый день в 12:00)
echo "Setting up cron job..."
echo "0 12 * * * cd /IMP && /build >> /var/log/app-cron.log 2>&1" > /etc/crontabs/root

# Даем права на выполнение cron файлу
chmod 0644 /etc/crontabs/root

# Запускаем cron в фоновом режиме
echo "Starting cron daemon..."
crond

# Показываем текущую настройку cron
echo "Current cron configuration:"
crontab -l

# Выводим сообщение о том, что контейнер готов
echo "Container is ready! App will run daily at 12:00 Moscow time."
echo "Logs will be written to /var/log/app-cron.log"

# Запускаем приложение один раз при старте (опционально)
echo "Running application once at startup..."
/build
