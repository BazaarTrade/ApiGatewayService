package models

type Order struct {
	ID         int64  `json:"orderID"`
	UserID     int64  `json:"userID"`
	IsBid      bool   `json:"isBid"`
	Symbol     string `json:"symbol"`
	Price      string `json:"price"`
	Qty        string `json:"qty"`
	SizeFilled string `json:"sizeFilled"`
	Status     string `json:"status"`
	Type       string `json:"type"`
	CreatedAt  string `json:"created_at"`
	ClosedAt   string `json:"closed_at"`
}

type PlaceOrderReq struct {
	UserID int64  `json:"userID"`
	IsBid  bool   `json:"isBid"`
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
	Qty    string `json:"qty"`
	Type   string `json:"type"` // Market или Limit
}

type SubscriptionRequest struct {
	Topic    string `json:"topic"`
	Action   string `json:"action"`
	Symbol   string `json:"symbol"`
	Interval string `json:"interval"`
}
