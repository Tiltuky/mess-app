# API Endpoints for Messages and Chats

## Get Messages
- **Endpoint:** `GET /messages/:chat_id`
- **Description:** Получение сообщений в чате по ID чата.

## Send Message
- **Endpoint:** `POST /messages/:chat_id`
- **Description:** Отправка сообщения в чат по ID чата.

## Delete Message
- **Endpoint:** `DELETE /messages/:chat_id/:msg_id`
- **Description:** Удаление сообщения по ID сообщения в чате.

## Get Chats
- **Endpoint:** `GET /chats`
- **Description:** Получение списка чатов текущего пользователя.

## Create Chat
- **Endpoint:** `POST /chats`
- **Description:** Создание нового чата.

## Get Chat Information
- **Endpoint:** `GET /chats/:id`
- **Description:** Получение информации о чате по ID.

## Delete Chat
- **Endpoint:** `DELETE /chats/:id`
- **Description:** Удаление чата по ID.
