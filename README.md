# REST API QR-GENERATOR
Простой REST API сервер с аутентификацией для генерации qr-кодов 
## Конечные точки API
|HTTP-метод|Путь        |Описание                     |
|----------|------------|-----------------------------|
|POST      |/user/create|Создание пользователя        |
|POST      |/user/auth  |Аутентификация пользователя  |
|GET       |/user/me    |Кто я                        |
|DELETE    |/user/delete|Удаление пользователя        |
|POST      |/qr/create  |Создание qr-кода             |
|GET       |/qr/get     |Получение qr-кода по значению|
|GET       |/qr/all     |Список всех занчений qr-кодов|
|DELETE    |/qr/delete  |Удаление qr-кода             |
## Начало работы
### 1. Клонируйте репозиторий
``` bash
git clone github.com/konocha/qr-generate
```
### 2. Установите зависимости
``` bash
go mod download
```
### 3. Cоберите и запустите проект
Перейдите в директорию проекта  
* Сборка
``` bash
go build -v ./cmd/apiserver
```
* Запуск
``` bash
./apiserver
```
## Документация