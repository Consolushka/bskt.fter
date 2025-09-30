#!/bin/bash

echo "Starting application setup..."

# Выполняем миграции один раз при старте
echo "Running database migrations..."
#goose up

# Создаем лог файл для cron
touch /var/log/app-cron.log

# Создаем wrapper скрипт для запуска приложения
echo "Creating app runner script..."
cat > /app-runner.sh << 'EOF'
#!/bin/bash

# Загружаем переменные окружения
if [ -f /etc/environment ]; then
    . /etc/environment
fi

# Переходим в рабочую директорию
cd /imp

# Логируем информацию о запуске
echo "=== App execution started at $(date) ===" >> /var/log/app-cron.log
echo "Current PATH: $PATH" >> /var/log/app-cron.log
echo "Current PWD: $(pwd)" >> /var/log/app-cron.log
echo "Environment variables:" >> /var/log/app-cron.log
env | sort >> /var/log/app-cron.log
echo "========================" >> /var/log/app-cron.log

# Запускаем приложение
/build >> /var/log/app-cron.log 2>&1

echo "=== App execution finished at $(date) ===" >> /var/log/app-cron.log
EOF

# Делаем скрипт исполняемым
chmod +x /app-runner.sh

# Сохраняем текущие переменные окружения
echo "Saving environment variables..."
printenv > /etc/environment

# Настраиваем простую cron задачу
echo "Setting up cron job..."
echo "*/5 * * * * cd ../ && /build >> /var/log/app-cron.log 2>&1" > /etc/crontabs/root

# Даем права на выполнение cron файлу
chmod 0644 /etc/crontabs/root

# Загружаем crontab для root пользователя
crontab /etc/crontabs/root

# Запускаем cron в фоновом режиме с логированием
echo "Starting cron daemon..."
crond -f -s 2 &

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

# Тестируем wrapper скрипт один раз
echo "Testing app runner script..."
/app-runner.sh

# Держим контейнер живым и показываем логи в реальном времени
tail -f /var/log/app-cron.log
