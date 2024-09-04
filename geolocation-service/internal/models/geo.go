package models

import "time"

type User struct {
	ID                  int64             `json:"id"  db:"id"`
	CustomerID          string            `json:"customer_id"  db:"customer_id"`
	StatusSubscription  string            `json:"status_subscription" db:"status_subscription"`
	SubscriptionEndDate time.Time         `json:"subscription_end_date" db:"subscription_end_date"`
	Name                string            `json:"name"  db:"name"`
	Email               string            `json:"email"  db:"email"`
	Address             Address           `json:"address" db:"address"`
	Privacy             string            `json:"privacy" db:"privacy"`
	H3Index             string            `json:"h3_index"  db:"h3_index"`
	CreatedAt           time.Time         `json:"created_at"  db:"created_at"`
	UpdatedAt           time.Time         `json:"updated_at"  db:"updated_at"`
	LHistory            []LocationHistory `json:"location_history"  db:"location_history"`
	Nearby              []NearbyUser      `json:"nearby"  db:"nearby"`
	Sharing             []LocationSharing `json:"sharing"  db:"sharing"`
}

type LocationHistory struct {
	Timestamp time.Time `json:"timestamp"  db:"timestamp"`
	Address   Address   `json:"address"  db:"adress"`
	H3Index   string    `json:"h3_index"  db:"h3_index"`
}

type LocationSharing struct {
	ID         int64     `json:"id" db:"id"`
	SharerID   int64     `json:"sharer_id" db:"sharer_id"`
	ReceiverID int64     `json:"receiver_id" db:"receiver_id"`
	StartTime  time.Time `json:"start_time" db:"start_time"`
	EndTime    time.Time `json:"end_time" db:"end_time"`
}

type Address struct {
	Lat float64 `json:"lat" db:"lat"`
	Lng float64 `json:"lng" db:"lng"`
}

type NearbyUser struct {
	ID       int64   `json:"id"`
	Distance float64 `json:"distance"`
}
