package entities

import "time"
type Transaction struct {
	Id        int       `json:"id"`
	OdaNumber int       `json:"oda_number"`
	Status    string    `json:"status"`
	Price     float32   `json:"price"`
	TotalData int       `json:"total_data"`
	CreatedAt time.Time `json:"created_at"`
}