# Удаляем лишние сообщения вк

Бот для превращения беседы в своебразный канал аля telegram

удаляет все сообщения кроме сообщений администраторов беседы или пользователей указанных в конфиге


# Быстрый старт

- Получить токен сообщества
- Включить longpoll в настройках
- Создать `.env` файл по образцу `.env.example`
- Запустить 

## Запуск
Запуск
```shell
docker-compose up
```
Запуск в фоне
```shell
docker-compose up -d
```

### Запуск без docker
```shell
export TOKEN='*VK TOKEN*'
export ADMINS=*список админов*

go run main.go
```


# Формат логов
```
timestamp user_id: message (peer_id:conversation_message_id)
```
or
```
timestamp error
```