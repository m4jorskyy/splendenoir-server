package order

type OrderResponse struct {
	ClientSecret string `json:"client_secret"`
	OrderID      int64  `json:"order_id"`
}
