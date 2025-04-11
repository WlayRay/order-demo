package dto

type CreateOrderResponse struct {
	CustomerID  string `json:"customerID"`
	OrderID     string `json:"orderID"`
	RedirectURL string `json:"redirectURL"`
}
