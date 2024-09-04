package models

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type LoginResponse struct {
	Code      int    `json:"code"`
	Token     string `json:"token"`
	ExpiredAt int64  `json:"expiredAt"`
}
