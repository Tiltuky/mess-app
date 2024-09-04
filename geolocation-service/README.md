# API Endpoints for Geolocation

## Update Geolocation
- **Endpoint:** `POST /geolocation/update`
- **Description:** Обновление геолокации текущего пользователя.

## Find Nearby Users
- **Endpoint:** `GET /geolocation/nearby`
- **Description:** Поиск пользователей поблизости (на основе текущей геолокации).

## Get User Location
- **Endpoint:** `GET /geolocation/user/:id`
- **Description:** Получение текущей геолокации пользователя по его ID.

## Share Location
- **Endpoint:** `POST /geolocation/share`
- **Description:** Поделиться своей геолокацией с другим пользователем.

## Stop Sharing Location
- **Endpoint:** `POST /geolocation/share/stop`
- **Description:** Прекратить делиться своей геолокацией с другим пользователем.

## Set Location Privacy
- **Endpoint:** `POST /geolocation/privacy`
- **Description:** Настройка конфиденциальности геолокации (например, видимость для других пользователей).

## Get Location History
- **Endpoint:** `GET /geolocation/history`
- **Description:** Получение истории геолокаций текущего пользователя.

## Clear Location History
- **Endpoint:** `DELETE /geolocation/history`
- **Description:** Очистка истории геолокаций текущего пользователя.
