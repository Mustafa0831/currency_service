# Currency Service

## Описание

Этот проект представляет собой веб-сервис на Golang, который использует библиотеку `gorilla/mux` для создания API. Сервис получает данные о курсах валют из публичного API Национального банка Казахстана и сохраняет их в локальную базу данных MS SQL Server. Также предоставляет возможность извлекать данные из базы данных по заданным параметрам.

## Установка

### Требования

- Golang 1.18+
  - Docker
  - MS SQL Server

### Шаги по установке

1. Клонируйте репозиторий:

    ```sh
    git clone https://github.com/ваш_репозиторий/currency_service.git
    cd currency_service
    ```

   2. Запустите MS SQL Server в Docker:

       ```sh
       docker run -e "ACCEPT_EULA=Y" -e "SA_PASSWORD=YourStrong@Passw0rd" -p 1433:1433 --name sql1 -d mcr.microsoft.com/mssql/server:2019-latest
       ```

   3. Подключитесь к MS SQL Server и создайте базу данных TEST:

       ```sh
       sudo docker exec -it sql1 /opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P 'YourStrong@Passw0rd'
       ```

      В консоли `sqlcmd` выполните следующие команды:

       ```sql
       CREATE DATABASE TEST;
       GO
       ```

   4. Настройте файл конфигурации:

      Создайте файл `config.json` в корне проекта со следующим содержимым:

       ```json
       {
         "port": "8080",
         "db_connection": "sqlserver://sa:YourStrong@Passw0rd@localhost:1433?database=TEST"
       }
       ```

   5. Установите зависимости и запустите сервис:

       ```sh
       go mod tidy
       go run main.go
       ```

## API Методы

### Сохранение данных о курсах валют

**URL:** `/currency/save/{date}`

**Метод:** `GET`

**Параметры:**

- `date` - Дата в формате `DD.MM.YYYY`

**Пример запроса:**

```sh
curl -X GET http://localhost:8080/currency/save/15.04.2021
```
Успешный ответ:
```sh
{
"success": true
}
```

### Получение данных о курсах валют
**URL:** `/currency/{date}/{code}`

**Метод:** `GET`

**Параметры:**  


- `date` - Дата в формате `DD.MM.YYYY`
- `code` - Код валюты (опционально)
**Пример запроса:**
```sh
curl -X GET http://localhost:8080/currency/15.04.2021/AUD
```


**Успешный ответ:**
```json 
[
  {
    "ID": 1,
    "Title": "Австралийский доллар",
    "Code": "AUD",
    "Value": 267.39,
    "ADate": "2021-04-15"
  }
]
```
## Тестирование
Для запуска тестов выполните команду:
```sh
go test ./handlers
```
## Документация
API документировано с использованием Swagger. Для генерации документации запустите:
```sh 
swag init
```
Откройте swagger.yaml или swagger.json в любом Swagger UI инструменте.

## Контакт
Для вопросов и предложений, пожалуйста, свяжитесь по адресу mr.mus0831@gmail.com.
```sh 

Теперь файл `README.md` форматирован правильно и готов к использованию.
```