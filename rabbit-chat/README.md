# Консольный сетевой чат на RabbitMQ

▶ Сервер для обмена сообщениями, о котором договариваются
клиенты

▶ Центральный сервер, задаваемый как параметр командной
строки (127.0.0.1 по умолчанию)

▶ Есть именованные каналы, на которые можно переключаться и
постить туда

▶ Должна быть команда подписки на канал, типа «!switch channel1»

▶ Начальный канал принимается как аргумент командной строки

▶ Переключение на несуществующий канал должно его создавать

▶ Нет истории, получать только те сообщения, что были опубликованы с момента подключения

## Запуск клиента
И корня проекта необходимо вызвать следующую команду:
```bash
docker pull rabbitmq
docker run -d --name rabbit-chat -p 15672:15672 -p 5672:5672 rabbitmq
go run .
```