
# Для запуска достаточно сделать `docker compose up --build`

## NB! Так как запускается все в докере и файлы сохраняются в докере, то после перезапуска docker compose все файлы потеряются

## Для проверки фукционала
1. POST http://localhost:7000/upload с полем file и заголовком `Content-Type: multipart/form-data` - роут для загрузки файла, вернет имя файла, по которому его можно запросить
2. GET http://localhost:7000/:filename - filename - имя файла из прошлого запроса, возвращает загруженный файл

### Допущения
1. Равномерное заполнение серверов хранения - файл делится на 6 примерно равных частей, в данной реализации не учитывается свободное место на серверах
2. Для добавления сервера, нужно добавить его в таблицу `servers` и в `docker-compose.yml`
