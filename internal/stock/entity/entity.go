package entity

type ItemWithQuantity struct {
	ID       string `json:"ID,omitempty"`
	Quantity int32  `json:"Quantity,omitempty"`
}

type Item struct {
	ID       string `json:"ID,omitempty"`
	Name     string `json:"Name,omitempty"`
	Quantity int32  `json:"Quantity,omitempty"`
	PriceID  string `json:"PriceID,omitempty"`
}
type ItemInfo struct {
	Name     string `json:"Name"`
	Price    string `json:"Price"`
	CreateAt string `json:"CreateAt"`
	UpdateAt string `json:"UpdateAt"`
}
