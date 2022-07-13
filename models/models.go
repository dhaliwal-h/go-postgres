package models

type Stock struct {
	StockID int64  `json:"stockid"`
	Name    string `json:"name"`
	Company string `json:"company"`
	Price   int64  `json:"price"`
}
