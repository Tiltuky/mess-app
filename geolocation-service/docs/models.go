package docs

// Модели данных

type UpdateGeolocationRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type UpdateGeolocationResponse struct {
	Message string `json:"message"`
}

type NearbyUser struct {
	ID       string  `json:"id"`
	Username string  `json:"username"`
	Distance float64 `json:"distance"`
}

type UserLocation struct {
	ID        string  `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type ShareLocationRequest struct {
	UserID    string  `json:"userId"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type ShareLocationResponse struct {
	Message string `json:"message"`
}

type StopSharingLocationRequest struct {
	UserID string `json:"userId"`
}

type StopSharingLocationResponse struct {
	Message string `json:"message"`
}

type SetLocationPrivacyRequest struct {
	Visibility string `json:"visibility"`
}

type SetLocationPrivacyResponse struct {
	Message string `json:"message"`
}

type LocationHistory struct {
	Timestamp string  `json:"timestamp"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type ClearLocationHistoryResponse struct {
	Message string `json:"message"`
}
