package entity

type Item struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Quantity int32  `json:"quantity"`
	PriceID  string `json:"priceID"`
}

type ItemWithQuantity struct {
	ID       string `json:"id"`
	Quantity int32  `json:"quantity"`
}
