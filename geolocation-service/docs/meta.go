// Package classification GeolocationService.
//
// Documentation of GeolocationService API.
//
//    Schemes:
//    - http
//    - https
//    BasePath: /
//    Version: 1.0.0
//
//    Consumes:
//    - application/json
//    - multipart/form-data
//
//    Produces:
//    - application/json
//
//    Security:
//    - basic
//
//    SecurityDefinitions:
//      Bearer:
//        type: apiKey
//        name: Authorization
//        in: header
//
// swagger:meta

package docs

//go:generate swagger generate spec -o ./docs/swagger.yaml --scan-models

// Эндпоинты

// swagger:route POST /geolocation/update geolocation updateGeolocationRequest
// Обновление геолокации текущего пользователя.
// responses:
//   200: updateGeolocationResponse
//   400: description: Bad request
//   500: description: Internal server error

// swagger:parameters updateGeolocationRequest
type updateGeolocationRequest struct {
	// in:body
	Body UpdateGeolocationRequest
}

// swagger:response updateGeolocationResponse
type updateGeolocationResponse struct {
	// in:body
	Body UpdateGeolocationResponse
}

// swagger:route GET /geolocation/nearby geolocation findNearbyUsersRequest
// Поиск пользователей поблизости.
// responses:
//   200: findNearbyUsersResponse
//   400: description: Bad request
//   500: description: Internal server error

// swagger:parameters findNearbyUsersRequest
type findNearbyUsersRequest struct {
	// in:query
	Radius float64 `json:"radius"`
}

// swagger:response findNearbyUsersResponse
type findNearbyUsersResponse struct {
	// in:body
	Body []NearbyUser
}

// swagger:route GET /geolocation/user/{id} geolocation getUserLocationRequest
// Получение текущей геолокации пользователя по его ID.
// responses:
//   200: getUserLocationResponse
//   400: description: Bad request
//   404: description: Not found
//   500: description: Internal server error

// swagger:parameters getUserLocationRequest
type getUserLocationRequest struct {
	// in:path
	ID string `json:"id"`
}

// swagger:response getUserLocationResponse
type getUserLocationResponse struct {
	// in:body
	Body UserLocation
}

// swagger:route POST /geolocation/share geolocation shareLocationRequest
// Поделиться своей геолокацией с другим пользователем.
// responses:
//   200: shareLocationResponse
//   400: description: Bad request
//   500: description: Internal server error

// swagger:parameters shareLocationRequest
type shareLocationRequest struct {
	// in:body
	Body ShareLocationRequest
}

// swagger:response shareLocationResponse
type shareLocationResponse struct {
	// in:body
	Body ShareLocationResponse
}

// swagger:route POST /geolocation/share/stop geolocation stopSharingLocationRequest
// Прекратить делиться своей геолокацией с другим пользователем.
// responses:
//   200: stopSharingLocationResponse
//   400: description: Bad request
//   500: description: Internal server error

// swagger:parameters stopSharingLocationRequest
type stopSharingLocationRequest struct {
	// in:body
	Body StopSharingLocationRequest
}

// swagger:response stopSharingLocationResponse
type stopSharingLocationResponse struct {
	// in:body
	Body StopSharingLocationResponse
}

// swagger:route POST /geolocation/privacy geolocation setLocationPrivacyRequest
// Настройка конфиденциальности геолокации.
// responses:
//   200: setLocationPrivacyResponse
//   400: description: Bad request
//   500: description: Internal server error

// swagger:parameters setLocationPrivacyRequest
type setLocationPrivacyRequest struct {
	// in:body
	Body SetLocationPrivacyRequest
}

// swagger:response setLocationPrivacyResponse
type setLocationPrivacyResponse struct {
	// in:body
	Body SetLocationPrivacyResponse
}

// swagger:route GET /geolocation/history geolocation getLocationHistoryRequest
// Получение истории геолокаций текущего пользователя.
// responses:
//   200: getLocationHistoryResponse
//   500: description: Internal server error

// swagger:response getLocationHistoryResponse
type getLocationHistoryResponse struct {
	// in:body
	Body []LocationHistory
}

// swagger:route DELETE /geolocation/history geolocation clearLocationHistoryRequest
// Очистка истории геолокаций текущего пользователя.
// responses:
//   200: clearLocationHistoryResponse
//   500: description: Internal server error

// swagger:response clearLocationHistoryResponse
type clearLocationHistoryResponse struct {
	// in:body
	Body ClearLocationHistoryResponse
}
