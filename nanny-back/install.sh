#!/bin/bash
# Скрипт быстрой установки Nanny Backend

set -e

echo "🚀 Установка Nanny Backend..."
echo ""

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Проверка Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go не установлен. Установите Go 1.21+ и попробуйте снова.${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Go найден: $(go version)${NC}"

# Проверка PostgreSQL
if ! command -v psql &> /dev/null; then
    echo -e "${RED}❌ PostgreSQL не установлен. Установите PostgreSQL и попробуйте снова.${NC}"
    exit 1
fi
echo -e "${GREEN}✅ PostgreSQL найден${NC}"

# Установка зависимостей
echo ""
echo "📦 Установка зависимостей..."
go mod download
echo -e "${GREEN}✅ Зависимости установлены${NC}"

# Проверка .env файла
if [ ! -f .env ]; then
    echo ""
    echo "⚙️  Создание .env файла..."
    cp .env.example .env
    echo -e "${YELLOW}⚠️  Не забудьте отредактировать .env файл с вашими настройками!${NC}"
    echo "   Особенно: DB_PASSWORD"
fi

# Создание базы данных
echo ""
read -p "Создать базу данных nanny_db? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "📊 Создание базы данных..."
    
    # Попытка создать БД
    if createdb nanny_db 2>/dev/null; then
        echo -e "${GREEN}✅ База данных nanny_db создана${NC}"
    else
        echo -e "${YELLOW}⚠️  База данных уже существует или нет прав${NC}"
    fi
    
    # Применение схемы
    echo "📊 Применение SQL схемы..."
    if psql -d nanny_db -f scripts/schema.sql > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Схема применена${NC}"
    else
        echo -e "${YELLOW}⚠️  Не удалось применить схему автоматически${NC}"
        echo "   Выполните вручную: psql -d nanny_db -f scripts/schema.sql"
    fi
fi

# Проверка установки
echo ""
echo "🔍 Проверка установки..."
if go build -o /tmp/nanny-test cmd/api/main.go 2>/dev/null; then
    echo -e "${GREEN}✅ Проект компилируется успешно${NC}"
    rm /tmp/nanny-test
else
    echo -e "${RED}❌ Ошибка компиляции${NC}"
    exit 1
fi

# Инструкции
echo ""
echo -e "${GREEN}🎉 Установка завершена!${NC}"
echo ""
echo "Следующие шаги:"
echo "1. Отредактируйте .env файл (особенно DB_PASSWORD)"
echo "2. Запустите сервер:"
echo "   ${GREEN}make run${NC}"
echo "   или"
echo "   ${GREEN}go run cmd/api/main.go${NC}"
echo ""
echo "3. Проверьте работу:"
echo "   ${GREEN}curl http://localhost:8080/api/auth/login${NC}"
echo ""
echo "📚 Документация:"
echo "   - README.md - Основная документация"
echo "   - QUICKSTART.md - Быстрый старт"
echo "   - docs/API_EXAMPLES.md - Примеры API"
echo ""
echo "💡 Полезные команды:"
echo "   ${GREEN}make help${NC} - Список всех команд"
echo ""
