package models

type PlaceOrderReq struct {
	UserID int64  `json:"user_id"`
	IsBid  bool   `json:"is_bid"`
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
	Qty    string `json:"qty"`
	Type   string `json:"type"` // Market или Limit
}
